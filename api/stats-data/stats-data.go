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
	MostConnections     []StatsDataPoint `json:"most_connections"`
	MostActiveCities    []StatsDataPoint `json:"most_active_cities"`
	MostActiveCountries []StatsDataPoint `json:"most_active_countries"`
	MostIpAddresses     []StatsDataPoint `json:"most_ip_addresses"`
	MostIngressPorts    []StatsDataPoint `json:"most_ingress_ports"`
}

// struct for stats data point
type StatsDataPoint struct {
	Value    interface{}             `json:"value"`
	MapData  []MapDataPoint          `json:"map_data"`
	Metadata map[string]MetadataItem `json:"metadata"`
}

// struct for map data point
type MapDataPoint struct {
	Lat      float64                 `json:"lat"`
	Lon      float64                 `json:"lon"`
	Metadata map[string]MetadataItem `json:"metadata"`
}

// struct for a metadata item
type MetadataItem struct {
	Title string      `json:"title"`
	Value interface{} `json:"value"`
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

// type for storing three log entries which have the most connections
type AggregatedLogEntryPodium [3]AggregatedLogEntry

// type for storing three stats data points which have the most connections
type StatsDataPointPodium [3]StatsDataPoint

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
		MostConnections:     createMostConnectionsData(aggregatedLogs),
		MostActiveCities:    createMostActiveCitiesData(aggregatedLogs),
		MostActiveCountries: createMostActiveCountriesData(aggregatedLogs),
		MostIpAddresses:     createMostIpAddressesData(aggregatedLogs),
		MostIngressPorts:    createMostIngressPortsData(aggregatedLogs),
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
func createMostConnectionsData(aggregatedLogs AggregatedLogs) []StatsDataPoint {
	var podium AggregatedLogEntryPodium

	// create podium of most connections
	for _, aggregatedLogEntry := range aggregatedLogs {
		if aggregatedLogEntry.Count > podium[2].Count {
			podium.insertByCount(aggregatedLogEntry)
		}
	}

	// create stats data
	data := make([]StatsDataPoint, 3)
	for i, podiumEntry := range podium {
		data[i] = StatsDataPoint{
			Value: podiumEntry.Count,
			MapData: []MapDataPoint{
				{
					Lat: podiumEntry.Lat,
					Lon: podiumEntry.Lon,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: podiumEntry.IngressPorts,
				},
				"ip_addresses": {
					Title: "IP addresses",
					Value: podiumEntry.IpAddresses,
				},
			},
		}
	}

	return data
}

// generates most active cities stats
func createMostActiveCitiesData(aggregatedLogs AggregatedLogs) []StatsDataPoint {

	// combine log entries by city
	combinedLogs := make(map[string]StatsDataPoint)
	for _, aggregatedLogEntry := range aggregatedLogs {

		if combinedLogEntry, ok := combinedLogs[aggregatedLogEntry.City]; ok {
			// update combined log entry
			combinedLogEntry.MapData = append(combinedLogEntry.MapData, MapDataPoint{
				Lat: aggregatedLogEntry.Lat,
				Lon: aggregatedLogEntry.Lon,
				Metadata: map[string]MetadataItem{
					"connections": {
						Title: "Connections",
						Value: aggregatedLogEntry.Count,
					},
				},
			})

			combinedLogEntry.Metadata["connections"] = MetadataItem{
				Title: combinedLogEntry.Metadata["connections"].Title,
				Value: combinedLogEntry.Metadata["connections"].Value.(int) + aggregatedLogEntry.Count,
			}

			combinedLogs[aggregatedLogEntry.City] = combinedLogEntry

		} else {
			// create new combined log entry
			combinedLogs[aggregatedLogEntry.City] = StatsDataPoint{
				Value: aggregatedLogEntry.City,
				MapData: []MapDataPoint{
					{
						Lat: aggregatedLogEntry.Lat,
						Lon: aggregatedLogEntry.Lon,
						Metadata: map[string]MetadataItem{
							"connections": {
								Title: "Connections",
								Value: aggregatedLogEntry.Count,
							},
						},
					},
				},
				Metadata: map[string]MetadataItem{
					"connections": {
						Title: "Connections",
						Value: aggregatedLogEntry.Count,
					},
				},
			}
		}
	}

	// create podium of most connections
	var podium StatsDataPointPodium
	for _, combinedLogEntry := range combinedLogs {
		podium.insertByMetadataConnections(combinedLogEntry)
	}

	// create stats data
	data := make([]StatsDataPoint, 3)
	for i, podiumEntry := range podium {
		data[i] = podiumEntry
	}

	return data
}

// generates most active countries stats
func createMostActiveCountriesData(aggregatedLogs AggregatedLogs) []StatsDataPoint {

	// combine log entries by country
	combinedLogs := make(map[string]StatsDataPoint)
	for _, aggregatedLogEntry := range aggregatedLogs {

		if combinedLogEntry, ok := combinedLogs[aggregatedLogEntry.Country]; ok {
			// update combined log entry
			combinedLogEntry.MapData = append(combinedLogEntry.MapData, MapDataPoint{
				Lat: aggregatedLogEntry.Lat,
				Lon: aggregatedLogEntry.Lon,
				Metadata: map[string]MetadataItem{
					"connections": {
						Title: "Connections",
						Value: aggregatedLogEntry.Count,
					},
				},
			})

			combinedLogEntry.Metadata["connections"] = MetadataItem{
				Title: combinedLogEntry.Metadata["connections"].Title,
				Value: combinedLogEntry.Metadata["connections"].Value.(int) + aggregatedLogEntry.Count,
			}

			combinedLogs[aggregatedLogEntry.Country] = combinedLogEntry

		} else {
			// create new combined log entry
			combinedLogs[aggregatedLogEntry.Country] = StatsDataPoint{
				Value: aggregatedLogEntry.Country,
				MapData: []MapDataPoint{
					{
						Lat: aggregatedLogEntry.Lat,
						Lon: aggregatedLogEntry.Lon,
						Metadata: map[string]MetadataItem{
							"connections": {
								Title: "Connections",
								Value: aggregatedLogEntry.Count,
							},
						},
					},
				},
				Metadata: map[string]MetadataItem{
					"connections": {
						Title: "Connections",
						Value: aggregatedLogEntry.Count,
					},
				},
			}
		}
	}

	// create podium of most connections
	var podium StatsDataPointPodium
	for _, combinedLogEntry := range combinedLogs {
		podium.insertByMetadataConnections(combinedLogEntry)
	}

	// create stats data
	data := make([]StatsDataPoint, 3)
	for i, podiumEntry := range podium {
		data[i] = podiumEntry
	}

	return data
}

// generates most ip addresses stats
func createMostIpAddressesData(aggregatedLogs AggregatedLogs) []StatsDataPoint {
	var podium AggregatedLogEntryPodium

	// create podium of most ip addresses
	for _, aggregatedLogEntry := range aggregatedLogs {
		if len(aggregatedLogEntry.IpAddresses) > len(podium[2].IpAddresses) {
			podium.insertByIpAddresses(aggregatedLogEntry)
		}
	}

	// create stats data
	data := make([]StatsDataPoint, 3)
	for i, podiumEntry := range podium {
		data[i] = StatsDataPoint{
			Value: len(podiumEntry.IpAddresses),
			MapData: []MapDataPoint{
				{
					Lat: podiumEntry.Lat,
					Lon: podiumEntry.Lon,
				},
			},
			Metadata: map[string]MetadataItem{
				"ip_addresses": {
					Title: "IP addresses",
					Value: podiumEntry.IpAddresses,
				},
			},
		}
	}

	return data
}

// generates most ip addresses stats
func createMostIngressPortsData(aggregatedLogs AggregatedLogs) []StatsDataPoint {
	var podium AggregatedLogEntryPodium

	// create podium of most ip addresses
	for _, aggregatedLogEntry := range aggregatedLogs {
		if len(aggregatedLogEntry.IngressPorts) > len(podium[2].IngressPorts) {
			podium.insertByIngressPorts(aggregatedLogEntry)
		}
	}

	// create stats data
	data := make([]StatsDataPoint, 3)
	for i, podiumEntry := range podium {
		data[i] = StatsDataPoint{
			Value: len(podiumEntry.IngressPorts),
			MapData: []MapDataPoint{
				{
					Lat: podiumEntry.Lat,
					Lon: podiumEntry.Lon,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: podiumEntry.IngressPorts,
				},
			},
		}
	}

	return data
}

// inserts aggregated log entry into the top 3 based on count
func (podium *AggregatedLogEntryPodium) insertByCount(entry AggregatedLogEntry) {
	if entry.Count > podium[0].Count {
		podium[0], podium[1], podium[2] = entry, podium[0], podium[1]
	} else if entry.Count > podium[1].Count {
		podium[1], podium[2] = entry, podium[1]
	} else if entry.Count > podium[2].Count {
		podium[2] = entry
	}
}

// inserts aggregated log entry into the top 3 based on number of ip addresses
func (podium *AggregatedLogEntryPodium) insertByIpAddresses(entry AggregatedLogEntry) {
	if len(entry.IpAddresses) > len(podium[0].IpAddresses) {
		podium[0], podium[1], podium[2] = entry, podium[0], podium[1]
	} else if len(entry.IpAddresses) > len(podium[1].IpAddresses) {
		podium[1], podium[2] = entry, podium[1]
	} else if len(entry.IpAddresses) > len(podium[2].IpAddresses) {
		podium[2] = entry
	}
}

// inserts aggregated log entry into the top 3 based on number of ingress ports
func (podium *AggregatedLogEntryPodium) insertByIngressPorts(entry AggregatedLogEntry) {
	if len(entry.IngressPorts) > len(podium[0].IngressPorts) {
		podium[0], podium[1], podium[2] = entry, podium[0], podium[1]
	} else if len(entry.IngressPorts) > len(podium[1].IngressPorts) {
		podium[1], podium[2] = entry, podium[1]
	} else if len(entry.IngressPorts) > len(podium[2].IngressPorts) {
		podium[2] = entry
	}
}

// inserts combined log entry into the top 3 based on connections
func (podium *StatsDataPointPodium) insertByMetadataConnections(entry StatsDataPoint) {
	if podium[0].Metadata["connections"].Value == nil ||
		entry.Metadata["connections"].Value.(int) > podium[0].Metadata["connections"].Value.(int) {
		podium[0], podium[1], podium[2] = entry, podium[0], podium[1]
	} else if podium[1].Metadata["connections"].Value == nil ||
		entry.Metadata["connections"].Value.(int) > podium[1].Metadata["connections"].Value.(int) {
		podium[1], podium[2] = entry, podium[1]
	} else if podium[2].Metadata["connections"].Value == nil ||
		entry.Metadata["connections"].Value.(int) > podium[2].Metadata["connections"].Value.(int) {
		podium[2] = entry
	}
}
