// Copyright 2020 Datawire

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func imageServer(w http.ResponseWriter, r *http.Request) {
	// Get the file requested by the registration number
	imageFile, err := os.Open(fmt.Sprintf("data/%v.jpg", path.Base(r.URL.Path)))

	// Get the Host header
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		defer imageFile.Close()

		// Write the image to the response.  All our images are jpeg.
		w.Header().Set("Content-Type", "image/jpeg")

		// Write out the image file to the request.
		io.Copy(w, imageFile)
	}
}

func main() {
	// For the demo, we can disable security checks.  Not normally recommended!
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Wire up the paths to their respective handlers
	http.HandleFunc("/", imageServer)

	// Start listening
	fmt.Println("listening at localhost:8080")
	fmt.Println("Try http://localhost:8080/N567M")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
