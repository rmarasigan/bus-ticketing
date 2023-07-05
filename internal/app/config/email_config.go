package config

import (
	"context"
	"errors"
	"os"

	"github.com/rmarasigan/bus-ticketing/internal/app/email"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

// GetEmailConfig checks if the Email Secrets Manager is configured on the environment,
// fetches the Email Secret values, and returns the email configuration.
func GetEmailConfig(ctx context.Context) (cfg email.Configuration, err error) {
	var emailsecret = os.Getenv("EMAIL_SECRET")

	// Check if the Email SecretsManager is configured
	if emailsecret == "" {
		err = errors.New("secretsmanager EMAIL_SECRET environment variable is not set")
		cfg.Error(err, "SMError", "secretsmanager EMAIL_SECRET is not configured on the environment")

		return
	}

	// Get the email secret values
	result, err := awswrapper.SecretGetValue(ctx, emailsecret)
	if err != nil {
		cfg.Error(err, "SMError", "failed to fetch the email secret")
		return
	}

	// Unmarshal the email secret values
	err = utility.ParseJSON([]byte(*result.SecretString), &cfg)
	if err != nil {
		cfg.Error(err, "JSONError", "failed to unmarshal the JSON-encoded secret")
		return
	}

	return
}
