package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// struct for ip location data
type IpLocationData struct {
	City    string
	Country string
	Lat     float64
	Lon     float64
}

// retrieves ip location data from ip-api.com
func getIpLocationData(ipAddress string) IpLocationData {

	// return cached location data if possible
	if ipLocationData, ok := ipApiCache[ipAddress]; ok {
		return ipLocationData
	}

	response, err := http.Get("http://ip-api.com/json/" + ipAddress + "?fields=49361")
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
		return getIpLocationData(ipAddress)
	}

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

	ipLocationData := IpLocationData{
		City:    getOrEmptyString(bodyParams, "city"),
		Country: getOrEmptyString(bodyParams, "country"),
		Lat:     getOrEmptyFloat(bodyParams, "lat"),
		Lon:     getOrEmptyFloat(bodyParams, "lon"),
	}

	// write to cache
	ipApiCache[ipAddress] = ipLocationData

	return ipLocationData
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
