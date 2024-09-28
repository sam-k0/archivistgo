package main

import (
	"os"
	"strconv"
)

func setup() {
	println("Setup...")
	// Check if directories exist
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", os.ModePerm)
	}
	if _, err := os.Stat("out"); os.IsNotExist(err) {
		os.Mkdir("out", os.ModePerm)
	}
	if _, err := os.Stat("bunkr"); os.IsNotExist(err) {
		os.Mkdir("bunkr", os.ModePerm)
	}
	print("Setup done.")
}

func main() {
	setup() // Create directories if not exist

	if len(os.Args) == 2 {
		println("Single argument mode: A single nhentai or bunkr album will be downloaded.")
		arg := os.Args[1]
		// decide if it is a nhentai or bunkr album, if its only numbers, its nhentai
		if _, err := strconv.Atoi(arg); err == nil {
			hentai := getByID(arg)
			path := DownloadHentai(hentai, "temp")
			CompilePdf(path, "out")
		} else {
			album := GetBunkrAlbum(arg)
			DownloadBunkrAlbum(album, "bunkr")
		}
		return
	}
	println("Batch mode: nhentai and bunkr albums will be downloaded from the links in the files.")
	// Download bunkr albums
	links := ReadBunkrDownloadFile("bunkrlinks.txt")
	for _, link := range links {
		album := GetBunkrAlbum(link)
		DownloadBunkrAlbum(album, "bunkr")
	}

	// Download nhentai albums
	links = ReadHentaiDownloadFile("hentailinks.txt")
	for _, link := range links {
		hentai := getByID(link)
		path := DownloadHentai(hentai, "temp")
		CompilePdf(path, "out")
	}
}
