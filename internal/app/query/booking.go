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

// getBooking returns the specific booking information.
func getBooking(ctx context.Context, tablename, id, busRouteId string) (schema.Bookings, error) {
	var booking schema.Bookings

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

// getBookingList returns all the booking information.
func getBookingList(ctx context.Context, tablename string) ([]schema.Bookings, error) {
	var bookings []schema.Bookings

	// Use the build expression to populate the DynamoDB Scan API
	var params = &dynamodb.ScanInput{TableName: aws.String(tablename)}

	result, err := awswrapper.DynamoDBScan(ctx, params)
	if err != nil {
		return nil, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual booking struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&bookings, result.Items)
		if err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

// GetBookingRecords checks if the DynamoDB Table is configured on the environment, and
// returns either a specific booking record or a list of booking records.
func GetBookingRecords(ctx context.Context, id, busRouteId string) ([]schema.Bookings, error) {
	var tablename = env.BOOKING_TABLE

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_TABLE environment variable is not set")

		return nil, err
	}

	// ********** Fetching a specific booking record ********** //
	if id != "" && busRouteId != "" {
		var bookings []schema.Bookings

		booking, err := getBooking(ctx, tablename, id, busRouteId)
		if err != nil {
			return nil, err
		}

		if booking == (schema.Bookings{}) {
			return bookings, nil
		}

		bookings = append(bookings, booking)
		return bookings, nil
	}

	// **************** List of booking records **************** //
	return getBookingList(ctx, tablename)
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

// FilterBookings checks if the DynamoDB Table is configured on the environment,
// fetches and returns a list of bookings information.
func FilterBookings(ctx context.Context, busId, routeId, status string) ([]schema.Bookings, error) {
	var (
		bookings   []schema.Bookings
		tablename  = env.BOOKING_TABLE
		filter     expression.ConditionBuilder
		filterList []expression.ConditionBuilder
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_TABLE environment variable is not set")

		return nil, err
	}

	// Check if the "bus_id" and "route_id" query parameters are not set and
	// if the "status" query parameter is set to fetch ALL records.
	if busId == "" && routeId == "" && status == "ALL" {
		return getBookingList(ctx, tablename)
	}

	if busId != "" {
		filterList = append(filterList, expression.Name("bus_id").Equal(expression.Value(busId)))
	}

	if routeId != "" {
		filterList = append(filterList, expression.Name("bus_route_id").Equal(expression.Value(routeId)))
	}

	if len(filterList) > 0 {
		filter = expression.And(expression.Name("status").Equal(expression.Value(status)), filterList[0], filterList[len(filterList)-1:]...)
	} else {
		filter = expression.Name("status").Equal(expression.Value(status))
	}

	result, err := FilterItems(ctx, tablename, filter)
	if err != nil {
		return bookings, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual bus route struct which the front-end can
		// understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&bookings, result.Items)
		if err != nil {
			return bookings, err
		}
	}

	return bookings, nil
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

// getCancelledBooking returns the specific cancelled booking information.
func getCancelledBooking(ctx context.Context, tablename, bookingId string) (schema.BookingCancelled, error) {
	var booking schema.BookingCancelled

	// Create a composite key expression
	key := expression.Key("booking_id").Equal(expression.Value(bookingId))

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

	// Unmarshal a map into actual Booking Cancelled struct which front-end
	// can understand as a JSON
	if result.Count > 0 {
		err := awswrapper.DynamoDBUnmarshalMap(&booking, result.Items[0])
		if err != nil {
			return booking, err
		}
	}

	return booking, nil
}

// GetCancelledBookingRecords checks if the DynamoDB Table is configured on the environment, and
// returns either a specific cancelled booking or a list of cancelled bookings.
func GetCancelledBookingRecords(ctx context.Context, bookingId string) ([]schema.BookingCancelled, error) {
	var (
		bookings  []schema.BookingCancelled
		tablename = env.BOOKING_CANCELLED_TABLE
	)

	// Check if the DynamoDB Table is configured
	if tablename == "" {
		trail.Error("dynamodb BOOKING_CANCELLED_TABLE is not configured on the environment")
		err := errors.New("dynamodb BOOKING_CANCELLED_TABLE environment is not set")

		return bookings, err
	}

	// ********** Fetching a specific cancelled booking record ********** //
	if bookingId != "" {
		booking, err := getCancelledBooking(ctx, tablename, bookingId)
		if err != nil {
			return nil, err
		}

		if booking == (schema.BookingCancelled{}) {
			return bookings, nil
		}

		bookings = append(bookings, booking)
		return bookings, nil
	}

	// **************** List of booking records **************** //
	// Use the build expression to populate the DynamoDB Scan API
	var params = &dynamodb.ScanInput{TableName: aws.String(tablename)}

	result, err := awswrapper.DynamoDBScan(ctx, params)
	if err != nil {
		return nil, err
	}

	if result.Count > 0 {
		// Unmarshal a map into actual Booking Cancelled struct which the
		// front-end can understand as a JSON.
		err = awswrapper.DynamoDBUnmarshalListOfMaps(&bookings, result.Items)
		if err != nil {
			return nil, err
		}
	}

	return bookings, nil
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
