package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var client *http.Client

type CatFact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

func getCatFact() {
	url := "https://catfact.ninja/fact"

	var catFact CatFact

	err := getJson(url, &catFact)
	if err != nil {
		fmt.Printf("error getting cat fact: %s\n", err.Error())
		return
	} else {
		fmt.Printf("A super interesting Cat Fact: %s\n", catFact.Fact)
	}
}

// for APIs with many sections and subsections
type RandomUser struct {
	Results []UserResult
}

type UserResult struct {
	Name    UserName
	Email   string
	Picture UserPicture
}

type UserName struct {
	Title string
	First string
	Last  string
}

type UserPicture struct {
	Large     string
	Medium    string
	Thumbnail string
}

func getRandomUser() {
	url := "https://randomuser.me/api"

	var user RandomUser

	err := getJson(url, &user)
	if err != nil {
		fmt.Printf("errpr getting json: %s\n", err.Error())
	} else {
		fmt.Printf("User: %s %s %s\nEmail: %s\nThumbnail: %s",
			user.Results[0].Name.Title,
			user.Results[0].Name.First,
			user.Results[0].Name.Last,
			user.Results[0].Email,
			user.Results[0].Picture.Thumbnail,
		)
	}
}

func getJson(url string, target interface{}) error { //allows us to pass through any interface
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func main() {
	client = &http.Client{Timeout: 10 * time.Second} //times out if no response for 10 seconds
	getCatFact()
	getRandomUser()
}
