package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	A Artist
	R Relation
	L Location
	D Date
}

type Artist struct {
	Id           uint     `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate uint     `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	Locations []string `json:"locations"`
}

type Date struct {
	Dates []string `json:"dates"`
}

type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

var (
	artistInfo  []Artist                   // slice of artist structs
	locationMap map[string]json.RawMessage // maps a string key to a json.RawMessage value
	//RawMessage = byte slice that represents a JSON value, doesn't need to be parsed
	locationInfo []Location
	datesMap     map[string]json.RawMessage
	datesInfo    []Date
	relationMap  map[string]json.RawMessage
	relationInfo []Relation
)

// handles 404, 500, 400 errors
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	t, err := template.ParseFiles("error.html")
	errorMssg := ""
	if status == http.StatusNotFound { // if status = 404
		if err != nil {
			// if it is not giving an error (when there is one),
			// return a 500 error instead
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		errorMssg = "Error: HTTP status 404"
		fmt.Println(errorMssg)
		t.Execute(w, errorMssg)
	}
	if status == http.StatusInternalServerError { // if status = 500
		errorMssg = "Error: HTTP status 500"
		if err != nil {
			fmt.Fprint(w, errorMssg)
		}
		fmt.Println(errorMssg)
		t.Execute(w, errorMssg)
	}
	if status == http.StatusBadRequest { // if status = 400
		errorMssg = "Error: HTTP status 400"
		if err != nil {
			errorMssg = "Error: HTTP status 500"
			fmt.Fprint(w, errorMssg)
		}
		fmt.Println(errorMssg)
		t.Execute(w, errorMssg)
	}
}

func ArtistData() []Artist {
	// The code will read the data from a JSON response from GroupieTracker's API

	artist, err := http.Get("https://groupietrackers.herokuapp.com/api/artists") //grabs list of artists from link
	if err != nil {
		log.Fatal()
	}
	artistData, err := io.ReadAll(artist.Body) //reads data using ReadAll, stores in artistData
	if err != nil {
		log.Fatal()
	}
	json.Unmarshal(artistData, &artistInfo) //unmarshalls the data from artistData into the artistinfo struct
	return artistInfo
}

func LocationData() []Location {
	//  The code will take the JSON response from GroupieTracker and parse it into a map of Location data.

	var bytes []byte                                                                  // empty array of bytes
	location, err2 := http.Get("https://groupietrackers.herokuapp.com/api/locations") // gets locations from link, stores in location
	if err2 != nil {
		log.Fatal()
	}
	locationData, err3 := io.ReadAll(location.Body) // reads location data from JSON file, stores in locationData
	if err3 != nil {
		log.Fatal()
	}
	err := json.Unmarshal(locationData, &locationMap) // unmarshalls locationData, stores in locationMap struct
	if err != nil {
		fmt.Println("error :", err)
	}
	for _, m := range locationMap { // for every value in locationMap, m is created
		for _, v := range m { // for every value in m, v is created
			bytes = append(bytes, v) // each value is appended into the array of bytes from before
		}
	}
	err = json.Unmarshal(bytes, &locationInfo)
	if err != nil {
		fmt.Println("error :", err)
	}
	return locationInfo
}

func DatesData() []Date {
	var bytes []byte
	dates, err2 := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err2 != nil {
		log.Fatal()
	}
	datesData, err3 := ioutil.ReadAll(dates.Body)
	if err3 != nil {
		log.Fatal()
	}
	err := json.Unmarshal(datesData, &datesMap)
	if err != nil {
		fmt.Println("error :", err)
	}
	for _, m := range datesMap {
		for _, v := range m {
			bytes = append(bytes, v)
		}
	}
	err = json.Unmarshal(bytes, &datesInfo)
	if err != nil {
		fmt.Println("error :", err)
	}
	return datesInfo
}

func RelationData() []Relation {
	var bytes []byte
	relation, err2 := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err2 != nil {
		log.Fatal()
	}
	relationData, err3 := ioutil.ReadAll(relation.Body)
	if err3 != nil {
		log.Fatal()
	}
	err := json.Unmarshal(relationData, &relationMap)
	if err != nil {
		fmt.Println("error :", err)
	}

	for _, m := range relationMap {
		for _, v := range m {
			bytes = append(bytes, v)
		}
	}

	err = json.Unmarshal(bytes, &relationInfo)
	if err != nil {
		fmt.Println("error :", err)
	}
	return relationInfo
}

func collectData() []Data {
	// The code is used to collect data about the artist, relation, location and date

	// calls functions from before
	ArtistData()
	RelationData()
	LocationData()
	DatesData()

	dataData := make([]Data, len(artistInfo)) // an empty array of Data objects that will be used to temporarily store names, locations etc.
	for i := 0; i < len(artistInfo); i++ {    // iterates through artistInfo values
		dataData[i].A = artistInfo[i]   // uses i to assign values from artistInfo to the A field in dataData
		dataData[i].R = relationInfo[i] // same is done for R and relationInfo, L for locationInfo etc.
		dataData[i].L = locationInfo[i]
		dataData[i].D = datesInfo[i]
	}
	return dataData
}

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	data := ArtistData()
	t, err := template.ParseFiles("index.html") // parse thru data
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, data) // executes template
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artistInfo" { // checks if URL ends with 'artistInfo'
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	value := r.FormValue("ArtistName") // value variable stores the artist name as a form value
	if value == "" {                   // checks if value is empty
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	a := collectData()                  // calls collectData, stores as a new variable
	var b Data                          // creates new variable named b
	for i, ele := range collectData() { // ranges over collectData using i and v
		if value == ele.A.Name { // checks if value is equal to v (in collectData)
			// of the A field (Data struct), of Name (Artist struct)
			b = a[i] // assigns b variable to the collectData element at i
		}
	}
	t, err := template.ParseFiles("artistPage.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, b) // executes template using data from b
}

// collection of webpage handlers
func HandleRequests() {
	fmt.Println("Fetching server at port 8080...")
	http.HandleFunc("/", homePage)
	http.HandleFunc("/artistInfo", artistPage)
	http.ListenAndServe(":8080", nil)
}

func main() {
	HandleRequests()
}
