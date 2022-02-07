package daemon // import "github.com/docker/docker/daemon"

import "github.com/docker/docker/api/types/versions/v1p20"

// containerInspect120 serializes the master version of a container into a json type.
func (daemon *Daemon) containerInspect120(name string) (*v1p20.ContainerJSON, error) {
	ctr, err := daemon.GetContainer(name)
	if err != nil {
		return nil, err
	}

	ctr.Lock()
	defer ctr.Unlock()

	base, err := daemon.getInspectData(ctr)
	if err != nil {
		return nil, err
	}

	return &v1p20.ContainerJSON{
		ContainerJSONBase: base,
		Mounts:            ctr.GetMountPoints(),
		Config: &v1p20.ContainerConfig{
			Config:          ctr.Config,
			MacAddress:      ctr.Config.MacAddress,
			NetworkDisabled: ctr.Config.NetworkDisabled,
			ExposedPorts:    ctr.Config.ExposedPorts,
			VolumeDriver:    ctr.HostConfig.VolumeDriver,
		},
		NetworkSettings: daemon.getBackwardsCompatibleNetworkSettings(ctr.NetworkSettings),
	}, nil
}
