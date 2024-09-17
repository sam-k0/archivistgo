package main

import (
	"os"
	"path/filepath"
	"strconv"

	gopdf "github.com/signintech/gopdf"
)

func addImage(pdf *gopdf.GoPdf, path string) {
	pdf.AddPage()
	pdf.Image(path, 0, 0, nil)
}

// Compile a directory of images into a PDF
// inpath: The directory of images
// outpath: The path to save the PDF
func CompilePdf(inpath string, outpath string) {
	println("Compiling PDF")
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	// Iterate through the directory
	files, err := os.ReadDir(inpath)
	if err != nil {
		panic(err)
	}
	// get number of files in directory
	amount := len(files)
	for i := 0; i < amount; i++ {
		filename := strconv.Itoa(i) + ".jpg"
		addImage(&pdf, filepath.Join(inpath, filename))
	}

	outpath = filepath.Join(outpath, "output.pdf")
	pdf.WritePdf(outpath)
	// delete the directory
	os.RemoveAll(inpath)
}
