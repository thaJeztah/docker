package containerd

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/content"
	cerrdefs "github.com/containerd/containerd/errdefs"
	containerdimages "github.com/containerd/containerd/images"
	"github.com/containerd/containerd/images/converter"
	"github.com/containerd/containerd/platforms"
	"github.com/docker/distribution/reference"
	containertypes "github.com/docker/docker/api/types/container"
	imagetype "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/image"
	"github.com/docker/docker/layer"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var shortID = regexp.MustCompile(`^([a-f0-9]{4,64})$`)

// GetContainerdImage returns the containerd image corresponding to the image referred to by refOrID.
// The platform parameter is currently ignored
func (i *ImageService) GetContainerdImage(ctx context.Context, refOrID string, platform *ocispec.Platform) (containerdimages.Image, error) {
	return i.resolveImageName2(ctx, refOrID)
}

// GetImage returns an image corresponding to the image referred to by refOrID.
func (i *ImageService) GetImage(ctx context.Context, refOrID string, options imagetype.GetImageOpts) (*image.Image, error) {
	ii, img, err := i.getImage(ctx, refOrID)
	if err != nil {
		return nil, err
	}

	if options.Details {
		size, err := ii.Size(ctx)
		if err != nil {
			return nil, err
		}

		tagged, err := i.client.ImageService().List(ctx, fmt.Sprintf("target.digest==%s", ii.Target().Digest.String()))
		if err != nil {
			return nil, err
		}
		tags := make([]reference.Named, 0, len(tagged))
		for _, i := range tagged {
			name, err := reference.ParseNamed(i.Name)
			if err != nil {
				return nil, err
			}
			tags = append(tags, name)
		}

		img.Details = &image.Details{
			References:  tags,
			Size:        size,
			Metadata:    nil,
			Driver:      i.StorageDriver(),
			LastUpdated: ii.Metadata().UpdatedAt,
		}
	}
	return img, err
}

func (i *ImageService) getImage(ctx context.Context, refOrID string) (containerd.Image, *image.Image, error) {
	desc, err := i.ResolveImage(ctx, refOrID)
	if err != nil {
		return nil, nil, err
	}

	ctrdimg, err := i.resolveImageName2(ctx, refOrID)
	if err != nil {
		return nil, nil, err
	}
	ii := containerd.NewImage(i.client, ctrdimg)
	provider := i.client.ContentStore()
	conf, err := ctrdimg.Config(ctx, provider, ii.Platform())
	if err != nil {
		return nil, nil, err
	}

	var ociimage ocispec.Image
	imageConfigBytes, err := content.ReadBlob(ctx, ii.ContentStore(), conf)
	if err != nil {
		return nil, nil, err
	}

	if err := json.Unmarshal(imageConfigBytes, &ociimage); err != nil {
		return nil, nil, err
	}

	fs, err := ii.RootFS(ctx)
	if err != nil {
		return nil, nil, err
	}
	rootfs := image.NewRootFS()
	for _, id := range fs {
		rootfs.Append(layer.DiffID(id))
	}
	exposedPorts := make(nat.PortSet, len(ociimage.Config.ExposedPorts))
	for k, v := range ociimage.Config.ExposedPorts {
		exposedPorts[nat.Port(k)] = v
	}
	return ii, &image.Image{
		V1Image: image.V1Image{
			ID:           string(desc.Digest),
			OS:           ociimage.OS,
			Architecture: ociimage.Architecture,
			Config: &containertypes.Config{
				Entrypoint:   ociimage.Config.Entrypoint,
				Env:          ociimage.Config.Env,
				Cmd:          ociimage.Config.Cmd,
				User:         ociimage.Config.User,
				WorkingDir:   ociimage.Config.WorkingDir,
				ExposedPorts: exposedPorts,
				Volumes:      ociimage.Config.Volumes,
			},
		},
		RootFS: rootfs,
	}, nil
}

// ResolveImage searches for an image based on the given
// reference or identifier. Returns the descriptor of
// the image, could be manifest list, manifest, or config.
func (i *ImageService) ResolveImage(ctx context.Context, refOrID string) (d ocispec.Descriptor, err error) {
	d, _, err = i.resolveImageName(ctx, refOrID)
	return
}

func (i *ImageService) resolveImageName2(ctx context.Context, refOrID string) (img containerdimages.Image, err error) {
	parsed, err := reference.ParseAnyReference(refOrID)
	if err != nil {
		return img, errdefs.InvalidParameter(err)
	}

	is := i.client.ImageService()

	namedRef, ok := parsed.(reference.Named)
	if !ok {
		digested, ok := parsed.(reference.Digested)
		if !ok {
			return img, errdefs.InvalidParameter(errors.New("bad reference"))
		}

		imgs, err := is.List(ctx, fmt.Sprintf("target.digest==%s", digested.Digest()))
		if err != nil {
			return img, errors.Wrap(err, "failed to lookup digest")
		}
		if len(imgs) == 0 {
			return img, errdefs.NotFound(errors.New("image not found with digest"))
		}

		return imgs[0], nil
	}

	namedRef = reference.TagNameOnly(namedRef)

	// If the identifier could be a short ID, attempt to match
	if shortID.MatchString(refOrID) {
		ref := namedRef.String()
		filters := []string{
			fmt.Sprintf("name==%q", ref),
			fmt.Sprintf(`target.digest~=/sha256:%s[0-9a-fA-F]{%d}/`, refOrID, 64-len(refOrID)),
		}
		imgs, err := is.List(ctx, filters...)
		if err != nil {
			return img, err
		}

		if len(imgs) == 0 {
			return img, errdefs.NotFound(errors.New("list returned no images"))
		}
		if len(imgs) > 1 {
			digests := map[digest.Digest]struct{}{}
			for _, img := range imgs {
				if img.Name == ref {
					return img, nil
				}
				digests[img.Target.Digest] = struct{}{}
			}

			if len(digests) > 1 {
				return img, errdefs.NotFound(errors.New("ambiguous reference"))
			}
		}

		if imgs[0].Name != ref {
			namedRef = nil
		}
		return imgs[0], nil
	}
	img, err = is.Get(ctx, namedRef.String())
	if err != nil {
		// TODO(containerd): error translation can use common function
		if !cerrdefs.IsNotFound(err) {
			return img, err
		}
		return img, errdefs.NotFound(errors.New("id not found"))
	}

	return img, nil
}

func (i *ImageService) resolveImageName(ctx context.Context, refOrID string) (ocispec.Descriptor, reference.Named, error) {
	parsed, err := reference.ParseAnyReference(refOrID)
	if err != nil {
		return ocispec.Descriptor{}, nil, errdefs.InvalidParameter(err)
	}

	is := i.client.ImageService()

	namedRef, ok := parsed.(reference.Named)
	if !ok {
		digested, ok := parsed.(reference.Digested)
		if !ok {
			return ocispec.Descriptor{}, nil, errdefs.InvalidParameter(errors.New("bad reference"))
		}

		imgs, err := is.List(ctx, fmt.Sprintf("target.digest==%s", digested.Digest()))
		if err != nil {
			return ocispec.Descriptor{}, nil, errors.Wrap(err, "failed to lookup digest")
		}
		if len(imgs) == 0 {
			return ocispec.Descriptor{}, nil, errdefs.NotFound(errors.New("image not found with digest"))
		}

		return imgs[0].Target, nil, nil
	}

	namedRef = reference.TagNameOnly(namedRef)

	// If the identifier could be a short ID, attempt to match
	if shortID.MatchString(refOrID) {
		ref := namedRef.String()
		filters := []string{
			fmt.Sprintf("name==%q", ref),
			fmt.Sprintf(`target.digest~=/sha256:%s[0-9a-fA-F]{%d}/`, refOrID, 64-len(refOrID)),
		}
		imgs, err := is.List(ctx, filters...)
		if err != nil {
			return ocispec.Descriptor{}, nil, err
		}

		if len(imgs) == 0 {
			return ocispec.Descriptor{}, nil, errdefs.NotFound(errors.New("list returned no images"))
		}
		if len(imgs) > 1 {
			digests := map[digest.Digest]struct{}{}
			for _, img := range imgs {
				if img.Name == ref {
					return img.Target, namedRef, nil
				}
				digests[img.Target.Digest] = struct{}{}
			}

			if len(digests) > 1 {
				return ocispec.Descriptor{}, nil, errdefs.NotFound(errors.New("ambiguous reference"))
			}
		}

		if imgs[0].Name != ref {
			namedRef = nil
		}
		return imgs[0].Target, namedRef, nil
	}
	img, err := is.Get(ctx, namedRef.String())
	if err != nil {
		// TODO(containerd): error translation can use common function
		if !cerrdefs.IsNotFound(err) {
			return ocispec.Descriptor{}, nil, err
		}
		return ocispec.Descriptor{}, nil, errdefs.NotFound(errors.New("id not found"))
	}

	return img.Target, namedRef, nil
}

// PresentChildrenHandler traverses recursively all children descriptors that are present in the store.
func (i *ImageService) presentChildrenHandler() containerdimages.HandlerFunc {
	store := i.client.ContentStore()

	return func(ctx context.Context, desc ocispec.Descriptor) (subdescs []ocispec.Descriptor, err error) {
		_, err = store.ReaderAt(ctx, desc)
		if err != nil {
			if cerrdefs.IsNotFound(err) {
				return nil, nil
			}
			return nil, err
		}

		return containerdimages.Children(ctx, store, desc)
	}
}

// createFullestImage converts the given manifest list to a fullest possible with the current store content.
// If all manifests blobs are present in the content store, then no conversion is performed and it returns nil.
// Otherwise it creates a new image which targets a manifest list referencing only manifests with present blobs.
func (i *ImageService) createFullestImage(ctx context.Context, img containerdimages.Image) (*containerdimages.Image, error) {
	if containerdimages.IsIndexType(img.Target.MediaType) {
		// Check which platform manifests we have blobs for.
		presentPlatforms, missingPlatforms, err := i.checkPresentPlatforms(ctx, img.Target)

		// If we have all the manifests, no need to convert the image.
		if len(missingPlatforms) == 0 {
			return nil, nil
		}

		// Create a new manifest list which contains only the manifests we have in store.
		srcRef := img.Name
		targetRef := srcRef + "-tmp-" + uuid.NewString()
		for _, p := range presentPlatforms {
			targetRef += "-" + platforms.Format(p)
		}

		newImg, err := converter.Convert(ctx, i.client, targetRef, srcRef,
			converter.WithPlatform(platforms.Any(presentPlatforms...)))
		if err != nil {
			return nil, err
		}
		return newImg, nil
	}

	return nil, nil
}

// checkPresentPlatforms returns the platforms for which the image targeting an index/manifest list
// has all the blobs present in the content store.
// Both present and missing platforms list returned if there wasn't any error.
func (i *ImageService) checkPresentPlatforms(ctx context.Context, index ocispec.Descriptor) ([]ocispec.Platform, []ocispec.Platform, error) {
	if !containerdimages.IsIndexType(index.MediaType) {
		return nil, nil, errors.New("not an index/manifest list")
	}

	store := i.client.ContentStore()

	children, err := containerdimages.Children(ctx, store, index)
	if err != nil {
		return nil, nil, err
	}

	missingPlatforms := []ocispec.Platform{}
	presentPlatforms := []ocispec.Platform{}
	for _, child := range children {
		if containerdimages.IsManifestType(child.MediaType) {
			missing := false

			_, err := store.ReaderAt(ctx, child)
			if cerrdefs.IsNotFound(err) {
				missing = true
			} else if err != nil {
				return nil, nil, err
			}

			// If the manifest is not missing, also check for its blobs.
			if !missing {
				manifestChildren, err := containerdimages.Children(ctx, store, child)
				if err != nil {
					return nil, nil, err
				}

				// If any blob is missing, mark the manifest as missing.
				for _, manifestChild := range manifestChildren {
					_, err := store.ReaderAt(ctx, manifestChild)
					if cerrdefs.IsNotFound(err) {
						missing = true
						break
					} else if err != nil {
						return nil, nil, err
					}
				}
			}

			if missing {
				missingPlatforms = append(missingPlatforms, *child.Platform)
				logrus.WithField("digest", child.Digest.String()).Debug("missing blob")
			} else {
				presentPlatforms = append(presentPlatforms, *child.Platform)
			}
		}
	}
	return presentPlatforms, missingPlatforms, nil
}
