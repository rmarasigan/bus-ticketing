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
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus_unit/create
//
// Sample API Payload:
// 	{
// 		"code": "RLBSWV1_0606",
// 		"bus_id": "RLBSW-856996",
// 		"active": true,
// 		"min_capacity": 40,
// 		"max_capacity": 50
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var unit = new(schema.BusUnit)

	// Unmarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), unit)
	if err != nil {
		unit.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusBadRequest(err)
	}

	// Validate if the required fields are not empty
	err = validate.CreateBusUnitFields(*unit)
	if err != nil {
		unit.Error(err, "CreateBusUnitFields", "missing required field(s)")
		return api.StatusBadRequest(err)
	}

	// Checks whether the bus unit exist or not
	busUnitExist, err := validate.IsBusUnitExisting(ctx, unit.BusID, unit.Code)
	if err != nil {
		unit.Error(err, "IsBusUnitExisting", "failed to validate bus unit if it exist")
		return api.StatusInternalServerError()
	}

	// If the bus unit exists, return a 400 BadRequest HTTP Status
	if busUnitExist {
		err := fmt.Errorf("%s bus unit from %s already exist", unit.Code, unit.BusID)
		unit.Error(err, "IsBusUnitExisting", "already existing bus unit")

		return api.StatusBadRequest(err)
	}

	// Set default values of the bus line information
	unit.SetValues()

	// Inserts a new bus unit record to the DynamoDB
	err = query.CreateBusUnit(ctx, unit)
	if err != nil {
		unit.Error(err, "DynamoDBError", "failed to create a new bus unit record")
		return api.StatusInternalServerError()
	}

	return api.StatusOKWithoutBody()
}
