package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
	"github.com/rmarasigan/bus-ticketing/pkg/handlers/user"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, events *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var method = events.HTTPMethod

	switch method {
	case "GET":
		return user.Get(ctx, events)

	case "POST":
		return user.Post(ctx, events)

	default:
		return api.StatusUnhandledMethod()
	}
}
