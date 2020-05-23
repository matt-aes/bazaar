// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

// An individual aircraft
type Aircraft struct {
	Registration string
	Model		string
	Price		int
}

// Result of an inventory search
type SearchResults struct {
	Title		string
	Items		[]Aircraft
}

// Specifications for a given aircraft model
type Specification struct {
	Model		string
	Type		string
	HP			int
	Seats		int
	Speed		int
	Range		int
	Load		int
}

// Page templates
var templates = template.Must(template.ParseFiles(
	"./static/html/edit.html", "./static/html/view.html",  "./static/html/results.html",
	"./static/html/home.html", "./static.html/detail.html"))

func helloHandler1(w http.ResponseWriter, r *http.Request) {
	log.Printf("/ => Hello Go World!")
	fmt.Fprintf(w, "/ => Hello Go World! v3")
}


func helloHandler2(w http.ResponseWriter, r *http.Request) {
	log.Printf("/hello => Hello Go World!")
	fmt.Fprintf(w, "/hello => Hello Go World! v3")
}

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

func getResultsPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("/fwd => about to forward to inventory service")
	res, err := http.Get("http://inventoryservice/all")
	log.Printf("     => returned from http.Get")

	data, _ := ioutil.ReadAll(res.Body)
	log.Printf("     => returned from ioutil.ReadAll()")
	log.Printf("data was %v", data)

	// We are returned a list of aircraft in JSON.  Create an array
	// and unmarshal the data.
	aircraft := make([]Aircraft, 0)
	json.Unmarshal(data, &aircraft)

	var results = &SearchResults{ Title: "Results", Items: aircraft }

	// Now, generate the page from the template.
	err = templates.ExecuteTemplate(w, "results.html", results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func forwardHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("/fwd => about to forward to python service")
	res, err := http.Get("http://pyservice/hello")

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("     => returned from http.Get")
	data, _ := ioutil.ReadAll(res.Body)
	log.Printf("     => returned from ioutil.ReadAll()")

	defer res.Body.Close()
	log.Printf("     => returned from res.Body.Close()")

	// Write to response
	fmt.Fprintf(w, "pyservice/hello returned %d", res.StatusCode)
	fmt.Fprintf(w, "\n  => %v", string(data))

	log.Printf("     => finished writing to ResponseWriter")
}


func main() {
	// For the demo, we can disable security checks.  Not normally recommended!
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Wire up the paths to their respective handlers
	http.HandleFunc("/", helloHandler1)
	http.HandleFunc("/fwd", forwardHandler)
	http.HandleFunc("/hello", helloHandler2)
	http.HandleFunc("/getjson", testJsonResponse)
	http.HandleFunc("/results", getResultsPage)

	// Start listening
	fmt.Println("listening at localhost:8080")
	fmt.Println("Try http://localhost:8080/hello")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
