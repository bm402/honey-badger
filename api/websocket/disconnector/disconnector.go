package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// struct for the api gateway http response
type DisconnectorResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

type DisconnectionItem struct {
	ConnectionID string `json:"connection_id"`
}

// packages connector response into a http response format for the api gateway
func handleDisconnectionWithResponse(disconnectionEvent events.APIGatewayWebsocketProxyRequest) (DisconnectorResponse, error) {

	// http headers
	headers := map[string]string{
		/*"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "OPTIONS,GET",
		"Content-Type": "application/json",*/
	}

	err := addConnectionToConnectionsTable(disconnectionEvent.RequestContext.ConnectionID)
	if err != nil {
		return DisconnectorResponse{
			StatusCode:      400,
			Headers:         headers,
			Body:            "Failed to disconnect: " + err.Error(),
			IsBase64Encoded: false,
		}, nil
	}

	return DisconnectorResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            "Disconnected",
		IsBase64Encoded: false,
	}, nil
}

// puts new connection id into the connections table
func addConnectionToConnectionsTable(connectionId string) error {
	disconnectionItem := DisconnectionItem{
		ConnectionID: connectionId,
	}

	keyValues, err := dynamodbattribute.MarshalMap(disconnectionItem)
	if err != nil {
		return err
	}

	deleteItemInput := &dynamodb.DeleteItemInput{
		Key:       keyValues,
		TableName: aws.String(connectionsTableName),
	}

	_, err = dynamoClient.DeleteItem(deleteItemInput)
	return err
}
