package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/pkg/api"
	busroute "github.com/rmarasigan/bus-ticketing/pkg/handlers/bus_route"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, events *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	method := events.HTTPMethod

	switch method {
	case "GET":
		return busroute.Get(ctx, events)

	case "POST":
		return busroute.Post(ctx, events)

	default:
		return api.StatusUnhandledMethod()
	}
}
