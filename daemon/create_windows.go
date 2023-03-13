package daemon // import "github.com/docker/docker/daemon"

import (
	"context"
	"path/filepath"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/container"
	volumeopts "github.com/docker/docker/volume/service/opts"
	"github.com/sirupsen/logrus"
)

// createContainerOSSpecificSettings performs host-OS specific container create functionality
func (daemon *Daemon) createContainerOSSpecificSettings(container *container.Container, config *containertypes.Config, hostConfig *containertypes.HostConfig) error {
	if containertypes.Isolation.IsDefault(hostConfig.Isolation) {
		// Make sure the host config has the default daemon isolation if not specified by caller.
		hostConfig.Isolation = daemon.defaultIsolation
	}

	// Look for volume-paths defined in the image config, and attach anonymous
	// volumes for those paths that do not have a volume or bind-mount mounted
	// at the location.
	for dest := range config.Volumes {
		// TODO(thaJeztah): should we allow `/some-path` to be used on Windows, and treat it as `C:/some-path` ?
		destination := filepath.Clean(dest)

		// Skip volumes for which we already have something mounted on that
		// destination because of a --volume-from.
		// TODO(thaJeztah): on Unix this is HasMountFor() - why do we have different names for Unix and Windows?
		if container.IsDestinationMounted(destination) {
			logrus.WithField("container", container.ID).WithField("destination", dest).Debug("mountpoint already exists, skipping anonymous volume")
			// Not an error, this could easily have come from the image config.
			continue
		}

		// TODO(thaJeztah); is this indeed only intended to be used for anonymous volumes?
		v, err := daemon.volumes.Create(context.TODO(), "", hostConfig.VolumeDriver, volumeopts.WithCreateReference(container.ID))
		if err != nil {
			return err
		}

		// FIXME Windows: This code block is present in the Linux version and
		// allows the contents to be copied to the container FS prior to it
		// being started. However, the function utilizes the FollowSymLinkInScope
		// path which does not cope with Windows volume-style file paths. There
		// is a separate effort to resolve this (@swernli), so this processing
		// is deferred for now. A case where this would be useful is when
		// a dockerfile includes a VOLUME statement, but something is created
		// in that directory during the dockerfile processing. What this means
		// on Windows for TP5 is that in that scenario, the contents will not
		// copied, but that's (somewhat) OK as HCS will bomb out soon after
		// at it doesn't support mapped directories which have contents in the
		// destination path anyway.
		//
		// Example for repro later:
		//   FROM windowsservercore
		//   RUN mkdir c:\myvol
		//   RUN copy c:\windows\system32\ntdll.dll c:\myvol
		//   VOLUME "c:\myvol"
		//
		// Then
		//   docker build -t vol .
		//   docker run -it --rm vol cmd  <-- This is where HCS will error out.
		//
		//	// never attempt to copy existing content in a container FS to a shared volume
		//	if v.DriverName() == volume.DefaultDriverName {
		//		if err := container.CopyImagePathContent(v, mp.Destination); err != nil {
		//			return err
		//		}
		//	}

		// Add it to container.MountPoints
		container.AddMountPointWithVolume(destination, &volumeWrapper{v: v, s: daemon.volumes}, true)
	}
	return nil
}
