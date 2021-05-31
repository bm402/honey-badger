package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// struct for the api gateway http response
type ConnectorResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type ConnectionItem struct {
	ConnectionID string `json:"connection_id"`
}

// packages connector response into a http response format for the api gateway
func handleNewConnectionWithResponse(connectionEvent events.APIGatewayWebsocketProxyRequest) (ConnectorResponse, error) {

	// http headers
	headers := map[string]string{
		/*"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "OPTIONS,GET",
		"Content-Type": "application/json",*/
	}

	err := addConnectionToConnectionsTable(connectionEvent.RequestContext.ConnectionID)
	if err != nil {
		return ConnectorResponse{
			StatusCode:      400,
			Headers:         headers,
			Body:            "Failed to connect: " + err.Error(),
			IsBase64Encoded: false,
		}, nil
	}

	return ConnectorResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            "Connected",
		IsBase64Encoded: false,
	}, nil
}

// puts new connection id into the connections table
func addConnectionToConnectionsTable(connectionId string) error {
	connectionItem := ConnectionItem{
		ConnectionID: connectionId,
	}

	attributeValues, err := dynamodbattribute.MarshalMap(connectionItem)
	if err != nil {
		return err
	}

	putItemInput := &dynamodb.PutItemInput{
		Item:      attributeValues,
		TableName: aws.String(connectionsTableName),
	}

	_, err = dynamoClient.PutItem(putItemInput)
	return err
}
