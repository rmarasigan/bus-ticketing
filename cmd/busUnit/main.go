package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, events *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return api.StatusOK("")
}
