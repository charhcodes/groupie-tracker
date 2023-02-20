package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
	artistInfo   []Artist
	locationMap  map[string]json.RawMessage
	locationInfo []Location
	datesMap     map[string]json.RawMessage
	datesInfo    []Date
	relationMap  map[string]json.RawMessage
	relationInfo []Relation
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	t, err := template.ParseFiles("error.html")
	errorMssg := ""
	if status == http.StatusNotFound {
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
		errorMssg = "Error: HTTP status 404"
		fmt.Println(errorMssg)
		t.Execute(w, errorMssg)
	}
	if status == http.StatusInternalServerError {
		errorMssg = "Error: HTTP status 500"
		if err != nil {
			fmt.Fprint(w, errorMssg)
		}
		fmt.Println(errorMssg)
		t.Execute(w, errorMssg)
	}
	if status == http.StatusBadRequest {
		errorMssg = "Error: HTTP status 400"
		if err != nil {
			fmt.Fprint(w, errorMssg)
		}
		fmt.Println(errorMssg)
		t.Execute(w, errorMssg)
	}
}

func ArtistData() []Artist {
	artist, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal()
	}
	artistData, err := ioutil.ReadAll(artist.Body)
	if err != nil {
		log.Fatal()
	}
	json.Unmarshal(artistData, &artistInfo)
	return artistInfo
}

func LocationData() []Location {
	var bytes []byte
	location, err2 := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err2 != nil {
		log.Fatal()
	}
	locationData, err3 := ioutil.ReadAll(location.Body)
	if err3 != nil {
		log.Fatal()
	}
	err := json.Unmarshal(locationData, &locationMap)
	if err != nil {
		fmt.Println("error :", err)
	}
	for _, m := range locationMap {
		for _, v := range m {
			bytes = append(bytes, v)
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
	ArtistData()
	RelationData()
	LocationData()
	DatesData()

	dataData := make([]Data, len(artistInfo))
	for i := 0; i < len(artistInfo); i++ {
		dataData[i].A = artistInfo[i]
		dataData[i].R = relationInfo[i]
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
	t, err := template.ParseFiles("index.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/artistInfo" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}
	value := r.FormValue("ArtistName")
	if value == "" {
		errorHandler(w, r, http.StatusBadRequest)
		return
	}
	a := collectData()
	var b Data
	for i, ele := range collectData() {
		if value == ele.A.Name {
			b = a[i]
		}
	}
	t, err := template.ParseFiles("artistPage.html")
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}
	t.Execute(w, b)
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
