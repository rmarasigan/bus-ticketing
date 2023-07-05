package query

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/rmarasigan/bus-ticketing/api/schema"
	"github.com/rmarasigan/bus-ticketing/internal/app/env"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// GetBooking checks if the DynamoDB Table is configured on the environment, and
// fetch and returns the booking information.
func GetBooking(ctx context.Context, id, busRouteId string) (schema.Bookings, error) {
	var (
		booking   schema.Bookings
		tablename = env.BOOKING_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_TABLE environment variable is not set")

		return booking, err
	}

	// Create a composite key expression
	key := expression.KeyAnd(expression.Key("id").Equal(expression.Value(id)),
		expression.Key("bus_route_id").Equal(expression.Value(busRouteId)))

	// Build an expression to retrieve the item from the DynamoDB
	expr, err := expression.NewBuilder().WithKeyCondition(key).Build()
	if err != nil {
		return booking, err
	}

	// Build the query params
	params := &dynamodb.QueryInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := awswrapper.DynamoDBQuery(ctx, params)
	if err != nil {
		return booking, err
	}

	// Unmarshal a map into actual Booking struct which front-end can
	// understand as a JSON
	if result.Count > 0 {
		err := awswrapper.DynamoDBUnmarshalMap(&booking, result.Items[0])
		if err != nil {
			return booking, err
		}
	}

	return booking, nil
}

// CreateBooking checks if the DynamoDB Table is configured on the environment, and
// creates a new booking record.
func CreateBooking(ctx context.Context, data interface{}) error {
	var tablename = env.BOOKING_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_TABLE environment is not set")

		return err
	}

	// Save the Booking record into the DynamoDB Table
	err := InsertItem(ctx, tablename, data)
	if err != nil {
		trail.Error("failed to insert a new booking record")
		return err
	}

	return nil
}

// UpdateBooking checks if the DynamoDB Table is configured on the environment and
// updates the booking record.
func UpdateBooking(ctx context.Context, key map[string]types.AttributeValue, update expression.UpdateBuilder) (schema.Bookings, error) {
	var (
		booking   schema.Bookings
		tablename = env.BOOKING_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_TABLE environment variable is not set")

		return booking, err
	}

	result, err := UpdateItem(ctx, tablename, key, update)
	if err != nil {
		trail.Error("failed to update the booking record")
		return booking, err
	}

	// Unmarshal a map into actual booking struct which the front-end can
	// understand as a JSON.
	err = awswrapper.DynamoDBUnmarshalMap(&booking, result.Attributes)
	if err != nil {
		return booking, err
	}

	return booking, nil
}

// RecordBookingCancelled updates the existing item's attribute or adds a new item
// to the table if it does not exist.
func RecordBookingCancelled(ctx context.Context, key map[string]types.AttributeValue, update expression.UpdateBuilder) (schema.BookingCancelled, error) {
	var (
		booking   schema.BookingCancelled
		tablename = env.BOOKING_CANCELLED_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_CANCELLED_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_CANCELLED_TABLE environment is not set")

		return booking, err
	}

	result, err := UpdateItem(ctx, tablename, key, update)
	if err != nil {
		trail.Error("failed to update the booking record")
		return booking, err
	}

	// Unmarshal a map into actual booking struct which the front-end can
	// understand as a JSON.
	err = awswrapper.DynamoDBUnmarshalMap(&booking, result.Attributes)
	if err != nil {
		return booking, err
	}

	return booking, nil
}
