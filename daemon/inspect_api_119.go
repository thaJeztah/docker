package daemon // import "github.com/docker/docker/daemon"

import "github.com/docker/docker/api/types/versions/v1p19"

// containerInspectPre120 gets containers for pre 1.20 APIs.
func (daemon *Daemon) containerInspectPre120(name string) (*v1p19.ContainerJSON, error) {
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

	volumes := make(map[string]string)
	volumesRW := make(map[string]bool)
	for _, m := range ctr.MountPoints {
		volumes[m.Destination] = m.Path()
		volumesRW[m.Destination] = m.RW
	}

	return &v1p19.ContainerJSON{
		ContainerJSONBase: base,
		Volumes:           volumes,
		VolumesRW:         volumesRW,
		Config: &v1p19.ContainerConfig{
			Config:          ctr.Config,
			MacAddress:      ctr.Config.MacAddress,
			NetworkDisabled: ctr.Config.NetworkDisabled,
			ExposedPorts:    ctr.Config.ExposedPorts,
			VolumeDriver:    ctr.HostConfig.VolumeDriver,
			Memory:          ctr.HostConfig.Memory,
			MemorySwap:      ctr.HostConfig.MemorySwap,
			CPUShares:       ctr.HostConfig.CPUShares,
			CPUSet:          ctr.HostConfig.CpusetCpus,
		},
		NetworkSettings: daemon.getBackwardsCompatibleNetworkSettings(ctr.NetworkSettings),
	}, nil
}
