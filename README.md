# ood-braket-adapter

Open OnDemand adapter for [Amazon Braket](https://aws.amazon.com/braket/) — submit,
monitor, and cancel quantum tasks on QPUs and managed simulators from the OOD portal.

Part of the [aws-openondemand](https://github.com/scttfrdmn/aws-openondemand) ecosystem.

## What it does

Translates OOD job submissions into Amazon Braket quantum tasks. Users write circuits
in OpenQASM 3.0 or Python (via PennyLane/Qiskit), submit them through OOD, and get
results back in their home directory as S3-backed files.

Supported backends:
| Device | ARN prefix | Notes |
|--------|-----------|-------|
| SV1 (state vector simulator) | `arn:aws:braket:::device/quantum-simulator/amazon/sv1` | 34 qubits, exact |
| TN1 (tensor network simulator) | `arn:aws:braket:::device/quantum-simulator/amazon/tn1` | 50 qubits, large circuits |
| DM1 (density matrix simulator) | `arn:aws:braket:::device/quantum-simulator/amazon/dm1` | 17 qubits, noisy |
| IonQ Aria | `arn:aws:braket:us-east-1::device/qpu/ionq/Aria-1` | Real QPU |
| Rigetti Ankaa | `arn:aws:braket:us-west-1::device/qpu/rigetti/Ankaa-9Q-3` | Real QPU |

## Installation

```bash
go install github.com/scttfrdmn/ood-braket-adapter@latest
```

Or download a pre-built binary from [Releases](https://github.com/scttfrdmn/ood-braket-adapter/releases).

## Usage

### submit

Reads a JSON job spec from stdin, submits a quantum task, and prints the task ARN.

```bash
echo '{
  "circuit": "OPENQASM 3.0; qubit[2] q; h q[0]; cnot q[0], q[1];",
  "device": "arn:aws:braket:::device/quantum-simulator/amazon/sv1",
  "shots": 1000,
  "output_bucket": "my-braket-results",
  "job_name": "bell-state"
}' | ood-braket-adapter submit --region us-east-1
# arn:aws:braket:us-east-1:123456789012:quantum-task/abc123...
```

### status

```bash
ood-braket-adapter status arn:aws:braket:us-east-1:123456789012:quantum-task/abc123
# {"id":"arn:...","status":"completed"}
```

### delete

Cancels a quantum task (only possible while QUEUED or RUNNING):

```bash
ood-braket-adapter delete arn:aws:braket:us-east-1:123456789012:quantum-task/abc123
```

## Job spec fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `circuit` | string | yes | OpenQASM 3.0 circuit or Braket IR JSON |
| `device` | string | yes* | Braket device ARN (* or use `--device` flag) |
| `shots` | int | no | Measurement shots (default: 1000) |
| `output_bucket` | string | yes* | S3 bucket for results (* or use `--output-bucket`) |
| `job_name` | string | no | Display name in OOD job table |

## IAM requirements

The OOD instance role needs:
```json
{
  "Effect": "Allow",
  "Action": [
    "braket:CreateQuantumTask",
    "braket:GetQuantumTask",
    "braket:CancelQuantumTask",
    "braket:SearchQuantumTasks"
  ],
  "Resource": "*"
}
```
Plus S3 write access to the output bucket.

## OOD cluster configuration

```yaml
# /etc/ood/config/clusters.d/braket.yml
v2:
  metadata:
    title: Amazon Braket
  job:
    adapter: script
    submit_host: localhost
    bin: /usr/local/bin/ood-braket-adapter
    submit:
      script: |
        #!/bin/bash
        echo "$OOD_JOB_SPEC" | /usr/local/bin/ood-braket-adapter submit \
          --region "$AWS_REGION" \
          --device "$BRAKET_DEVICE_ARN" \
          --output-bucket "$BRAKET_OUTPUT_BUCKET"
```
