package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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

// struct for aggregated log entry
type AggregatedLogEntry struct {
	LatLon       string   `json:"lat_lon"`
	Lat          float64  `json:"lat"`
	Lon          float64  `json:"lon"`
	City         string   `json:"city"`
	Country      string   `json:"country"`
	IpAddresses  []string `json:"ip_addresses"`
	IngressPorts []string `json:"ingress_ports"`
	Inputs       []string `json:"inputs"`
	Count        int      `json:"count"`
}

func aggregateDynamoDBEvent(event events.DynamoDBEvent) error {
	for _, eventRecord := range event.Records {

		// get raw log entry data from event
		rawLogEntry, err := unmarshalStreamImageToRawLogEntry(eventRecord.Change.NewImage)
		if err != nil {
			return err
		}

		// query aggregate table for existing record for lat/lon
		aggregatedLogEntry, exists, err := getAggregatedLogEntryFromRawLogEntryData(rawLogEntry)
		if err != nil {
			return err
		}

		if exists {
			// if aggregate record exists update with new data
			aggregatesToUpdate := findAggregatesToUpdate(aggregatedLogEntry, rawLogEntry)

			err := writeUpdatedAggregatedLogEntry(aggregatedLogEntry.LatLon, aggregatesToUpdate)
			if err != nil {
				return err
			}

		} else {
			// if aggregate record does not exist write new record
			aggregatedLogEntry = AggregatedLogEntry{
				LatLon:       fmt.Sprintf("%f,%f", rawLogEntry.Lat, rawLogEntry.Lon),
				Lat:          rawLogEntry.Lat,
				Lon:          rawLogEntry.Lon,
				City:         rawLogEntry.City,
				Country:      rawLogEntry.Country,
				IpAddresses:  []string{rawLogEntry.IpAddress},
				IngressPorts: []string{rawLogEntry.IngressPort},
				Inputs:       []string{rawLogEntry.Input},
				Count:        1,
			}

			err := writeNewAggregatedLogEntry(aggregatedLogEntry)
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

// retrieves aggregated log entry (if exists) using lat/lon from raw log entry
func getAggregatedLogEntryFromRawLogEntryData(rawLogEntry RawLogEntry) (AggregatedLogEntry, bool, error) {
	var aggregatedLogEntry AggregatedLogEntry
	latLon := fmt.Sprintf("%f,%f", rawLogEntry.Lat, rawLogEntry.Lon)

	keyAttributes := map[string]string{
		"lat_lon": latLon,
	}
	key, err := dynamodbattribute.MarshalMap(keyAttributes)
	if err != nil {
		return aggregatedLogEntry, false, err
	}

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(aggregatedLogsTableName),
		Key:       key,
	}

	getItemOutput, err := dynamoClient.GetItem(getItemInput)
	if err != nil {
		return aggregatedLogEntry, false, err
	}

	if len(getItemOutput.Item) == 0 {
		return aggregatedLogEntry, false, nil
	}

	err = dynamodbattribute.UnmarshalMap(getItemOutput.Item, &aggregatedLogEntry)
	if err != nil {
		return aggregatedLogEntry, false, err
	}

	return aggregatedLogEntry, true, nil
}

// writes new aggregated log entry to the db
func writeNewAggregatedLogEntry(aggregatedLogEntry AggregatedLogEntry) error {
	item, err := dynamodbattribute.MarshalMap(aggregatedLogEntry)
	if err != nil {
		return err
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(aggregatedLogsTableName),
		Item:      item,
	}

	_, err = dynamoClient.PutItem(putItemInput)
	return err
}

// writes updated aggregated log entry to the db
func writeUpdatedAggregatedLogEntry(aggregatedLogEntryKey string, aggregatesToUpdate map[string][]string) error {
	keyAttributes := map[string]string{
		"lat_lon": aggregatedLogEntryKey,
	}
	key, err := dynamodbattribute.MarshalMap(keyAttributes)
	if err != nil {
		return err
	}

	// increments aggregate count
	update := expression.Set(expression.Name("count"), expression.Plus(expression.Name("count"), expression.Value(1)))
	for name, value := range aggregatesToUpdate {
		// appends new raw values to aggregate arrays
		update.Set(expression.Name(name), expression.ListAppend(expression.Name(name), expression.Value(value)))
	}

	expression, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(aggregatedLogsTableName),
		Key:                       key,
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		UpdateExpression:          expression.Update(),
	}

	_, err = dynamoClient.UpdateItem(updateItemInput)
	return err
}

// finds which values from raw log entry should be updated in the aggregated log entry
func findAggregatesToUpdate(aggregatedLogEntry AggregatedLogEntry, rawLogEntry RawLogEntry) map[string][]string {

	// checks whether an array contains a string
	contains := func(arr []string, str string) bool {
		for _, v := range arr {
			if v == str {
				return true
			}
		}
		return false
	}

	aggregatesToUpdate := make(map[string][]string)
	if !contains(aggregatedLogEntry.IpAddresses, rawLogEntry.IpAddress) {
		aggregatesToUpdate["ip_addresses"] = []string{rawLogEntry.IpAddress}
	}
	if !contains(aggregatedLogEntry.IngressPorts, rawLogEntry.IngressPort) {
		aggregatesToUpdate["ingress_ports"] = []string{rawLogEntry.IngressPort}
	}
	if !contains(aggregatedLogEntry.Inputs, rawLogEntry.Input) {
		aggregatesToUpdate["inputs"] = []string{rawLogEntry.Input}
	}

	return aggregatesToUpdate
}
