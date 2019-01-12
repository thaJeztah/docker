package node

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"gotest.tools/v3/poll"
)

// IsInState verifies if the current swarm node is in the given state
func IsInState(ctx context.Context, client client.APIClient, state swarm.LocalNodeState) func(log poll.LogT) poll.Result {
	return func(log poll.LogT) poll.Result {
		info, err := client.Info(ctx)
		if err != nil {
			return poll.Error(err)
		}
		if info.Swarm.LocalNodeState == state {
			return poll.Success()
		}
		return poll.Continue("waiting for node to be %q, currently %q", state, info.Swarm.LocalNodeState)
	}
}
