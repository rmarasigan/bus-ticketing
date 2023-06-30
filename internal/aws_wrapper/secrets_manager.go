package awswrapper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

var (
	smClient *secretsmanager.Client
)

// initSMClient initializes the Secrets Manager Client from the
// provided configuration.
func initSMClient(ctx context.Context) {
	if smClient != nil {
		return
	}

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AWS_REGION))
	if err != nil {
		utility.Error(err, "SMClientError", "failed to load the default config")
		return
	}

	// Using the cfg value to create the Secrets Manager client
	smClient = secretsmanager.NewFromConfig(cfg)
}

// SecretGetValue initializes the Secrets Manager client and retrieves
// the contents of the encrypted fields.
func SecretGetValue(ctx context.Context, secretId string) (*secretsmanager.GetSecretValueOutput, error) {
	// Initialize the SecretManager Client
	initSMClient(ctx)

	var input = &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretId),
	}

	output, err := smClient.GetSecretValue(ctx, input)
	if err != nil {
		return nil, err
	}

	return output, nil
}
