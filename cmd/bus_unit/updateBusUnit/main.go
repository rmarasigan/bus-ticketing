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
// request query and body, updates the bus unit's record and responds with a 200
// OK HTTP Status.
//
// Endpoint:
//  https://{api_id}.execute-api.{region}.amazonaws.com/prod/bus_unit/update?code=xxxxx&bus_id=xxxxx
//
// Sample API Params:
//  bus_id=RLBSW-856996
// 	code=RLBSWV1_0606
//
// Sample API Payload:
// 	{
// 		"active": false
// 	}
//
// Sample API Response:
// 	{
// 		"bus_id": "RLBSW-856996",
// 		"code": "RLBSWV1_0606",
// 		"active": false,
// 		"min_capacity": 40,
// 		"max_capacity": 50,
// 		"date_created": "1686039674"
// 	}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var (
		unit        = new(schema.BusUnit)
		code_query  = request.QueryStringParameters["code"]
		busId_query = request.QueryStringParameters["bus_id"]
	)

	// Unmarshal the received JSON-encoded data
	err := utility.ParseJSON([]byte(request.Body), unit)
	if err != nil {
		unit.Error(err, "JSONError", "failed to unmarshal the JSON-encoded data",
			utility.KVP{Key: "payload", Value: request.Body})

		return api.StatusInternalServerError(err)
	}

	// Fetch the existing bus unit record/information
	busUnit, err := query.GetBusUnit(ctx, code_query, busId_query)
	if err != nil {
		unit.Error(err, "DynamoDBError", "failed to fetch the bus unit record")
		return api.StatusInternalServerError(err)
	}

	if busUnit == (schema.BusUnit{}) {
		err := errors.New("the bus unit you're trying to update is non-existent")
		unit.Error(err, "APIError", "the bus unit does not exist")

		return api.StatusBadRequest(err)
	}

	err = unit.ValidateMinimumCapacity()
	if err != nil {
		unit.Error(err, "APIError", "the minimum capacity is less than the required capacity (25)")
		return api.StatusBadRequest(err)
	}

	err = unit.ValidateMaximumCapacity(*busUnit.MinCapacity)
	if err != nil {
		unit.Error(err, "APIError", "the max capacity is less than the minimum capacity")
		return api.StatusBadRequest(err)
	}

	// Create a composite key that has both the partition/primary key
	// and the sort key of the item.
	var compositeKey = map[string]types.AttributeValue{
		"code":   &types.AttributeValueMemberS{Value: code_query},
		"bus_id": &types.AttributeValueMemberS{Value: busId_query},
	}

	// Construct the update builder
	busUnit = validate.UpdateBusUnitFields(*unit, busUnit)
	var update = expression.Set(expression.Name("active"), expression.Value(busUnit.Active)).
		Set(expression.Name("min_capacity"), expression.Value(busUnit.MinCapacity)).
		Set(expression.Name("max_capacity"), expression.Value(busUnit.MaxCapacity))

	result, err := query.UpdateBusUnit(ctx, compositeKey, update)
	if err != nil {
		busUnit.Error(err, "DynamoDBError", "failed to update the bus unit record")
		return api.StatusInternalServerError(err)
	}

	return api.StatusOK(result)
}
