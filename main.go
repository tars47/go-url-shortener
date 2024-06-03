package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/tars47/go-url-shortener/urlshort"
)

var yamlFileName string

func init() {
	flag.StringVar(&yamlFileName, "yaml", "url.yaml", "a yaml file that contains path and url that maps to")
}

func main() {

	flag.Parse()

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/dgo":  "https://go.dev/",
		"/lgo":  "https://go.dev/learn/",
		"/pkgo": "https://pkg.go.dev/",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback

	yamlHandler := urlshort.YAMLHandler(yamlFileName, mapHandler)

	fmt.Println("Starting the server on :4747")
	log.Fatal(http.ListenAndServe(":4747", yamlHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", http.NotFound)
	return mux
}
