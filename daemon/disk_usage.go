package daemon // import "github.com/docker/docker/daemon"

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/docker/docker/api/server/router/system"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/errdefs"
)

type usageCollector func(ctx context.Context, daemon *Daemon, du *types.DiskUsage) error

var collectors = map[system.DiskUsageObject]usageCollector{
	system.ContainerObject:  appendContainers,
	system.ImageObject:      appendImages,
	system.VolumeObject:     appendVolumes,
	system.BuildCacheObject: nil, // build-cache is currently collected in systemRouter.getDiskUsage()
}

// SystemDiskUsage returns information about the daemon data disk usage
func (daemon *Daemon) SystemDiskUsage(ctx context.Context, opts system.DiskUsageOptions) (*types.DiskUsage, error) {
	if !atomic.CompareAndSwapInt32(&daemon.diskUsageRunning, 0, 1) {
		return nil, fmt.Errorf("a disk usage operation is already running")
	}
	defer atomic.StoreInt32(&daemon.diskUsageRunning, 0)

	if len(opts.ObjectTypes) == 0 {
		// Collect all objects by default
		opts.ObjectTypes = map[system.DiskUsageObject]bool{
			system.ContainerObject: true,
			system.ImageObject:     true,
			system.VolumeObject:    true,
		}
	}

	du := &types.DiskUsage{}
	for t := range opts.ObjectTypes {
		c, ok := collectors[t]
		if !ok {
			return nil, errdefs.InvalidParameter(fmt.Errorf("unknown object type: %s", t))
		}
		if c == nil {
			continue
		}
		if err := c(ctx, daemon, du); err != nil {
			return nil, err
		}
	}
	return du, nil
}

// appendContainers retrieves container list, and adds it to du
func appendContainers(_ context.Context, daemon *Daemon, du *types.DiskUsage) error {
	containers, err := daemon.Containers(&types.ContainerListOptions{
		Size: true,
		All:  true,
	})
	if err != nil {
		return fmt.Errorf("failed to retrieve container list: %v", err)
	}
	du.Containers = containers
	return nil
}

// appendImages retrieves image list and layersSize, and adds it to du
func appendImages(ctx context.Context, daemon *Daemon, du *types.DiskUsage) error {
	images, err := daemon.imageService.Images(filters.NewArgs(), false, true)
	if err != nil {
		return fmt.Errorf("failed to retrieve image list: %v", err)
	}

	layersSize, err := daemon.imageService.LayerDiskUsage(ctx)
	if err != nil {
		return err
	}
	du.Images = images
	du.LayersSize = layersSize
	return nil
}

// appendVolumes retrieves volumes list, and adds it to du
func appendVolumes(ctx context.Context, daemon *Daemon, du *types.DiskUsage) error {
	volumes, err := daemon.volumes.LocalVolumesSize(ctx)
	if err != nil {
		return err
	}
	du.Volumes = volumes
	return nil
}
