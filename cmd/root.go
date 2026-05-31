// Package cmd implements the ood-braket-adapter CLI.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	region string
	device string
)

var version = "dev" // overridden at release time via -ldflags -X .../cmd.version

var rootCmd = &cobra.Command{
	Version: version,
	Use:   "ood-braket-adapter",
	Short: "OOD adapter for Amazon Braket quantum/hybrid jobs",
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
