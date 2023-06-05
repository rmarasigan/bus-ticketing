package validate

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// CreateBusLine validates if the required fields are empty or not.
//
// Fields that are validated:
//  name, owner, email, company, mobile_number
func CreateBusLine(bus schema.Bus) error {
	var fields []string

	if bus.Name == "" {
		fields = append(fields, "name")
	}

	if bus.Owner == "" {
		fields = append(fields, "owner")
	}

	if bus.Email == "" {
		fields = append(fields, "email")
	}

	if bus.Company == "" {
		fields = append(fields, "company")
	}

	if bus.MobileNumber == "" {
		fields = append(fields, "mobile_number")
	}

	if len(fields) > 0 {
		return fmt.Errorf("missing %s field(s)", strings.Join(fields, ", "))
	}

	return nil
}

// UpdateBusLineFields validates if the field that are going to be updated
// are empty or not to set its previous value.
//
// Fields that are validated:
//  owner, email, address, mobile_number
func UpdateBusLineFields(bus schema.Bus, old schema.Bus) schema.Bus {
	if bus.Owner == "" {
		bus.Owner = old.Owner
	}

	if bus.Email == "" {
		bus.Email = old.Email
	}

	if bus.Address == "" {
		bus.Address = old.Address
	}

	if bus.MobileNumber == "" {
		bus.MobileNumber = old.MobileNumber
	}

	return bus
}

// IsBusLineExisting checks if the DynamoDB Table is configured on the environment, and
// returns a boolean and error value to check whether the bus line already exist or not.
func IsBusLineExisting(ctx context.Context, name, company string) (bool, error) {
	var tablename = env.BUS_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_TABLE environment variable is not set")

		return false, err
	}

	// Create a composite key expression
	key := expression.KeyAnd(expression.Key("name").Equal(expression.Value(name)), expression.Key("company").Equal(expression.Value(company)))

	// Build an expression to retrieve item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		return false, err
	}

	// Build the query params parameter
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return false, err
	}

	return (result.Count > 0), nil
}
