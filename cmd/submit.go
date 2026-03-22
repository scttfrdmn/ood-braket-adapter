package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-braket-adapter/internal/braket"
	"github.com/scttfrdmn/ood-braket-adapter/internal/ood"
	"github.com/spf13/cobra"
)

var (
	submitOutputBucket string
	submitOutputPrefix string
)

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a Braket quantum task or hybrid job",
	Long: `Reads a JSON job spec from stdin and submits it to Amazon Braket.

For quantum tasks (circuits), the spec must include an OpenQASM 3.0 or
Braket IR circuit. For hybrid jobs, the spec must include a container image
or algorithm script S3 URI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec ood.JobSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode job spec: %w", err)
		}

		dev := device
		if dev == "" {
			dev = spec.Device
		}
		if dev == "" {
			return fmt.Errorf("--device is required (or set device in job spec)")
		}

		bucket := submitOutputBucket
		if bucket == "" {
			bucket = spec.OutputBucket
		}
		if bucket == "" {
			return fmt.Errorf("--output-bucket is required for Braket task results")
		}

		ctx := context.Background()
		client, err := braket.New(ctx, region)
		if err != nil {
			return err
		}

		taskID, err := client.CreateQuantumTask(ctx, dev, spec.Circuit, bucket, submitOutputPrefix, spec.Shots)
		if err != nil {
			return err
		}

		fmt.Println(taskID)
		return nil
	},
}

func init() {
	submitCmd.Flags().StringVar(&submitOutputBucket, "output-bucket", "", "S3 bucket for task results")
	submitCmd.Flags().StringVar(&submitOutputPrefix, "output-prefix", "braket-results", "S3 key prefix for results")
	rootCmd.AddCommand(submitCmd)
}
