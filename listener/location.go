package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// struct for ip location data
type IpLocationData struct {
	Location string
	Lat      float64
	Lon      float64
}

// retrieves ip location data from ip-api.com
func getIpLocationData(ipAddress string) IpLocationData {

	ipLocationData := IpLocationData{}

	response, err := http.Get("http://ip-api.com/json/" + ipAddress + "?fields=49369")
	if err != nil {
		log.Fatal("Error querying the ip-api: ", err.Error())
	}

	// retry if necessary
	if response.StatusCode == 429 {

		backoffStr := response.Header.Get("X-Ttl")
		backoff, err := strconv.ParseInt(backoffStr, 10, 64)
		if err != nil {
			log.Fatal("Error converting backoff header to int: ", err.Error())
		}

		time.Sleep(time.Duration(backoff+1) * time.Second)
		ipLocationData = getIpLocationData(ipAddress)

	} else {

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal("Error reading response body from ip-api: ", err.Error())
		}
		defer response.Body.Close()

		bodyParams := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyParams)
		if err != nil {
			log.Fatal("Error unmarshalling body params from ip-api: ", err.Error())
		}

		city := getOrEmptyString(bodyParams, "city")
		state := getOrEmptyString(bodyParams, "regionName")
		country := getOrEmptyString(bodyParams, "country")

		ipLocationData.Location = createLocationString(city, state, country)
		ipLocationData.Lat = getOrEmptyFloat(bodyParams, "lat")
		ipLocationData.Lon = getOrEmptyFloat(bodyParams, "lon")
	}

	return ipLocationData
}

// generates a location string in the format "City, State, Country"
func createLocationString(city, state, country string) string {
	locations := make([]string, 0)
	if len(city) > 0 {
		locations = append(locations, city)
	}
	if len(state) > 0 {
		locations = append(locations, state)
	}
	if len(country) > 0 {
		locations = append(locations, country)
	}
	return strings.Join(locations, ", ")
}

// gets string from map or returns empty string if it does not exist
func getOrEmptyString(mp map[string]interface{}, key string) string {
	if str, ok := mp[key]; ok {
		return str.(string)
	}
	return ""
}

// gets float from map or returns 0 if it does not exist
func getOrEmptyFloat(mp map[string]interface{}, key string) float64 {
	if flt, ok := mp[key]; ok {
		return flt.(float64)
	}
	return 0
}
