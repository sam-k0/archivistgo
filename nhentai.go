package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type HentaiDict struct {
	ID      int    `json:"id"`
	MediaID string `json:"media_id"`
	Title   struct {
		English  string `json:"english"`
		Japanese string `json:"japanese"`
		Pretty   string `json:"pretty"`
	} `json:"title"`
	Images struct {
		Pages []struct {
			T   string `json:"t"`
			W   int    `json:"w"`
			H   int    `json:"h"`
			Url string
		} `json:"pages"`
		Cover struct {
			T   string `json:"t"`
			W   int    `json:"w"`
			H   int    `json:"h"`
			Url string
		} `json:"cover"`

		Thumbnail struct {
			T   string `json:"t"`
			W   int    `json:"w"`
			H   int    `json:"h"`
			Url string
		} `json:"thumbnail"`
	} `json:"images"`
	Scanlator  string `json:"scanlator"`
	UploadDate int    `json:"upload_date"`
	Tags       []struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		URL   string `json:"url"`
		Count int    `json:"count"`
	} `json:"tags"`
	NumPages     int `json:"num_pages"`
	NumFavorites int `json:"num_favorites"`
}

func ReadHentaiDownloadFile(fpath string) []string {
	// Open the file
	file, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Read the file
	links := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		links = append(links, scanner.Text())
	}
	return links
}

func getByID(id string) HentaiDict {
	// Perform a GET request
	resp, err := http.Get("http://nhentai.net/api/gallery/" + id)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// Serialize the response body into a HentaiDict struct
	var hentai HentaiDict
	err = json.Unmarshal(body, &hentai)
	if err != nil {
		panic(err)
	}

	// Fill the URL fields now that we have the MediaID
	for i := range hentai.Images.Pages {
		hentai.Images.Pages[i].Url = fmt.Sprintf("https://i.nhentai.net/galleries/%s/%d.jpg", hentai.MediaID, i+1)
	}
	// Cover and Thumbnail URLs
	hentai.Images.Cover.Url = fmt.Sprintf("https://t.nhentai.net/galleries/%s/cover.jpg", hentai.MediaID)
	hentai.Images.Thumbnail.Url = fmt.Sprintf("https://t.nhentai.net/galleries/%s/thumb.jpg", hentai.MediaID)
	return hentai
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
