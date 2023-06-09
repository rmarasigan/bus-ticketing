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

// UpdateBusUnitFields validates if the field that are going to be updated
// are empty or not to set its previous value.
//
// Fields that are valdiated:
//  active, min_capacity, max_capacity
func UpdateBusUnitFields(unit, old schema.BusUnit) schema.BusUnit {
	if unit.Active == nil {
		unit.Active = old.Active
	}

	if unit.MinCapacity == 0 {
		unit.MinCapacity = old.MinCapacity
	}

	if unit.MaxCapacity == 0 || unit.MaxCapacity < unit.MinCapacity {
		unit.MaxCapacity = old.MaxCapacity
	}

	return unit
}

// IsBusUnitExisting checks if the DynamoDB Table is configured on the environment, and
// returns a boolean and error value to check whether the bus unit already exist or not.
func IsBusUnitExisting(ctx context.Context, busId, code string) (bool, error) {
	var tablename = env.BUS_UNIT

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_UNIT_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_UNIT_TABLE environment variable is not set")

		return false, err
	}

	// Create a composite key expression
	key := expression.KeyAnd(expression.Key("code").Equal(expression.Value(code)), expression.Key("bus_id").Equal(expression.Value(busId)))

	result, err := query.IsExisting(ctx, tablename, key)
	if err != nil {
		return false, err
	}

	return result, nil
}
