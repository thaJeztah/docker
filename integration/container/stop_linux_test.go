package container // import "github.com/docker/docker/integration/container"

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/docker/docker/integration/internal/container"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/poll"
)

// TestStopContainerWithTimeout checks that ContainerStop with
// a timeout works as documented, i.e. in case of negative timeout
// waiting is not limited (issue #35311).
func TestStopContainerWithTimeout(t *testing.T) {
	defer setupTest(t)()
	client := testEnv.APIClient()
	ctx := context.Background()

	testCmd := container.WithCmd("sh", "-c", "sleep 2 && exit 42")
	testData := []struct {
		doc              string
		timeout          int
		expectedExitCode int
	}{
		// In case container is forcefully killed, 137 is returned,
		// otherwise the exit code from the above script
		{
			"zero timeout: expect forceful container kill",
			0, 137,
		},
		{
			"too small timeout: expect forceful container kill",
			1, 137,
		},
		{
			"big enough timeout: expect graceful container stop",
			3, 42,
		},
		{
			"unlimited timeout: expect graceful container stop",
			-1, 42,
		},
	}

	for _, d := range testData {
		d := d
		t.Run(strconv.Itoa(d.timeout), func(t *testing.T) {
			t.Parallel()
			id := container.Run(ctx, t, client, testCmd)

			err := client.ContainerStop(ctx, id, containertypes.StopOptions{Timeout: &d.timeout})
			assert.NilError(t, err)

			poll.WaitOn(t, container.IsStopped(ctx, client, id),
				poll.WithDelay(100*time.Millisecond))

			inspect, err := client.ContainerInspect(ctx, id)
			assert.NilError(t, err)
			assert.Equal(t, inspect.State.ExitCode, d.expectedExitCode)
		})
	}
}

// TestStopContainerWithTimeoutCancel checks that ContainerStop with a timeout
// cancels the stop, but does not forcefully kill the container.
// See issue https://github.com/moby/moby/issues/45731
func TestStopContainerWithTimeoutCancel(t *testing.T) {
	t.Parallel()
	defer setupTest(t)()

	const (
		cancelContext   = "cancel context"
		closeConnection = "close connection"
	)

	for _, tc := range []string{cancelContext, closeConnection} {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			t.Parallel()
			apiClient := testEnv.APIClient()
			defer apiClient.Close()
			ctx := context.Background()

			testCmd := container.WithCmd("sh", "-c", "sleep 30")
			id := container.Run(ctx, t, apiClient, testCmd)

			ctxCancel, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				stopTimeout := 15
				err := apiClient.ContainerStop(ctxCancel, id, containertypes.StopOptions{Timeout: &stopTimeout})
				switch tc {
				case cancelContext:
					assert.Check(t, is.ErrorType(err, errdefs.IsCancelled))
				case closeConnection:
					assert.Check(t, err)
				}
			}()

			// Give the ContainerStop some time to make sure it's being handled,
			// then cancel the context or close the client to cancel the stop.
			time.Sleep(1 * time.Second)
			switch tc {
			case cancelContext:
				// cancel the context
				cancel()
			case closeConnection:
				// close the client connection
				assert.Check(t, apiClient.Close())
			}

			// The stop should be cancelled, and the container should still
			// be running
			inspect, err := apiClient.ContainerInspect(ctx, id)
			assert.Check(t, err)
			assert.Check(t, inspect.State.Running)

			err = apiClient.ContainerRemove(ctx, id, types.ContainerRemoveOptions{Force: true})
			assert.Check(t, err)
		})

	}

}
