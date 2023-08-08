package main

import (
	"context"
	"net/http"

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
// Method: POST
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus-unit/create
//
// Sample API Payload:
// 	[
// 	  {
// 	    "bus_id": "BCBSCMPN-875011",
// 	    "code": "BCBSCMPNBUS001",
// 	    "active": true,
// 	    "min_capacity": 30,
// 	    "max_capacity": 60
// 	  }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		unitList    []schema.BusUnit
		failedUnits schema.FailedBusUnits
	)

	// Unmarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), &unitList)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus Unit"}, utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusBadRequest(err)
	}

	for _, unit := range unitList {
		err = unit.ValidateMaximumCapacity(*unit.MinCapacity)
		if err != nil {
			failedUnits.SetFailedUnits(unit, err.Error())
			unit.Error(err, "InvalidBusCapacity", "max_capacity is invalid")

			continue
		}

		// Checks whether the bus unit exist or not
		busUnitExist, err := validate.IsBusUnitExisting(ctx, unit.BusID, unit.Code)
		if err != nil {
			failedUnits.SetFailedUnits(unit, "failed to validate bus unit if it exist")
			unit.Error(err, "IsBusUnitExisting", "failed to validate bus unit if it exist")

			continue
		}

		// If the bus unit exists, continue to the next item
		if busUnitExist {
			utility.Info("BusUnitExisting", "already existing bus unit", utility.KVP{Key: "unit", Value: unit})
			continue
		}

		// Set default values of the bus line information
		unit.SetValues()

		// Inserts a new bus unit record to the DynamoDB
		err = query.CreateBusUnit(ctx, unit)
		if err != nil {
			failedUnits.SetFailedUnits(unit, "failed to create a new bus unit record")
			unit.Error(err, "DynamoDBError", "failed to create a new bus unit record")

			continue
		}
	}

	if len(failedUnits.Failed) > 0 {
		return api.Response(http.StatusBadRequest, failedUnits), nil
	}

	return api.StatusOKWithoutBody()
}
