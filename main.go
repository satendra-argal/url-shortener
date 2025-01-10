package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// UrlShortener defines the main structure for the URL shortener.
// It contains maps for URL mappings and a mutex for thread-safe operations.

type UrlShortener struct {
	urlMap      map[string]string
	reverseMap  map[string]string
	domainCount map[string]int
	mutex       sync.Mutex
}

// Initialize the Url Shortener instance

var shortener = UrlShortener{
	urlMap:      make(map[string]string),
	reverseMap:  make(map[string]string),
	domainCount: make(map[string]int),
}

// Initialize a local random number generator
var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateShortURL creates a random 6 character string for the short url

func generateShortURL() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[randGen.Intn(len(charset))]
	}
	return string(b)
}

// Constants for the base URL and the charset for short URL generation
const baseURL = "http://localhost:8080/"
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//DomainExtraction from the Given URL

func DomainExtract(url string) string {
	if strings.HasPrefix(url, "http://") {
		url = strings.TrimPrefix(url, "http://")

	} else if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
	}
	domain := strings.Split(url, "/")[0]
	return domain
}

// shortenURLHandler handles the POST request to shorten a URL
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var requestData struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	// Lock the mutex for thread-safe operations
	shortener.mutex.Lock()

	defer shortener.mutex.Unlock()

	if shortURL, exists := shortener.reverseMap[requestData.URL]; exists {
		response := map[string]string{"short_url": baseURL + shortURL}
		json.NewEncoder(w).Encode(response)
		return
	}
	shortURL := generateShortURL()
	for _, exists := shortener.urlMap[shortURL]; exists; {
		shortURL = generateShortURL()
	}

	shortener.urlMap[shortURL] = requestData.URL

	shortener.reverseMap[requestData.URL] = shortURL

	domain := DomainExtract(requestData.URL)
	shortener.domainCount[domain]++

	response := map[string]string{"short_url": baseURL + shortURL}
	json.NewEncoder(w).Encode(response)

}

// redirectHandler handles redirection requests for short URLs
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")

	// Lock the mutex for thread-safe operations
	shortener.mutex.Lock()
	defer shortener.mutex.Unlock()

	originalURL, exists := shortener.urlMap[shortURL]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

// topDomains handles requests to retrieve domain Count
func topDomains(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Lock the mutex for thread-safe operations
	shortener.mutex.Lock()
	defer shortener.mutex.Unlock()

	// Map to hold the count of each domain
	domainCounts := make(map[string]int)

	// Iterate through the reverse map and extract domain from each original URL
	for originalURL := range shortener.reverseMap {
		domain := DomainExtract(originalURL)
		domainCounts[domain]++
	}

	// Convert the domainCounts map to a slice of structs for sorting
	type domainCount struct {
		Domain string `json:"domain"`
		Count  int    `json:"count"`
	}

	var counts []domainCount
	for domain, count := range domainCounts {
		counts = append(counts, domainCount{Domain: domain, Count: count})
	}

	// Sort the counts in descending order by count
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	// Get the top 3 domains
	if len(counts) > 3 {
		counts = counts[:3]
	}

	// Encode the result as JSON and send the response
	json.NewEncoder(w).Encode(counts)

}

// main function up the HTTP server and routes
func main() {
	fmt.Println("Url Shortener started........ ")

	// Define routes for URL shortening, metrics, and redirection
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/top-domains", topDomains)

	log.Println("server started at : 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
