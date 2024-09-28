package main

import (
	"io"
	"net/http"
	"os"
)

// Download a file from a URL and save it to a path
// url: The URL of the file
// path: The path to save the file
func downloadFile(url string, path string) {
	// Perform a GET request
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// Create the file
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}
}
