//https://github.com/Jasonasante/Groupie-Tracker/blob/master/main.go

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

// collates the data taken from all API structs.
type Data struct {
	A Artist
	R Relation
	L Location
	D Date
}

// stores data from artist API struct.
type Artist struct {
	Id           uint     `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate uint     `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

// stores data from location API struct.
type Location struct {
	Locations []string `json:"locations"`
}

// stores data from date API struct.
type Date struct {
	Dates []string `json:"dates"`
}

// stores data from relation API struct.
type Relation struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

// the slices of structs are used to index the data of each artist from APIs.
// the map[string]json.RawMessage variables are used to unmarshal another layer
// when multiple layers are present.
var (
	artistInfo   []Artist
	locationMap  map[string]json.RawMessage
	locationInfo []Location
	datesMap     map[string]json.RawMessage
	datesInfo    []Date
	relationMap  map[string]json.RawMessage
	relationInfo []Relation
)

// handles error messages
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status) // sends an HTTP response with the status code, does not write
	if status == http.StatusNotFound {
		_, err := template.ParseFiles("error.html")
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		fmt.Println("HTTP status 404: Page Not Found")
		os.Exit(0)
	}
	if status == http.StatusInternalServerError {
		_, err := template.ParseFiles("error.html")
		if err != nil {
			fmt.Fprint(w, "HTTP status 500: Internal Server Error")
			return
		}
		fmt.Println("HTTP status 500: Internal Server Error")
		os.Exit(0)
	}
	if status == http.StatusBadRequest {
		_, err := template.ParseFiles("error.html")
		if err != nil {
			fmt.Fprint(w, "HTTP status 500: Internal Server Error")
			return
		}
		fmt.Println("HTTP status 400: Bad Request")
		os.Exit(0)
	}
}

// gets and stores data from Artist API
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

// gets and stores data from Location API
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

// gets and stores data from Dates API
func DatesData() []Date {
	var bytes []byte
	dates, err2 := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err2 != nil {
		log.Fatal()
	}
	datesData, err3 := io.ReadAll(dates.Body)
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

// gets and stores data from Relation API
func RelationData() []Relation {
	var bytes []byte
	relation, err2 := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err2 != nil {
		log.Fatal()
	}
	relationData, err3 := io.ReadAll(relation.Body)
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

// collates the data taken from all API slices into one data struct.
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

// home page handler
func homePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if r.URL.Path == "/artistInfo" {
			fmt.Println("Working 1")
			return
		}
		fmt.Println("Home page, error 1")
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	data := ArtistData()
	t, err := template.ParseFiles("index.html") // uses template.html to parse thru data
	if err != nil {
		fmt.Println("Home page, error 2")
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, data) // executes template.html
}

// handles the artist Page when artist image is clicked by receiving "ArtistName" value
func artistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artistInfo" { // checks if URL ends with 'artistInfo'
		fmt.Println("Error 1")
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	value := r.FormValue("ArtistName") // value variable stores the artist name as a form value
	if value == "" {                   // checks if value is empty
		fmt.Println("Error 2")
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	a := collectData()                // calls collectData, stores as a new variable
	var b Data                        // creates new variable named b
	for i, v := range collectData() { // ranges over collectData using i and v
		if value == v.A.Name { // checks if value is equal to v (in collectData), of the A field (Data struct), of Name (Artist struct)
			// aka it checks if value is equal to an artist in our database
			b = a[i] // assigns b variable to the collectData element at i
		}
	}
	t, err := template.ParseFiles("artistPage.html")
	if err != nil {
		fmt.Println("Error 3")
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, b) // executes template using data from b
}

// collection of webpage handlers
func main() {
	fmt.Println("Fetching server at port 8080...")
	http.HandleFunc("/", homePage)
	http.HandleFunc("/artistInfo", artistPage)
	http.ListenAndServe(":8080", nil)
}
