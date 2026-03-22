package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/braket/types"
	"github.com/scttfrdmn/ood-braket-adapter/internal/braket"
	"github.com/scttfrdmn/ood-braket-adapter/internal/ood"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <task-arn>",
	Short: "Get the status of a Braket quantum task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := braket.New(ctx, region)
		if err != nil {
			return err
		}

		task, err := client.GetQuantumTask(ctx, args[0])
		if err != nil {
			return err
		}

		js := ood.JobStatus{
			ID:     args[0],
			Status: braketStateToOod(task.Status),
		}
		if task.FailureReason != nil {
			js.Message = *task.FailureReason
		}

		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

func braketStateToOod(s types.QuantumTaskStatus) string {
	switch s {
	case types.QuantumTaskStatusCreated, types.QuantumTaskStatusQueued:
		return ood.StatusQueued
	case types.QuantumTaskStatusRunning:
		return ood.StatusRunning
	case types.QuantumTaskStatusCompleted:
		return ood.StatusCompleted
	case types.QuantumTaskStatusFailed:
		return ood.StatusFailed
	case types.QuantumTaskStatusCancelled, types.QuantumTaskStatusCancelling:
		return ood.StatusCancelled
	default:
		return ood.StatusUnknown
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
