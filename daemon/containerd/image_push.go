package containerd

import (
	"context"
	"io"

	"github.com/containerd/containerd"
	containerdimages "github.com/containerd/containerd/images"
	"github.com/containerd/containerd/images/converter"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/remotes"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types/registry"
	"github.com/google/uuid"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// PushImage initiates a push operation on the repository named localName.
func (i *ImageService) PushImage(ctx context.Context, image, tag string, metaHeaders map[string][]string, authConfig *registry.AuthConfig, outStream io.Writer) error {
	// TODO: Pass this from user?
	platformMatcher := platforms.All

	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return err
	}
	if tag != "" {
		// Push by digest is not supported, so only tags are supported.
		ref, err = reference.WithTag(ref, tag)
		if err != nil {
			return err
		}
	}

	is := i.client.ImageService()
	store := i.client.ContentStore()

	img, err := is.Get(ctx, ref.String())
	if err != nil {
		return errors.Wrap(err, "Failed to get image")
	}

	target := img.Target

	if containerdimages.IsIndexType(target.MediaType) {
		if platformMatcher == platforms.All {
			// Create an image with manifest list that references only present blobs.
			fullest, err := i.createFullestImage(ctx, img)
			if err != nil {
				return err
			}
			if fullest != nil {
				target = fullest.Target
				defer i.client.ImageService().Delete(ctx, fullest.Name, containerdimages.SynchronousDelete())
			}
		} else {
			// If user requests to push a specific platform, convert the index to one that
			// doesn't reference manifest of other platforms.
			tmpRef := ref.String() + "-tmp-" + uuid.NewString()
			platformImg, err := converter.Convert(ctx, i.client, tmpRef, ref.String(), converter.WithPlatform(platformMatcher))
			if err != nil {
				return errors.Wrap(err, "failed to convert image to platform specific")
			}

			target = platformImg.Target
			defer i.client.ImageService().Delete(ctx, platformImg.Name, containerdimages.SynchronousDelete())
		}
	}

	jobs := newJobs()

	imageHandler := containerdimages.HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) (subdescs []ocispec.Descriptor, err error) {
		logrus.WithField("desc", desc).Debug("Pushing")
		if desc.MediaType != containerdimages.MediaTypeDockerSchema1Manifest {
			children, err := containerdimages.Children(ctx, store, desc)

			if err == nil {
				for _, c := range children {
					jobs.Add(c)
				}
			}

			jobs.Add(desc)
		}
		return nil, nil
	})
	imageHandler = remotes.SkipNonDistributableBlobs(imageHandler)

	resolver, tracker := newResolverFromAuthConfig(authConfig)

	finishProgress := showProgress(ctx, jobs, outStream, pushProgress(tracker))
	defer finishProgress()

	logrus.WithField("desc", target).WithField("ref", ref.String()).Info("Pushing desc to remote ref")
	err = i.client.Push(ctx, ref.String(), target,
		containerd.WithResolver(resolver),
		containerd.WithPlatformMatcher(platformMatcher),
		containerd.WithImageHandler(imageHandler),
	)

	return err
}
