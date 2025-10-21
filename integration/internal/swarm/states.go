package swarm

import (
	"context"
	"fmt"

	swarmtypes "github.com/moby/moby/api/types/swarm"
	"github.com/moby/moby/client"
	"gotest.tools/v3/poll"
)

// NoTasksForService verifies that there are no more tasks for the given service
func NoTasksForService(ctx context.Context, apiClient client.ServiceAPIClient, serviceID string) func(log poll.LogT) poll.Result {
	return func(log poll.LogT) poll.Result {
		result, err := apiClient.TaskList(ctx, client.TaskListOptions{
			Filters: make(client.Filters).Add("service", serviceID),
		})
		if err == nil {
			if len(result.Tasks) == 0 {
				return poll.Success()
			}
			if len(result.Tasks) > 0 {
				return poll.Continue("task count for service %s at %d waiting for 0", serviceID, len(result.Tasks))
			}
			return poll.Continue("waiting for tasks for service %s to be deleted", serviceID)
		}
		// TODO we should not use an error as indication that the tasks are gone. There may be other reasons for an error to occur.
		return poll.Success()
	}
}

// NoTasks verifies that all tasks are gone
func NoTasks(ctx context.Context, apiClient client.ServiceAPIClient) func(log poll.LogT) poll.Result {
	return func(log poll.LogT) poll.Result {
		result, err := apiClient.TaskList(ctx, client.TaskListOptions{})
		switch {
		case err != nil:
			return poll.Error(err)
		case len(result.Tasks) == 0:
			return poll.Success()
		default:
			return poll.Continue("waiting for all tasks to be removed: task count at %d", len(result.Tasks))
		}
	}
}

// RunningTasksCount verifies there are `instances` tasks running for `serviceID`
func RunningTasksCount(ctx context.Context, apiClient client.ServiceAPIClient, serviceID string, instances uint64) func(log poll.LogT) poll.Result {
	return func(log poll.LogT) poll.Result {
		result, err := apiClient.TaskList(ctx, client.TaskListOptions{
			Filters: make(client.Filters).Add("service", serviceID),
		})
		var running int
		var taskError string
		for _, task := range result.Tasks {
			switch task.Status.State {
			case swarmtypes.TaskStateRunning:
				running++
			case swarmtypes.TaskStateFailed, swarmtypes.TaskStateRejected:
				if task.Status.Err != "" {
					log.Logf("task %v on node %v %v: %v", task.ID, task.NodeID, task.Status.State, task.Status.Err)
					taskError = task.Status.Err
				}
			default:
				// not interested in other states.
			}
		}

		switch {
		case err != nil:
			return poll.Error(err)
		case running > int(instances):
			return poll.Continue("waiting for tasks to terminate")
		case running < int(instances) && taskError != "":
			return poll.Continue("waiting for tasks to enter run state. task failed with error: %s", taskError)
		case running == int(instances):
			return poll.Success()
		default:
			return poll.Continue("running task count at %d waiting for %d (total tasks: %d)", running, instances, len(result.Tasks))
		}
	}
}

// JobComplete is a poll function for determining that a ReplicatedJob is
// completed additionally, while polling, it verifies that the job never
// exceeds MaxConcurrent running tasks
func JobComplete(ctx context.Context, apiClient client.ServiceAPIClient, service swarmtypes.Service) func(log poll.LogT) poll.Result {
	filter := make(client.Filters).Add("service", service.ID)

	var jobIteration swarmtypes.Version
	if service.JobStatus != nil {
		jobIteration = service.JobStatus.JobIteration
	}

	maxConcurrent := int(*service.Spec.Mode.ReplicatedJob.MaxConcurrent)
	totalCompletions := int(*service.Spec.Mode.ReplicatedJob.TotalCompletions)
	previousResult := ""

	return func(log poll.LogT) poll.Result {
		result, err := apiClient.TaskList(ctx, client.TaskListOptions{
			Filters: filter,
		})
		if err != nil {
			poll.Error(err)
		}

		var running int
		var completed int

		var runningSlot []int
		var runningID []string

		for _, task := range result.Tasks {
			// make sure the task has the same job iteration
			if task.JobIteration == nil || task.JobIteration.Index != jobIteration.Index {
				continue
			}
			switch task.Status.State {
			case swarmtypes.TaskStateRunning:
				running++
				runningSlot = append(runningSlot, task.Slot)
				runningID = append(runningID, task.ID)
			case swarmtypes.TaskStateComplete:
				completed++
			default:
				// not interested in other states.
			}
		}

		switch {
		case running > maxConcurrent:
			return poll.Error(fmt.Errorf(
				"number of running tasks (%v) exceeds max (%v)", running, maxConcurrent,
			))
		case (completed + running) > totalCompletions:
			return poll.Error(fmt.Errorf(
				"number of tasks exceeds total (%v), %v running and %v completed",
				totalCompletions, running, completed,
			))
		case completed == totalCompletions && running == 0:
			return poll.Success()
		default:
			newRes := fmt.Sprintf(
				"Completed: %2d Running: %v\n\t%v",
				completed, runningSlot, runningID,
			)
			if newRes == previousResult {
			} else {
				previousResult = newRes
			}

			return poll.Continue(
				"Job not yet finished, %v completed and %v running out of %v total",
				completed, running, totalCompletions,
			)
		}
	}
}

func HasLeader(ctx context.Context, apiClient client.NodeAPIClient) func(log poll.LogT) poll.Result {
	return func(log poll.LogT) poll.Result {
		nodes, err := apiClient.NodeList(ctx, client.NodeListOptions{
			Filters: make(client.Filters).Add("role", "manager"),
		})
		if err != nil {
			return poll.Error(err)
		}
		for _, node := range nodes {
			if node.ManagerStatus != nil && node.ManagerStatus.Leader {
				return poll.Success()
			}
		}
		return poll.Continue("no leader elected yet")
	}
}
