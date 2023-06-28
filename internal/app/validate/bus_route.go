package validate

import (
	"context"
	"errors"

	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	"github.com/rmarasigan/bus-ticketing/internal/app/query"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// IsBusRouteExisting checks if the DynamoDB Table is configured on the environment, and
// returns a boolean and error value to check whether the bus route alreadu exist or not.
func IsBusRouteExisting(ctx context.Context, routefilter schema.BusRouteFilter) (bool, error) {
	var tablename = env.BUS_ROUTE_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BUS_ROUTE_TABLE is not configured on the environment")
		err := errors.New("dynamodb BUS_ROUTE_TABLE environment variable is not set")

		return false, err
	}

	result, err := query.FilterBusRoute(ctx, routefilter)
	if err != nil {
		return false, err
	}

	return (len(result) > 0), nil
}
