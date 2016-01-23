package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type Result struct {
	Title string `json:"Title"`
	URL   string `json:"Url"`
}

type BingResponse struct {
	D struct {
		Next    string   `json:"__next"`
		Results []Result `json:"results"`
	} `json:"d"`
}

var key = flag.String("key", "", "Key you get at azure market place.")

const (
	numItems = 10
)

func main() {
	flag.Parse()
	http.HandleFunc("/search", handleSearch)
	fmt.Println("Servindo em http://localhost:8080/search")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	// Query validation.
	query := r.FormValue("q")
	if query == "" {
		http.Error(w, `missing "q" URL parameter`, http.StatusBadRequest)
		return
	}
	// Invoking bing
	bingURL := fmt.Sprintf("https://api.datamarket.azure.com/Bing/Search/Web?$format=json&Adult=%%27Off%%27&$top=%d&Query=%%27%s%%27", query, numItems)
	req, _ := http.NewRequest("GET", bingURL, nil)
	req.SetBasicAuth(*key, *key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error calling bing: %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Parsing response.
	var bingResponse BingResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&bingResponse); err != nil {
		log.Printf("Error parsing response: %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Printing and returning results.
	for i, r := range bingResponse.D.Results {
		log.Printf("%d  -- %+v\n", i, r)
	}
	results, err := json.Marshal(bingResponse.D.Results)
	if err != nil {
		log.Printf("Error parsing response: %q", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(results))
	w.Header().Set("Context", "application/json")
}
