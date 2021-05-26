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

	podium.insert(AggregatedLogEntry{Count: 4})

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

	podium.insert(AggregatedLogEntry{Count: 3})

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

	podium.insert(AggregatedLogEntry{Count: 2})

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

	want := []map[string]interface{}{
		{
			"location": map[string]interface{}{
				"lat": 1,
				"lon": 1,
			},
			"data": map[string]interface{}{
				"connections":   4,
				"ingress_ports": []string{"1"},
				"ip_addresses":  []string{"1.1.1.1"},
			},
		},
		{
			"location": map[string]interface{}{
				"lat": 0,
				"lon": 0,
			},
			"data": map[string]interface{}{
				"connections":   3,
				"ingress_ports": []string{"0"},
				"ip_addresses":  []string{"0.0.0.0"},
			},
		},
		{
			"location": map[string]interface{}{
				"lat": 3,
				"lon": 3,
			},
			"data": map[string]interface{}{
				"connections":   2,
				"ingress_ports": []string{"3"},
				"ip_addresses":  []string{"3.3.3.3"},
			},
		},
	}

	got := createMostConnectionsData(aggregatedLogs)

	if got[0]["data"].(map[string]interface{})["connections"].(int) != 4 ||
		got[1]["data"].(map[string]interface{})["connections"].(int) != 3 ||
		got[2]["data"].(map[string]interface{})["connections"].(int) != 2 {

		t.Error("Most connections data was incorrect, got: ", got, ", want: ", want)
	}
}
