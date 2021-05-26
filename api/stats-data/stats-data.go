package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// struct for the api gateway http response
type StatsDataResponse struct {
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
}

// struct for stats data
type StatsData struct {
	MostConnections     []map[string]interface{} `json:"most_connections"`
	MostActiveCities    []map[string]interface{} `json:"most_active_cities"`
	MostActiveCountries []map[string]interface{} `json:"most_active_countries"`
	MostIpAddresses     []map[string]interface{} `json:"most_ip_addresses"`
	MostIngressPorts    []map[string]interface{} `json:"most_ingress_ports"`
}

// type for unmarshalling aggregated logs from dynamodb
type AggregatedLogs []AggregatedLogEntry

// struct for aggregated log entries from dynamodb
type AggregatedLogEntry struct {
	Lat          float64  `json:"lat"`
	Lon          float64  `json:"lon"`
	Count        int      `json:"count"`
	IngressPorts []string `json:"ingress_ports"`
	IpAddresses  []string `json:"ip_addresses"`
	City         string   `json:"city"`
	Country      string   `json:"country"`
}

// type for storing three log entries which are the most of something
type AggregatedLogEntryPodium [3]AggregatedLogEntry

// packages stats data into a http response format for the api gateway
func getStatsDataResponse() (StatsDataResponse, error) {

	// http headers including cors
	headers := map[string]string{
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "OPTIONS,GET",
		"Content-Type":                 "application/json",
	}

	// function for creating an error response
	createErrorResponse := func(err error) StatsDataResponse {
		return StatsDataResponse{
			StatusCode:      400,
			Headers:         headers,
			Body:            "{\"error\":\"" + err.Error() + "\"}",
			IsBase64Encoded: false,
		}
	}

	statsData, err := createStatsData()
	if err != nil {
		return createErrorResponse(err), nil
	}

	body, err := json.Marshal(statsData)
	if err != nil {
		return createErrorResponse(err), nil
	}

	return StatsDataResponse{
		StatusCode:      200,
		Headers:         headers,
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}

// generates stats data
func createStatsData() (StatsData, error) {

	aggregatedLogs, err := getAggregatedLogs()
	if err != nil {
		return StatsData{}, err
	}

	statsData := StatsData{
		MostConnections: createMostConnectionsData(aggregatedLogs),
	}

	return statsData, nil
}

// gets aggregated logs from dynamodb
func getAggregatedLogs() (AggregatedLogs, error) {

	projection := expression.NamesList(expression.Name("lat"), expression.Name("lon"),
		expression.Name("count"), expression.Name("ingress_ports"), expression.Name("ip_addresses"),
		expression.Name("city"), expression.Name("country"))

	expression, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return AggregatedLogs{}, err
	}

	scanInput := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expression.Names(),
		ExpressionAttributeValues: expression.Values(),
		ProjectionExpression:      expression.Projection(),
		TableName:                 aws.String(aggregatedLogsTableName),
	}

	scanOutput, err := dynamoClient.Scan(scanInput)
	if err != nil {
		return AggregatedLogs{}, err
	}

	var aggregatedLogs AggregatedLogs
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &aggregatedLogs)
	if err != nil {
		return AggregatedLogs{}, err
	}

	return aggregatedLogs, nil
}

// generates most connections stats
func createMostConnectionsData(aggregatedLogs AggregatedLogs) []map[string]interface{} {
	var podium AggregatedLogEntryPodium

	// create podium of most connections
	for _, aggregatedLogEntry := range aggregatedLogs {
		if aggregatedLogEntry.Count > podium[2].Count {
			podium.insert(aggregatedLogEntry)
		}
	}

	// create stats data
	data := make([]map[string]interface{}, 3)
	for i, podiumEntry := range podium {
		data[i] = map[string]interface{}{
			"location": map[string]interface{}{
				"lat": podiumEntry.Lat,
				"lon": podiumEntry.Lon,
			},
			"data": map[string]interface{}{
				"connections":   podiumEntry.Count,
				"ingress_ports": podiumEntry.IngressPorts,
				"ip_addresses":  podiumEntry.IpAddresses,
			},
		}
	}

	return data
}

// inserts aggregated log entry into the top 3 based on count
func (podium *AggregatedLogEntryPodium) insert(entry AggregatedLogEntry) {
	if entry.Count > podium[0].Count {
		podium[0], podium[1], podium[2] = entry, podium[0], podium[1]
	} else if entry.Count > podium[1].Count {
		podium[1], podium[2] = entry, podium[1]
	} else if entry.Count > podium[2].Count {
		podium[2] = entry
	}
}
