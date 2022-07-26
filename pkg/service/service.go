package service

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	cw "github.com/rmarasigan/bus-ticketing/pkg/cw/logger"
)

// DynamoDBAttributeResponse returns an item with the values of result object.
func DynamoDBAttributeResponse(v interface{}, result map[string]*dynamodb.AttributeValue) error {
	// Unmarshal it into actual interface which front-end can understand as a JSON
	err := dynamodbattribute.UnmarshalMap(result, v)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBUnmarshalMap", Message: "Failed to unmarshal result map to interface"})
		return err
	}

	return nil
}

// DynamoDBAttributeResponse returns a list of item with the values of result object.
func DynamoDBAttributesResponse(v interface{}, result []map[string]*dynamodb.AttributeValue) error {
	// Unmarshal it into actual interface which front-end can understand as a JSON
	err := dynamodbattribute.UnmarshalListOfMaps(result, v)
	if err != nil {
		cw.Error(err, &cw.Logs{Code: "DynamoDBUnmarshalListOfMaps", Message: "Failed to unmarshal result map to interface"})
		return err
	}

	return nil
}
