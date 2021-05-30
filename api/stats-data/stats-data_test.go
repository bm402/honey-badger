package main

import (
	"testing"
)

func TestAggregatedLogEntryPodiumInsertFirst(t *testing.T) {
	podium := AggregatedLogEntryPodium{
		AggregatedLogEntry{Count: 3},
		AggregatedLogEntry{Count: 2},
		AggregatedLogEntry{Count: 1},
	}

	want := AggregatedLogEntryPodium{
		AggregatedLogEntry{Count: 4},
		AggregatedLogEntry{Count: 3},
		AggregatedLogEntry{Count: 2},
	}

	podium.insertByCount(AggregatedLogEntry{Count: 4})

	if podium[0].Count != 4 || podium[1].Count != 3 || podium[2].Count != 2 {
		t.Error("Aggregated log entry podium insert was incorrect, got: ", podium, ", want: ", want)
	}
}

func TestAggregatedLogEntryPodiumInsertSecond(t *testing.T) {
	podium := AggregatedLogEntryPodium{
		AggregatedLogEntry{Count: 4},
		AggregatedLogEntry{Count: 2},
		AggregatedLogEntry{Count: 1},
	}

	want := AggregatedLogEntryPodium{
		AggregatedLogEntry{Count: 4},
		AggregatedLogEntry{Count: 3},
		AggregatedLogEntry{Count: 2},
	}

	podium.insertByCount(AggregatedLogEntry{Count: 3})

	if podium[0].Count != 4 || podium[1].Count != 3 || podium[2].Count != 2 {
		t.Error("Aggregated log entry podium insert was incorrect, got: ", podium, ", want: ", want)
	}
}

func TestAggregatedLogEntryPodiumInsertThird(t *testing.T) {
	podium := AggregatedLogEntryPodium{
		AggregatedLogEntry{Count: 4},
		AggregatedLogEntry{Count: 3},
		AggregatedLogEntry{Count: 1},
	}

	want := AggregatedLogEntryPodium{
		AggregatedLogEntry{Count: 4},
		AggregatedLogEntry{Count: 3},
		AggregatedLogEntry{Count: 2},
	}

	podium.insertByCount(AggregatedLogEntry{Count: 2})

	if podium[0].Count != 4 || podium[1].Count != 3 || podium[2].Count != 2 {
		t.Error("Aggregated log entry podium insert was incorrect, got: ", podium, ", want: ", want)
	}
}

func TestCreateMostConnectionsData(t *testing.T) {
	aggregatedLogs := AggregatedLogs{
		AggregatedLogEntry{Lat: 0, Lon: 0, Count: 3, IngressPorts: []string{"0"}, IpAddresses: []string{"0.0.0.0"}},
		AggregatedLogEntry{Lat: 1, Lon: 1, Count: 4, IngressPorts: []string{"1"}, IpAddresses: []string{"1.1.1.1"}},
		AggregatedLogEntry{Lat: 2, Lon: 2, Count: 1, IngressPorts: []string{"2"}, IpAddresses: []string{"2.2.2.2"}},
		AggregatedLogEntry{Lat: 3, Lon: 3, Count: 2, IngressPorts: []string{"3"}, IpAddresses: []string{"3.3.3.3"}},
	}

	want := []StatsDataPoint{
		{
			Value: 4,
			MapData: []MapDataPoint{
				{
					Lat: 1,
					Lon: 1,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: []string{"1"},
				},
				"ip_addresses": {
					Title: "IP addresses",
					Value: []string{"1.1.1.1"},
				},
			},
		},
		{
			Value: 3,
			MapData: []MapDataPoint{
				{
					Lat: 0,
					Lon: 0,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: []string{"0"},
				},
				"ip_addresses": {
					Title: "IP addresses",
					Value: []string{"0.0.0.0"},
				},
			},
		},
		{
			Value: 2,
			MapData: []MapDataPoint{
				{
					Lat: 3,
					Lon: 3,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: []string{"3"},
				},
				"ip_addresses": {
					Title: "IP addresses",
					Value: []string{"3.3.3.3"},
				},
			},
		},
	}

	got := createMostConnectionsData(aggregatedLogs)

	if got[0].Value.(int) != want[0].Value.(int) ||
		got[1].Value.(int) != want[1].Value.(int) ||
		got[2].Value.(int) != want[2].Value.(int) {

		t.Error("Most connections data was incorrect, got: ", got, ", want: ", want)
	}
}

func TestCreateMostActiveCitiesData(t *testing.T) {
	aggregatedLogs := AggregatedLogs{
		AggregatedLogEntry{Lat: 0, Lon: 0, Count: 1, City: "Birmingham"},
		AggregatedLogEntry{Lat: 1, Lon: 1, Count: 8, City: "Leeds"},
		AggregatedLogEntry{Lat: 2, Lon: 2, Count: 4, City: "London"},
		AggregatedLogEntry{Lat: 3, Lon: 3, Count: 6, City: "Bristol"},
		AggregatedLogEntry{Lat: 4, Lon: 4, Count: 5, City: "London"},
	}

	want := []StatsDataPoint{
		{
			Value: "London",
			MapData: []MapDataPoint{
				{
					Lat: 2,
					Lon: 2,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 4,
						},
					},
				},
				{
					Lat: 4,
					Lon: 4,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 5,
						},
					},
				},
			},
			Metadata: map[string]MetadataItem{
				"connections": {
					Title: "Connections",
					Value: 9,
				},
			},
		},
		{
			Value: "Leeds",
			MapData: []MapDataPoint{
				{
					Lat: 1,
					Lon: 1,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 8,
						},
					},
				},
			},
			Metadata: map[string]MetadataItem{
				"connections": {
					Title: "Connections",
					Value: 8,
				},
			},
		},
		{
			Value: "Bristol",
			MapData: []MapDataPoint{
				{
					Lat: 3,
					Lon: 3,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 6,
						},
					},
				},
			},
			Metadata: map[string]MetadataItem{
				"connections": {
					Title: "Connections",
					Value: 6,
				},
			},
		},
	}

	got := createMostActiveCitiesData(aggregatedLogs)

	if got[0].Value.(string) != want[0].Value.(string) ||
		got[1].Value.(string) != want[1].Value.(string) ||
		got[2].Value.(string) != want[2].Value.(string) {

		t.Error("Most active cities data was incorrect, got: ", got, ", want: ", want)
	}
}

func TestCreateMostActiveCountriesData(t *testing.T) {
	aggregatedLogs := AggregatedLogs{
		AggregatedLogEntry{Lat: 0, Lon: 0, Count: 1, Country: "Belgium"},
		AggregatedLogEntry{Lat: 1, Lon: 1, Count: 8, Country: "Mexico"},
		AggregatedLogEntry{Lat: 2, Lon: 2, Count: 4, Country: "United Kingdom"},
		AggregatedLogEntry{Lat: 3, Lon: 3, Count: 6, Country: "Italy"},
		AggregatedLogEntry{Lat: 4, Lon: 4, Count: 5, Country: "United Kingdom"},
	}

	want := []StatsDataPoint{
		{
			Value: "United Kingdom",
			MapData: []MapDataPoint{
				{
					Lat: 2,
					Lon: 2,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 4,
						},
					},
				},
				{
					Lat: 4,
					Lon: 4,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 5,
						},
					},
				},
			},
			Metadata: map[string]MetadataItem{
				"connections": {
					Title: "Connections",
					Value: 9,
				},
			},
		},
		{
			Value: "Mexico",
			MapData: []MapDataPoint{
				{
					Lat: 1,
					Lon: 1,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 8,
						},
					},
				},
			},
			Metadata: map[string]MetadataItem{
				"connections": {
					Title: "Connections",
					Value: 8,
				},
			},
		},
		{
			Value: "Italy",
			MapData: []MapDataPoint{
				{
					Lat: 3,
					Lon: 3,
					Metadata: map[string]MetadataItem{
						"connections": {
							Title: "Connections",
							Value: 6,
						},
					},
				},
			},
			Metadata: map[string]MetadataItem{
				"connections": {
					Title: "Connections",
					Value: 6,
				},
			},
		},
	}

	got := createMostActiveCountriesData(aggregatedLogs)

	if got[0].Value.(string) != want[0].Value.(string) ||
		got[1].Value.(string) != want[1].Value.(string) ||
		got[2].Value.(string) != want[2].Value.(string) {

		t.Error("Most active countries data was incorrect, got: ", got, ", want: ", want)
	}
}

func TestCreateMostIpAddressesData(t *testing.T) {
	aggregatedLogs := AggregatedLogs{
		AggregatedLogEntry{Lat: 0, Lon: 0, IpAddresses: []string{"0.0.0.0", "1.1.1.1"}},
		AggregatedLogEntry{Lat: 1, Lon: 1, IpAddresses: []string{"2.2.2.2"}},
		AggregatedLogEntry{Lat: 2, Lon: 2, IpAddresses: []string{"3.3.3.3", "4.4.4.4", "5.5.5.5"}},
	}

	want := []StatsDataPoint{
		{
			Value: 3,
			MapData: []MapDataPoint{
				{
					Lat: 2,
					Lon: 2,
				},
			},
			Metadata: map[string]MetadataItem{
				"ip_addresses": {
					Title: "IP addresses",
					Value: []string{"3.3.3.3", "4.4.4.4", "5.5.5.5"},
				},
			},
		},
		{
			Value: 2,
			MapData: []MapDataPoint{
				{
					Lat: 0,
					Lon: 0,
				},
			},
			Metadata: map[string]MetadataItem{
				"ip_addresses": {
					Title: "IP addresses",
					Value: []string{"0.0.0.0", "1.1.1.1"},
				},
			},
		},
		{
			Value: 1,
			MapData: []MapDataPoint{
				{
					Lat: 1,
					Lon: 1,
				},
			},
			Metadata: map[string]MetadataItem{
				"ip_addresses": {
					Title: "IP addresses",
					Value: []string{"2.2.2.2"},
				},
			},
		},
	}

	got := createMostIpAddressesData(aggregatedLogs)

	if got[0].Value.(int) != want[0].Value.(int) ||
		got[1].Value.(int) != want[1].Value.(int) ||
		got[2].Value.(int) != want[2].Value.(int) {

		t.Error("Most ip addresses data was incorrect, got: ", got, ", want: ", want)
	}
}

func TestCreateMostIngressPortsData(t *testing.T) {
	aggregatedLogs := AggregatedLogs{
		AggregatedLogEntry{Lat: 0, Lon: 0, IngressPorts: []string{"0", "1"}},
		AggregatedLogEntry{Lat: 1, Lon: 1, IngressPorts: []string{"2"}},
		AggregatedLogEntry{Lat: 2, Lon: 2, IngressPorts: []string{"3", "4", "5"}},
	}

	want := []StatsDataPoint{
		{
			Value: 3,
			MapData: []MapDataPoint{
				{
					Lat: 2,
					Lon: 2,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: []string{"3", "4", "5"},
				},
			},
		},
		{
			Value: 2,
			MapData: []MapDataPoint{
				{
					Lat: 0,
					Lon: 0,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: []string{"0", "1"},
				},
			},
		},
		{
			Value: 1,
			MapData: []MapDataPoint{
				{
					Lat: 1,
					Lon: 1,
				},
			},
			Metadata: map[string]MetadataItem{
				"ingress_ports": {
					Title: "Ingress ports",
					Value: []string{"2"},
				},
			},
		},
	}

	got := createMostIngressPortsData(aggregatedLogs)

	if got[0].Value.(int) != want[0].Value.(int) ||
		got[1].Value.(int) != want[1].Value.(int) ||
		got[2].Value.(int) != want[2].Value.(int) {

		t.Error("Most ingress ports data was incorrect, got: ", got, ", want: ", want)
	}
}
