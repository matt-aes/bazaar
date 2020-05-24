// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"sort"
)

// An individual aircraft, from inventoryservice/
type Aircraft struct {
	ImageURL     string `json: "imageUrl"`
	Registration string `json: "registration"`
	Model        string `json: "model"`
	Price        int    `json: "price"`
	DetailURL    string
	LocalPrice   string
}

// Specifications for a given aircraft model, from specsservice/
type Specification struct {
	Model string `json: "model"`
	Type  string `json: "type"`
	HP    int    `json: "hp"`
	Seats int    `json: "seats"`
	Speed int    `json: "speed"`
	Range int    `json: "range"`
	Load  int    `json: "load"`
}

// The home page parameters to the home.html template.
type HomePage struct {
	ResultsURL    string
	TitleImageURL string
}

// Result of an inventory search
type SearchResults struct {
	Title string
	Items []Aircraft
}

// Particular aircraft with its specifications.
type Detail struct {
	Title    string
	Aircraft Aircraft
	Specs    Specification
}

// Local currency: USD, EUR, NOK
var localCurrency = "USD"

// Page templates
var templates = template.Must(template.ParseFiles(
	"./static/html/home.html", "./static/html/results.html", "./static/html/detail.html"))

// Convert from US dollars to a local currency
func localizePrice(price int, currency string) string {
	// Default: price as given
	result := string(price)

	switch currency {
	case "USD":
		// Price is initially in US Dollars
		p := message.NewPrinter(language.English)
		result = p.Sprintf("$%d", price)
	case "EUR":
		// Convert from Dollars to Euros; use German formatting.
		exchangeRate :=  0.92 // USD to EUR
		p := message.NewPrinter(language.German)
		result = p.Sprintf("€%d", int(float64(price)/exchangeRate))
	case "NOK":
		// Convert from Dollars to Kroner; use Norwegian formatting.
		exchangeRate := 10.29 // USD to NOK
		p := message.NewPrinter(language.Norwegian)
		result = p.Sprintf("%d kr", int(float64(price)*exchangeRate))

	}

	return result
}

// ========== Test code ==========
type PersonalProfile struct {
	Name    string
	Hobbies []string
}

func testJsonResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("/testJsonHandler => a cool json structure")

	profile := PersonalProfile{"Bruce", []string{"flying", "telemark skiing", "travel", "running"}}

	// Return Json: marshal the struct.
	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Tell the client that the content type is json
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Home Page: show aircraft image and link to inventory
func getHomePage(w http.ResponseWriter, r *http.Request) {
	home := HomePage{
		ResultsURL:    filepath.Join(r.Host, "results"),
		TitleImageURL: filepath.Join(r.Host, "static", "images", "staggerwing.jpg")}

	err := templates.ExecuteTemplate(w, "home.html", home)

	if err != nil {
		log.Printf("template failed execution: %v", err)
	}
}

// Results Page: list entire inventory
func getResultsPage(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("http://inventoryservice/all")
	data, _ := ioutil.ReadAll(res.Body)

	// We are returned a list of aircraft in JSON.  Create an array
	// and unmarshal the data.
	aircraft := make([]Aircraft, 0)
	json.Unmarshal(data, &aircraft)

	// Sort the array of aircraft by price, low to high.
	sort.SliceStable(aircraft[:], func(i, j int) bool {
		return aircraft[i].Price < aircraft[j].Price
	})

	// Each aircraft has an ImageURL, but since we called the inventory service from inside
	// the cluster we have to rewrite it based on the host in our request.  Also, convert the
	// aircraft price to the local currency and format correctly.
	for i, ac := range aircraft {
		aircraft[i].ImageURL = filepath.Join(r.Host, "image", ac.Registration)
		aircraft[i].DetailURL = filepath.Join(r.Host, "detail", ac.Registration)
		aircraft[i].LocalPrice = localizePrice(ac.Price, localCurrency)
	}

	var results = &SearchResults{Title: "Aircraft Bazaar: Our Inventory", Items: aircraft}

	// Now, generate the page from the template.
	err = templates.ExecuteTemplate(w, "results.html", results)
	if err != nil {
		log.Printf("template failed execution: %v", err)
	}
}

// Detail Page: show the details of an individual aircraft
func getDetailPage(w http.ResponseWriter, r *http.Request) {
	// Get the registration number from the URL
	registration := path.Base(r.URL.Path)

	// Look up the individual aircraft from the inventory service
	res, err := http.Get("http://inventoryservice/one/" + registration)
	data, _ := ioutil.ReadAll(res.Body)

	aircraft := Aircraft{}
	json.Unmarshal(data, &aircraft)

	// Look up the specifications from the specsservice
	res, err = http.Get("http://specsservice/" + aircraft.Model)
	data, _ = ioutil.ReadAll(res.Body)

	specs := Specification{}
	json.Unmarshal(data, &specs)

	// Fix the imageURL so it will use our host path.
	aircraft.ImageURL = filepath.Join(r.Host, "image", aircraft.Registration)

	// Convert price to local currency.
	aircraft.LocalPrice = localizePrice(aircraft.Price, localCurrency)

	// Create the detail object to pass to the template.
	title := fmt.Sprintf("Aircraft Bazaar: %s %s", aircraft.Model, aircraft.Registration)
	detail := Detail{Title: title, Aircraft: aircraft, Specs: specs}

	// Now, generate the page from the template.
	err = templates.ExecuteTemplate(w, "detail.html", detail)

	if err != nil {
		log.Printf("template failed execution: %v", err)
	}
}


// Main: Set up transport, http handlers, static file serving, anmd start the listener.
func main() {
	// For the demo, we can disable security checks.  Not normally recommended!
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Wire up the paths to their respective handlers
	http.HandleFunc("/testjson", testJsonResponse)

	http.HandleFunc("/", getHomePage)
	http.HandleFunc("/results", getResultsPage)
	http.HandleFunc("/detail/", getDetailPage)

	// Handle static assets (images primarily)
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start listening
	fmt.Println("listening at localhost:8080")
	fmt.Println("Try http://localhost:8080/results")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
