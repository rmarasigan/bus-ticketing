package main

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
// request query and body, updates the bus line's information/record and responds
// with a 200 OK HTTP Status.
//
// Method: POST
//
// Endpoint: https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus/update?id=xxxxx&name=xxxxx
//
// Sample API Params:
//  id=BCBSCMPN-875011
//  name=Blue Horizon
//
// Sample API Payload:
// 	{
// 		"address": "Långbro, Stockholm",
// 		"mobile_number": "0567-8809105"
// 	}
//
// Sample API Response:
// 	{
// 	  "id": "BCBSCMPN-875011",
// 	  "name": "Blue Horizon",
// 	  "owner": "Daniel Martinez",
// 	  "email": "daniel.martinez@example.com",
// 	  "address": "Långbro, Stockholm",
// 	  "company": "ABC Bus Company",
// 	  "mobile_number": "0567-8809105",
// 	  "date_created": "1687501112"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		bus        = new(schema.Bus)
		id_query   = request.QueryStringParameters["id"]
		name_query = request.QueryStringParameters["name"]
	)

	err := bus.IsEmptyPayload(request.Body)
	if err != nil {
		return api.StatusBadRequest(err)
	}

	// Unmarshal the received JSON-encoded data
	err = utility.ParseJSON([]byte(request.Body), bus)
	if err != nil {
		bus.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusInternalServerError(err)
	}

	// Fetch the existing bus line record
	busLines, err := query.GetBusLineRecords(ctx, id_query, name_query)
	if err != nil {
		bus.Error(err, "DynamoDBError", "failed to fetch the bus line record")
		return api.StatusInternalServerError(err)
	}

	busLine := busLines[0]
	if busLine == (schema.Bus{}) {
		err := errors.New("the bus line record you're trying to update is non-existent")
		bus.Error(err, "APIError", "the bus line does not exist")

		return api.StatusBadRequest(err)
	}

	// Create a composite key that has both the partition/primary key
	// and the sort key of the item.
	var compositeKey = map[string]types.AttributeValue{
		"name":    &types.AttributeValueMemberS{Value: name_query},
		"company": &types.AttributeValueMemberS{Value: busLine.Company},
	}

	// Construct the update builder
	busLine = validate.UpdateBusLineFields(*bus, busLine)
	var update = expression.Set(expression.Name("owner"), expression.Value(busLine.Owner)).
		Set(expression.Name("email"), expression.Value(busLine.Email)).
		Set(expression.Name("address"), expression.Value(busLine.Address)).
		Set(expression.Name("mobile_number"), expression.Value(busLine.MobileNumber))

	result, err := query.UpdateBusLine(ctx, compositeKey, update)
	if err != nil {
		busLine.Error(err, "DynamoDBError", "failed to update the bus line record")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(result)
}
