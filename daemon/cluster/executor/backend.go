package executor // import "github.com/docker/docker/daemon/cluster/executor"

import (
	"context"
	"io"
	"time"

	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/backend"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	swarmtypes "github.com/docker/docker/api/types/swarm"
	containerpkg "github.com/docker/docker/container"
	clustertypes "github.com/docker/docker/daemon/cluster/provider"
	networkSettings "github.com/docker/docker/daemon/network"
	"github.com/docker/docker/libnetwork"
	"github.com/docker/docker/libnetwork/cluster"
	networktypes "github.com/docker/docker/libnetwork/types"
	"github.com/docker/docker/plugin"
	volumeopts "github.com/docker/docker/volume/service/opts"
	"github.com/docker/swarmkit/agent/exec"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Backend defines the executor component for a swarm agent.
type Backend interface {
	stateBackend
	CreateManagedNetwork(clustertypes.NetworkCreateRequest) error
	DeleteManagedNetwork(networkID string) error
	FindNetwork(idName string) (libnetwork.Network, error)
	SetupIngress(clustertypes.NetworkCreateRequest, string) (<-chan struct{}, error)
	ReleaseIngress() (<-chan struct{}, error)
	CreateManagedContainer(config types.ContainerCreateConfig) (container.ContainerCreateCreatedBody, error)
	// ContainerStart(name string, hostConfig *container.HostConfig, checkpoint string, checkpointDir string) error
	// ContainerStop(name string, seconds *int) error
	ContainerLogs(context.Context, string, *types.ContainerLogsOptions) (msgs <-chan *backend.LogMessage, tty bool, err error)
	ConnectContainerToNetwork(containerName, networkName string, endpointConfig *network.EndpointSettings) error
	ActivateContainerServiceBinding(containerName string) error
	DeactivateContainerServiceBinding(containerName string) error
	UpdateContainerServiceConfig(containerName string, serviceConfig *clustertypes.ServiceConfig) error
	ContainerInspectCurrent(name string, size bool) (*types.ContainerJSON, error)
	// ContainerWait(ctx context.Context, name string, condition containerpkg.WaitCondition) (<-chan containerpkg.StateStatus, error)
	// ContainerRm(name string, config *types.ContainerRmConfig) error
	// ContainerKill(name string, sig uint64) error
	SetContainerDependencyStore(name string, store exec.DependencyGetter) error
	SetContainerSecretReferences(name string, refs []*swarmtypes.SecretReference) error
	SetContainerConfigReferences(name string, refs []*swarmtypes.ConfigReference) error
	SystemInfo() *types.Info
	Containers(config *types.ContainerListOptions) ([]*types.Container, error)
	SetNetworkBootstrapKeys([]*networktypes.EncryptionKey) error
	DaemonJoinsCluster(provider cluster.Provider)
	DaemonLeavesCluster()
	IsSwarmCompatible() error
	SubscribeToEvents(since, until time.Time, filter filters.Args) ([]events.Message, chan interface{})
	UnsubscribeFromEvents(listener chan interface{})
	UpdateAttachment(string, string, string, *network.NetworkingConfig) error
	WaitForDetachment(context.Context, string, string, string, string) error
	PluginManager() *plugin.Manager
	PluginGetter() *plugin.Store
	GetAttachmentStore() *networkSettings.AttachmentStore
	HasExperimental() bool
}

// stateBackend includes functions to implement to provide container state lifecycle functionality.
type stateBackend interface {
	// ContainerCreate(config types.ContainerCreateConfig) (container.ContainerCreateCreatedBody, error)
	ContainerKill(name string, sig uint64) error
	// ContainerPause(name string) error
	// ContainerRename(oldName, newName string) error
	// ContainerResize(name string, height, width int) error
	// ContainerRestart(name string, seconds *int) error
	ContainerRm(name string, config *types.ContainerRmConfig) error
	ContainerStart(name string, hostConfig *container.HostConfig, checkpoint string, checkpointDir string) error
	ContainerStop(name string, seconds *int) error
	// ContainerUnpause(name string) error
	// ContainerUpdate(name string, hostConfig *container.HostConfig) (container.ContainerUpdateOKBody, error)
	ContainerWait(ctx context.Context, name string, condition containerpkg.WaitCondition) (<-chan containerpkg.StateStatus, error)
}

// monitorBackend includes functions to implement to provide containers monitoring functionality.
type monitorBackend interface {
	// ContainerChanges(name string) ([]archive.Change, error)
	ContainerInspect(name string, size bool, version string) (interface{}, error)
	ContainerLogs(ctx context.Context, name string, config *types.ContainerLogsOptions) (msgs <-chan *backend.LogMessage, tty bool, err error)
	ContainerStats(ctx context.Context, name string, config *backend.ContainerStatsConfig) error
	ContainerTop(name string, psArgs string) (*container.ContainerTopOKBody, error)

	Containers(config *types.ContainerListOptions) ([]*types.Container, error)
}

// VolumeBackend is used by an executor to perform volume operations
type VolumeBackend interface {
	Create(ctx context.Context, name, driverName string, opts ...volumeopts.CreateOption) (*types.Volume, error)
}

// ImageBackend is used by an executor to perform image operations
type ImageBackend interface {
	PullImage(ctx context.Context, image, tag string, platform *specs.Platform, metaHeaders map[string][]string, authConfig *types.AuthConfig, outStream io.Writer) error
	GetRepository(context.Context, reference.Named, *types.AuthConfig) (distribution.Repository, error)
	LookupImage(name string) (*types.ImageInspect, error)
}
