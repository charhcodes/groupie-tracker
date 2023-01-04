package main

// https://levelup.gitconnected.com/consuming-a-rest-api-using-golang-b323602ba9d8

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type artistData struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	Album        string              `json:"album"`
	Locations    string              `json:"locations`
	Concerts     map[string][]string `json:"concerts`
}

func main() {
	fmt.Println("Calling API...")
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://groupietrackers.herokuapp.com/api", nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject artistData
	json.Unmarshal(bodyBytes, &responseObject)

	fmt.Printf("%+v\n", responseObject)
}
