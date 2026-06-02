// Package cmd implements the ood-braket-adapter CLI.
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/scttfrdmn/ood-braket-adapter/internal/awscfg"
	"github.com/spf13/cobra"
)

var (
	region string
	device string
)

var (
	assumeRoleArn     string // #78: per-user role to assume (empty = instance role)
	assumeRoleExtID   string
	assumeRoleSession string
)

var version = "dev" // overridden at release time via -ldflags -X .../cmd.version

var rootCmd = &cobra.Command{
	Version: version,
	Use:     "ood-braket-adapter",
	Short:   "OOD adapter for Amazon Braket quantum/hybrid jobs",
	Long: `ood-braket-adapter submits, monitors, and cancels Amazon Braket
quantum tasks and hybrid jobs on behalf of Open OnDemand.

Braket supports three device types:
  QPU        — real quantum hardware (IonQ, Rigetti, Oxford, IQM)
  Simulator  — managed simulators (SV1, TN1, DM1) and on-demand (local)
  Hybrid     — PennyLane hybrid quantum-classical jobs on EC2

Job specs are read from stdin as JSON; results are written to stdout as JSON.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&region, "region", "us-east-1", "AWS region")
	rootCmd.PersistentFlags().StringVar(&device, "device", "", "Braket device ARN (e.g. arn:aws:braket:::device/quantum-simulator/amazon/sv1)")
}

// #78: per-user cross-account AssumeRole flags. Empty (default) = use the OOD instance role.
func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&assumeRoleArn, "assume-role-arn", "", "IAM role ARN to assume for AWS calls (empty = use the instance role)")
	pf.StringVar(&assumeRoleExtID, "assume-role-external-id", "", "sts:ExternalId for the assumed-role trust policy")
	pf.StringVar(&assumeRoleSession, "assume-role-session-name", "", "RoleSessionName for the assumed role (e.g. the OOD username)")
}

// awsOptions builds the AWS config options from the root flags (region + optional AssumeRole).
func awsOptions(ctx context.Context) []func(*config.LoadOptions) error {
	return awscfg.LoadOptions(ctx, awscfg.Options{
		Region:        region,
		AssumeRoleARN: assumeRoleArn,
		ExternalID:    assumeRoleExtID,
		SessionName:   assumeRoleSession,
	})
}
