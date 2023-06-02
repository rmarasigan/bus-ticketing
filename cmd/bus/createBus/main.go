package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rmarasigan/bus-ticketing/api"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/app/validate"
	"github.com/rmarasigan/bus-ticketing/internal/utility"
)

func main() {
	lambda.Start(handler)
}

// It receives the Amazon API Gateway event record data as input, validates the
// request body, saves the validated request body to the DynamoDB Table, and
// responds with a 200 OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/create
//
// Sample API Payload:
// 	{
// 		"name": "Thunder Rail Bus Line",
// 		"owner": "Thando Oyibo Emmett",
// 		"company": "Rail Bus Way",
// 		"address": "1986 Bogisich Junctions, Hamillhaven, Kansas",
// 		"email": "thando.emmet@outlook.com",
// 		"mobile_number": "+1-335-908-1432"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var bus = new(schema.Bus)

	// Unmarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), bus)
	if err != nil {
		bus.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusBadRequest(err)
	}

	// Validate if the required fields are not empty
	err = validate.CreateBusLine(*bus)
	if err != nil {
		bus.Error(err, "CreateBusLine", "missing required field(s)")
		return api.StatusBadRequest(err)
	}

	// Checks whether the bus line exist or not
	busLineExist, err := validate.IsBusLineExisting(ctx, bus.Name, bus.Company)
	if err != nil {
		bus.Error(err, "IsBusLineExisting", "failed to validate bus line if it exist")
		return api.StatusInternalServerError()
	}

	// If the bus line exists, return a 400 BadRequet HTTP Status
	if busLineExist {
		err := fmt.Errorf("%s bus line from %s company already exist", bus.Name, bus.Company)
		bus.Error(err, "IsBusLineExisting", "already existing bus line")

		return api.StatusBadRequest(err)
	}

	// Set default values of the bus line information
	bus.SetValues()

	// Inserts a new bus line record to the DynamoDB
	err = query.CreateBusLine(ctx, bus)
	if err != nil {
		bus.Error(err, "DynamoDBError", "failed to create a new bus line record")
		return api.StatusInternalServerError()
	}

	return api.StatusOKWithoutBody()
}
