package system // import "github.com/docker/docker/api/server/router/system"

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

// DiskUsageOptions holds parameters for system disk usage query.
type DiskUsageOptions struct {
	// ObjectTypes to include when collecting usage.
	ObjectTypes map[DiskUsageObject]bool
}

// Backend is the methods that need to be implemented to provide
// system specific functionality.
type Backend interface {
	SystemInfo() *types.Info
	SystemVersion() types.Version
	SystemDiskUsage(ctx context.Context, opts DiskUsageOptions) (*types.DiskUsage, error)
	SubscribeToEvents(since, until time.Time, ef filters.Args) ([]events.Message, chan interface{})
	UnsubscribeFromEvents(chan interface{})
	AuthenticateToRegistry(ctx context.Context, authConfig *types.AuthConfig) (string, string, error)
}

// ClusterBackend is all the methods that need to be implemented
// to provide cluster system specific functionality.
type ClusterBackend interface {
	Info() swarm.Info
}
