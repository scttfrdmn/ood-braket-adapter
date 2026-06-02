package cmd

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/ood-braket-adapter/internal/braket"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <task-id>",
	Short: "Cancel a Braket quantum task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := braket.New(ctx, region, awsOptions(ctx)...)
		if err != nil {
			return err
		}
		if err := client.CancelQuantumTask(ctx, args[0]); err != nil {
			return err
		}
		fmt.Printf("task %s cancelled\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
