// Package braket wraps the Amazon Braket API for the OOD adapter.
package braket

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/braket"
)

// Client wraps the Amazon Braket client.
type Client struct {
	svc    *braket.Client
	region string
}

// New creates a Braket client using the default AWS credential chain.
func New(ctx context.Context, region string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: braket.NewFromConfig(cfg), region: region}, nil
}

// CreateQuantumTask submits a circuit as a quantum task to the specified device.
// circuit should be an OpenQASM 3.0 string or Braket IR JSON.
// shots is the number of measurement shots (0 = adapter default of 1000).
func (c *Client) CreateQuantumTask(ctx context.Context, deviceARN, circuit, outputBucket, outputPrefix string, shots int32) (string, error) {
	if shots == 0 {
		shots = 1000
	}
	out, err := c.svc.CreateQuantumTask(ctx, &braket.CreateQuantumTaskInput{
		DeviceArn:         aws.String(deviceARN),
		Action:            aws.String(circuit),
		OutputS3Bucket:    aws.String(outputBucket),
		OutputS3KeyPrefix: aws.String(outputPrefix),
		Shots:             aws.Int64(int64(shots)),
	})
	if err != nil {
		return "", fmt.Errorf("braket CreateQuantumTask: %w", err)
	}
	return aws.ToString(out.QuantumTaskArn), nil
}

// GetQuantumTask returns the current status and metadata of a quantum task.
func (c *Client) GetQuantumTask(ctx context.Context, taskARN string) (*braket.GetQuantumTaskOutput, error) {
	out, err := c.svc.GetQuantumTask(ctx, &braket.GetQuantumTaskInput{
		QuantumTaskArn: aws.String(taskARN),
	})
	if err != nil {
		return nil, fmt.Errorf("braket GetQuantumTask: %w", err)
	}
	return out, nil
}

// CancelQuantumTask cancels a quantum task that has not yet completed.
func (c *Client) CancelQuantumTask(ctx context.Context, taskARN string) error {
	_, err := c.svc.CancelQuantumTask(ctx, &braket.CancelQuantumTaskInput{
		QuantumTaskArn: aws.String(taskARN),
	})
	if err != nil {
		return fmt.Errorf("braket CancelQuantumTask: %w", err)
	}
	return nil
}
