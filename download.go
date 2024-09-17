package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

// Download a whole hentai
// hentai: The hentai to download
// path: The directory to save the hentai, will be populated with the cover and pages
func DownloadHentai(hentai HentaiDict, path string) string {
	path = filepath.Join(path, strconv.Itoa(hentai.ID))
	// Check if the directory exists
	if _, err := os.Stat(path); err == nil {
		fmt.Println("Hentai already downloaded")
		return path
	}

	println("Downloading hentai")
	os.MkdirAll(path, os.ModePerm)
	// Download the pages
	for i, page := range hentai.Images.Pages {
		downloadFile(page.Url, fmt.Sprintf("%s/%d.jpg", path, i))
		// Sleep for a bit to avoid rate limiting
		time.Sleep(time.Millisecond * 100)

		println("Downloaded", i, "of", len(hentai.Images.Pages))
	}
	return path
}
