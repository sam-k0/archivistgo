package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		println("Usage: go run main.go <id>")
		os.Exit(1)
	}

	id := os.Args[1]
	hentai := getByID(id)
	path := DownloadHentai(hentai, "temp")
	CompilePdf(path, "out")
}
