package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

var (
	sqsClient *sqs.Client
)

// initSQSClient initializes the SQS Client from the provided
// configuration.
func initSQSClient(ctx context.Context) {
	if sqsClient != nil {
		return
	}

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "SQSClientError", "failed to load the default config")
		return
	}

	// Using the cfg value to create the SQS client
	sqsClient = sqs.NewFromConfig(cfg)
}
