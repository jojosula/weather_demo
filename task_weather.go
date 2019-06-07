package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	API_KEY     = "CWB-79BF19E7-ED3E-43DD-8745-1515EC55F73C"
	WEATHER_API = "https://opendata.cwb.gov.tw/api/v1/rest/datastore/F-C0032-001"

	TAIPEI     = "臺北市"
	NEW_TAIPEI = "新北市"
	TAOYUAN    = "桃園市"

	ONE_HOUR = 3600
)

func prepareRequest(location string) *http.Request {
	// new request
	req, err := http.NewRequest("GET", WEATHER_API, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	// add header
	req.Header.Add("Authorization", API_KEY)
	req.Header.Add("Accept", "application/json")
	// add query parameter
	query := req.URL.Query()
	query.Add("format", "JSON")
	query.Add("locationName", location)
	req.URL.RawQuery = query.Encode()
	return req
}

func makeRequest(location string) {
	client := &http.Client{}
	// prepare request
	req := prepareRequest(location)
	if req == nil {
		return
	}
	// do request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	// read body
	weatherData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	// unmarshal to struct
	var result WeatherAPIResponse
	err = json.Unmarshal(weatherData, &result)
	if err != nil {
		log.Fatal(err)
		return
	}
	// marshal to json string
	recordsStr, _ := json.Marshal(result.Records)
	data := GetWeatherData(location)
	if len(data.Data) == 0 {
		CreateWeatherData(location, string(recordsStr))
	} else {
		UpdateWeatherData(location, string(recordsStr))
	}
}

func TaskGetWeatherData() {
	for true {
		log.Printf("start get weather data task!")
		makeRequest(TAIPEI)
		makeRequest(NEW_TAIPEI)
		makeRequest(TAOYUAN)

		time.Sleep(ONE_HOUR * time.Second)
	}
}
