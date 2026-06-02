// Package ood defines OOD job spec types for the Braket adapter.
package ood

// JobSpec is the OOD job submission payload for Braket tasks.
type JobSpec struct {
	// Circuit is an OpenQASM 3.0 string or Braket IR JSON representing the quantum circuit.
	Circuit string `json:"circuit"`
	// Device is the Braket device ARN. Overridden by --device flag.
	Device string `json:"device,omitempty"`
	// Shots is the number of measurement shots. 0 = adapter default (1000).
	Shots int32 `json:"shots,omitempty"`
	// OutputBucket is the S3 bucket for task results. Overridden by --output-bucket flag.
	OutputBucket string `json:"output_bucket,omitempty"`
	// JobName is used for tagging and display in OOD.
	JobName string `json:"job_name,omitempty"`
	// Env contains additional environment variables passed to hybrid jobs.
	Env map[string]string `json:"env,omitempty"`
}

// JobStatus maps Braket task states to OOD status strings.
type JobStatus struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

const (
	StatusQueued    = "queued"
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
	StatusUnknown   = "undetermined"
)
