package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
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

func main() {
	fmt.Println("Url Shortener started........ ")

	log.Println("server started at : 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
