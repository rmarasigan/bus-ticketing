package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

var (
	evbClient *eventbridge.Client
)

// initEventBridgeClient initializes the EventBridge Service Client
// from the provided configuration.
func initEventBridgeClient(ctx context.Context) {
	if evbClient != nil {
		return
	}

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "EVBError", "failed to load the default config")
		return
	}

	// Using the cfg value to create the EventBridge client
	evbClient = eventbridge.NewFromConfig(cfg)
}
