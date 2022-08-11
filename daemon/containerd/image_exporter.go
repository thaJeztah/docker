package containerd

import (
	"context"
	"io"

	"github.com/containerd/containerd"
	containerdimages "github.com/containerd/containerd/images"
	"github.com/containerd/containerd/images/archive"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/platforms"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/container"
	"github.com/docker/docker/pkg/containerfs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (i *ImageService) PerformWithBaseFS(ctx context.Context, c *container.Container, fn func(containerfs.ContainerFS) error) error {
	snapshotter := i.client.SnapshotService(i.snapshotter)
	mounts, err := snapshotter.Mounts(ctx, c.ID)
	if err != nil {
		return err
	}
	return mount.WithTempMount(ctx, mounts, func(root string) error {
		return fn(containerfs.NewLocalContainerFS(root))
	})
}

// LoadImage uploads a set of images into the repository. This is the
// complement of ExportImage.  The input stream is an uncompressed tar
// ball containing images and metadata.
func (i *ImageService) LoadImage(ctx context.Context, inTar io.ReadCloser, outStream io.Writer, quiet bool) error {
	platform := platforms.All
	imgs, err := i.client.Import(ctx, inTar, containerd.WithImportPlatform(platform))

	if err != nil {
		logrus.WithError(err).Error("Failed to import image to containerd")
		return errors.Wrapf(err, "Failed to import image")
	}

	for _, img := range imgs {
		platformImg := containerd.NewImageWithPlatform(i.client, img, platform)

		unpacked, err := platformImg.IsUnpacked(ctx, i.snapshotter)
		if err != nil {
			logrus.WithError(err).WithField("image", img.Name).Error("IsUnpacked failed")
			continue
		}

		if !unpacked {
			err := platformImg.Unpack(ctx, i.snapshotter)
			if err != nil {
				logrus.WithError(err).WithField("image", img.Name).Error("Failed to unpack image")
				return errors.Wrapf(err, "Failed to unpack image")
			}
		}
	}
	return nil
}

// ExportImage exports a list of images to the given output stream. The
// exported images are archived into a tar when written to the output
// stream. All images with the given tag and all versions containing
// the same tag are exported. names is the set of tags to export, and
// outStream is the writer which the images are written to.
func (i *ImageService) ExportImage(ctx context.Context, names []string, outStream io.Writer) error {
	opts := []archive.ExportOpt{
		archive.WithSkipNonDistributableBlobs(),
	}

	is := i.client.ImageService()

	for _, name := range names {
		ref, err := reference.ParseDockerRef(name)
		if err != nil {
			return err
		}

		img, err := is.Get(ctx, ref.String())
		if err != nil {
			return err
		}

		converted, err := i.createFullestImage(ctx, img)
		if err != nil {
			return err
		}

		if converted != nil {
			opts = append(opts, archive.WithImage(is, converted.Name))
			defer is.Delete(context.Background(), converted.Name, containerdimages.SynchronousDelete())
			continue
		}

		opts = append(opts, archive.WithImage(is, img.Name))
	}

	return i.client.Export(ctx, outStream, opts...)
}
