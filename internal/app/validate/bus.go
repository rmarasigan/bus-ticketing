package validate

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

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

	result, err := query.IsExisting(ctx, tablename, key)
	if err != nil {
		return false, err
	}

	return result, nil
}
