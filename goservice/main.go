// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + "wiki.txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + "wiki.txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("./static/html/edit.html", "./static/html/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func helloHandler1(w http.ResponseWriter, r *http.Request) {
	log.Printf("/ => Hello Go World!")
	fmt.Fprintf(w, "/ => Hello Go World! v3")
}


func helloHandler2(w http.ResponseWriter, r *http.Request) {
	log.Printf("/hello => Hello Go World!")
	fmt.Fprintf(w, "/hello => Hello Go World! v3")
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
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	// Start listening
	fmt.Println("listening at localhost:8080")
	fmt.Println("Try http://localhost:8080/hello")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
