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
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/create
//
// Sample API Payload:
// 	[
// 	  {
// 	    "name": "Blue Horizon",
// 	    "owner": "John Doe",
// 	    "email": "john.doe@example.com",
// 	    "address": "123 Main Street, City",
// 	    "company": "ABC Bus Company",
// 	    "mobile_number": "123-456-7890"
// 	  },
// 	  {
// 	    "name": "Green Wave",
// 	    "owner": "Jane Smith",
// 	    "email": "jane.smith@example.com",
// 	    "address": "456 Elm Avenue, Town",
// 	    "company": "XYZ Bus Services",
// 	    "mobile_number": "987-654-3210"
// 	  }
// 	]
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		busList   []schema.Bus
		failedBus schema.FailedBus
	)

	// Unmarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), &busList)
	if err != nil {
		utility.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "Integration", Value: "Bus Ticketing – Bus"}, utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusInternalServerError(err)
	}

	for _, bus := range busList {
		// Checks whether the bus line exist or not
		busLineExist, err := validate.IsBusLineExisting(ctx, bus.Name, bus.Company)
		if err != nil {
			failedBus.SetFailedBus(bus, "failed to validate bus line if it exist")
			bus.Error(err, "IsBusLineExisting", "failed to validate bus line if it exist")

			continue
		}

		// If the bus line exists, continue to the next item
		if busLineExist {
			utility.Info("BusLineExisting", "already existing bus line", utility.KVP{Key: "bus", Value: bus})
			continue
		}

		// Set default values of the bus line information
		bus.SetValues()

		// Inserts a new bus line record to the DynamoDB
		err = query.CreateBusLine(ctx, bus)
		if err != nil {
			failedBus.SetFailedBus(bus, "failed to create a new bus line record")
			bus.Error(err, "DynamoDBError", "failed to create a new bus line record")

			continue
		}
	}

	if len(failedBus.Failed) > 0 {
		return api.Response(http.StatusBadRequest, failedBus), nil
	}

	return api.StatusOKWithoutBody()
}
