package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// struct for raw log entry
type RawLogEntry struct {
	IngressPort string  `json:"ingress_port"`
	Timestamp   int64   `json:"timestamp"`
	IpAddress   string  `json:"ip_address"`
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Input       string  `json:"input"`
}

// struct for connection item from connections table
type ConnectionItem struct {
	ConnectionID string `json:"connection_id"`
}

func broadcastDynamoDBEvent(event events.DynamoDBEvent) error {
	for _, eventRecord := range event.Records {

		// get raw log entry data from event
		rawLogEntry, err := unmarshalStreamImageToRawLogEntry(eventRecord.Change.NewImage)
		if err != nil {
			return err
		}

		// get websocket connection ids
		connectionIds, err := getConnectionIds()
		if err != nil {
			return err
		}

		// marshal websocket data to broadcast
		broadcastData, err := json.Marshal(rawLogEntry)
		if err != nil {
			return err
		}

		// send data to websocket connections
		for _, connectionId := range connectionIds {
			connectionInput := &apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: aws.String(connectionId),
				Data:         broadcastData,
			}

			_, err := apigwClient.PostToConnection(connectionInput)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// converts stream image to raw log entry struct
func unmarshalStreamImageToRawLogEntry(attribute map[string]events.DynamoDBAttributeValue) (RawLogEntry, error) {
	var rawLogEntry RawLogEntry
	dbAttributes := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attribute {
		var dbAttribute dynamodb.AttributeValue

		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			return rawLogEntry, marshalErr
		}

		json.Unmarshal(bytes, &dbAttribute)
		dbAttributes[k] = &dbAttribute
	}

	err := dynamodbattribute.UnmarshalMap(dbAttributes, &rawLogEntry)
	return rawLogEntry, err
}

// gets websocket connection ids from connections table
func getConnectionIds() ([]string, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String(connectionsTableName),
	}

	scanOutput, err := dynamoClient.Scan(scanInput)
	if err != nil {
		return []string{}, err
	}

	var connectionItems []ConnectionItem
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &connectionItems)
	if err != nil {
		return []string{}, err
	}

	connectionIds := []string{}
	for _, connectionItem := range connectionItems {
		connectionIds = append(connectionIds, connectionItem.ConnectionID)
	}

	return connectionIds, nil
}
