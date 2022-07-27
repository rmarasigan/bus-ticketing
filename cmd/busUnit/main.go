package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
	busunit "github.com/rmarasigan/bus-ticketing/pkg/handlers/bus_unit"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, events *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	method := events.HTTPMethod

	switch method {
	case "GET":
		return busunit.Get(ctx, events)

	case "POST":
		return busunit.Post(ctx, events)

	default:
		return api.StatusUnhandledMethod()
	}
}
