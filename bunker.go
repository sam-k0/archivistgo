package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/anaskhan96/soup"
)

type BunkrAlbum struct {
	AlbumName string
	AlbumURL  string
	ImageUrls []string
	VideoUrls []string
}

// Read the "bunkrlinks.txt" file and return the links
func ReadBunkrDownloadFile(fpath string) []string {
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

// Get all images and videos from a Bunkr album
// albumUrl: The URL of the album
func GetBunkrAlbum(albumUrl string) BunkrAlbum {
	// Make a GET request to the main page url
	resp, err := http.Get(albumUrl)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	currentAlbum := BunkrAlbum{
		AlbumURL: albumUrl,
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Get all links from the page
	medialinks := []string{}
	doc := soup.HTMLParse(string(body))
	links := doc.FindAll("a")
	for _, link := range links {
		// Save link if it contains a type of "https://bunkrrr.org/v" or "https://bunkrrr.org/i"
		if strings.Contains(link.Attrs()["href"], "https://bunkrrr.org/v") ||
			strings.Contains(link.Attrs()["href"], "https://bunkrrr.org/i") {
			medialinks = append(medialinks, link.Attrs()["href"])
		}
	}

	for _, link := range medialinks {
		println("Found link ", link)
	}

	// Every of the links link to a page with the download link
	for _, link := range medialinks {
		resp, err := http.Get(link)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		// Get all video html tags
		doc := soup.HTMLParse(string(body))
		videos := doc.FindAll("video")
		for _, video := range videos {
			source := video.Find("source")
			videourl := source.Attrs()["src"]
			videotype := source.Attrs()["type"]
			println("Adding video of type ", videotype, " with url ", videourl)
			currentAlbum.VideoUrls = append(currentAlbum.VideoUrls, videourl)
		}

		// Find all image html tags with link of "*.bunkr.ru"
		images := doc.FindAll("img")
		for _, image := range images {
			imageurl := image.Attrs()["src"]
			if strings.Contains(imageurl, ".bunkr.ru") {
				println("Adding image with url ", imageurl)
				currentAlbum.ImageUrls = append(currentAlbum.ImageUrls, imageurl)
			}
		}

	}
	return currentAlbum
}

// Download a file from a URL and save it to a path
// album: The album to download
// path: The parent directory to save the album
func DownloadBunkrAlbum(album BunkrAlbum, path string) {
	println("Downloading album ", album.AlbumURL)
	path = filepath.Join(path, strings.Split(album.AlbumURL, "/")[4])
	// Create the directory for the album
	os.MkdirAll(path, os.ModePerm)

	downloadpath := "" // Temporary variable to store the download path
	for i, image := range album.ImageUrls {
		downloadpath = filepath.Join(path, fmt.Sprintf("%d.jpg", i))
		println("Downloading image ", image)
		downloadFile(image, downloadpath)
	}
	for i, video := range album.VideoUrls {
		// Get the format of the video
		// The format is the last part of the URL
		// Example: https://bunkrrr.org/v/1234.mp4
		format := strings.Split(video, ".")
		downloadpath = filepath.Join(path, fmt.Sprintf("%d.%s", i, format[len(format)-1]))
		println("Downloading video ", video)
		downloadFile(video, downloadpath)
	}
}
