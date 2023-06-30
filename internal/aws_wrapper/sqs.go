package awswrapper

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

var (
	sqsClient            *sqs.Client
	BOOKING_MSG_GROUP_ID = "process.booking"
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

// generateMessageDeduplicationID returns a token to be used for the SQS message
// deduplication ID.
func generateMessageDeduplicationID(message string) string {
	hash := md5.Sum([]byte(message))
	return hex.EncodeToString(hash[:])
}

// SQSSendMessage initializes the SQS client and delivers message to the specified queue.
func SQSSendMessage(ctx context.Context, queue, message, groupdId string) error {
	// Initlaize the SQS client.
	initSQSClient(ctx)

	var input = &sqs.SendMessageInput{
		QueueUrl:    aws.String(queue),
		MessageBody: aws.String(message),
	}

	if strings.Contains(queue, ".fifo") {
		input.MessageGroupId = aws.String(groupdId)
		input.MessageDeduplicationId = aws.String(generateMessageDeduplicationID(message))
	}

	_, err := sqsClient.SendMessage(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
