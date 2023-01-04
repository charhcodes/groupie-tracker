package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

const API string = "https://groupietrackers.herokuapp.com/api"

type Data struct {
	artists   ArtistData
	locations LocationData
	dates     DateData
	relations RelationData
}

type ArtistData struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	//Locations    string   `json:"locations"`
	//ConcertDates string `json:"concertDates"`
	//Relations string `json:"relations"`
}

type LocationData struct {
	Locations []string `json:"locations"`
}

type DateData struct {
	Dates []string `json:"dates"`
}

type RelationData struct {
	Relations []string `json:"relations"`
}

var (
	artistData   []ArtistData
	locationData []LocationData
	locationMap  map[string]json.RawMessage
	dateData     []DateData
	dateMap      map[string]json.RawMessage
	relationData []RelationData
	relationMap  map[string]json.RawMessage
)

// Handle HTTP error messages

// unmarshals artist data, returns ArtistData slice
func artistGet() []ArtistData {
	url := API + "/artists"

	response, err := http.Get(url)
	if err != nil {
		//fmt.Println("error fetching API")
		os.Exit(0)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//fmt.Println("error reading API")
		os.Exit(0)
	}
	err = json.Unmarshal(body, &artistData)
	if err != nil {
		//fmt.Println("error decoding API")
		os.Exit(0)
	}

	return artistData
}

// unmarshals location data, returns LocationData slice
func locationGet() []LocationData {
	url := API + "/locations"
	//make a slice of bytes
	var bytes []byte

	response, err := http.Get(url)
	if err != nil {
		//fmt.Println("error fetching API")
		os.Exit(0)
	}
	locationData1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//fmt.Println("error reading API")
		os.Exit(0)
	}
	err = json.Unmarshal(locationData1, &locationMap)
	if err != nil {
		//fmt.Println("error decoding API")
		os.Exit(0)
	}
	for _, l := range locationMap {
		for _, v := range l {
			bytes = append(bytes, v)
		}
	}
	err = json.Unmarshal(bytes, &locationData)
	if err != nil {
		//fmt.Println("error decoding API")
		os.Exit(0)
	}

	return locationData
}

// unmarshals date data, returns DateData slice
func dateGet() []DateData {
	url := API + "/dates"
	//make a slice of bytes
	var bytes []byte

	response, err := http.Get(url)
	if err != nil {
		//fmt.Println("error fetching API")
		os.Exit(0)
	}
	dateData1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//fmt.Println("error reading API")
		os.Exit(0)
	}
	err = json.Unmarshal(dateData1, &dateMap)
	if err != nil {
		//fmt.Println("error decoding API")
		os.Exit(0)
	}
	for _, l := range dateMap {
		for _, v := range l {
			bytes = append(bytes, v)
		}
	}
	err = json.Unmarshal(bytes, &dateData)
	if err != nil {
		//fmt.Println("error decoding API")
		os.Exit(0)
	}

	return dateData
}

// unmarshals relation data, returns RelationData slice
func relationGet() []RelationData {
	url := API + "/relations"
	//make a slice of bytes
	var bytes []byte

	response, err := http.Get(url)
	if err != nil {
		//fmt.Println("error fetching API")
		os.Exit(0)
	}
	relationData1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		//fmt.Println("error reading API")
		os.Exit(0)
	}
	err = json.Unmarshal(relationData1, &relationMap)
	if err != nil {
		//fmt.Println("error decoding API")
		os.Exit(0)
	}
	for _, l := range relationMap {
		for _, v := range l {
			bytes = append(bytes, v)
		}
	}
	err = json.Unmarshal(bytes, &relationData)
	if err != nil {
		//fmt.Println("error decoding API")
		os.Exit(0)
	}

	return relationData
}

func collateData() []Data {

}

// func main() {
// 	Unmarshal(artistLink)
// }
