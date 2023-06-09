package query

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	awswrapper "github.com/rmarasigan/bus-ticketing/internal/aws_wrapper"
	"github.com/rmarasigan/bus-ticketing/internal/trail"
)

// InsertItem converts the data into a map of AttributeValues and performs DynamoDB Put
// Item Operation to create the new item into the DynamoDB Table.
func InsertItem(ctx context.Context, tablename string, data interface{}) error {
	// Marshal the user to a map of AttributeValues
	values, err := awswrapper.DynamoDBMarshalMap(data)
	if err != nil {
		trail.Error("failed to marshal data to a map of AttributeValues")
		return err
	}

	params := &dynamodb.PutItemInput{
		Item:      values,
		TableName: aws.String(tablename),
	}

	// Save the item into the DynamoDB Table
	_, err = awswrapper.DynamoDBPutItem(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

// IsExisting creates an expression, performs DyanmoDB Query Operation, and returns
// if the item from the DynamoDB Table exist or not.
func IsExisting(ctx context.Context, tablename string, key expression.KeyConditionBuilder) (bool, error) {
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

// UpdateItem creates an expression with UpdateBuilder, performs the DynamoDB UpdateItem
// operation, and returns all of the attributes of the item.
func UpdateItem(ctx context.Context, tablename string, key map[string]types.AttributeValue, update expression.UpdateBuilder) (*dynamodb.UpdateItemOutput, error) {
	// Using the update expression to create a DynamoDB Expression
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		trail.Error("failed to build DynamoDB Expression")
		return nil, err
	}

	// Use the build expression to populate the DynamoDB Update Item API
	var params = &dynamodb.UpdateItemInput{
		Key:                       key,
		TableName:                 aws.String(tablename),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
	}

	result, err := awswrapper.DynamoDBUpdateItem(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FilterItems creates an expression with ConditionBuilder, performs the DynamoDB Scan
// operation, and returns a list of attributes of the items.
func FilterItems(ctx context.Context, tablename string, filter expression.ConditionBuilder) (*dynamodb.ScanOutput, error) {
	// Using the update expression to create a DynamoDB Expression
	expr, err := expression.NewBuilder().WithCondition(filter).Build()
	if err != nil {
		trail.Error("failed to build DynamoDB Expression")
		return nil, err
	}

	// Use the build expression to populate the DynamoDB Scan API
	var params = &dynamodb.ScanInput{
		TableName:                 aws.String(tablename),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Condition(),
	}

	result, err := awswrapper.DynamoDBScan(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// BuildMultipleORConditionExpression builds the logical OR clause of the argument ConditionBuilders.
func BuildMultipleORConditionExpression(conditions []expression.ConditionBuilder) expression.ConditionBuilder {
	var filter expression.ConditionBuilder

	if len(conditions) == 0 {
		return filter
	}

	filter = filter.Or(conditions[0])

	for _, condition := range conditions[1:] {
		filter = filter.Or(conditions[0], condition)
	}

	return filter
}
