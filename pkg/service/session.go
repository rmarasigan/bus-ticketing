package service

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

var (
	SQSClient         sqsiface.SQSAPI
	EventBridgeClient *eventbridge.EventBridge
	DynamoDBClient    dynamodbiface.DynamoDBAPI
)

// newSession creates a new session with custom configuration value.
func newSession() *session.Session {
	// Create a config
	config := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	// Create Session with custom configuration.
	sess, err := session.NewSession(config)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "SessionError", Message: "Failed to create a session"})
		return nil
	}

	return sess
}

// SqsSession creates a new instance of the SQS client with a session if not set.
func SqsSession() {
	sess := newSession()

	// Checks if SQSClient session is not set
	if SQSClient == nil {
		// Create an SQS Client from just a session
		SQSClient = sqs.New(sess)
	}

	SQSClient = sqs.New(sess)
}

// DynamodbSession creates a new instance of the DynamoDB client with a session if not set.
func DynamodbSession() {
	sess := newSession()

	// Checks if DynamoDBClient session is not set
	if DynamoDBClient == nil {
		// Create a DynamoDB Client from just a session
		DynamoDBClient = dynamodb.New(sess)
	}

	DynamoDBClient = dynamodb.New(sess)
}

// EventBridgeSession creates a new instance of EventBridge client with a session if not set.
func EventBridgeSession() {
	sess := newSession()

	// Checks if EventBridgeClient session is not set
	if EventBridgeClient == nil {
		// Create an EventBridge Client from just a session
		EventBridgeClient = eventbridge.New(sess)
	}

	EventBridgeClient = eventbridge.New(sess)
}
