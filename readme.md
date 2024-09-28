# archivist go
[![Go](https://github.com/sam-k0/archivistgo/actions/workflows/go.yml/badge.svg)](https://github.com/sam-k0/archivistgo/actions/workflows/go.yml)

nhentai.net and bunkr.ru downloader and pdf compiler.

## Usage

For the single album mode, pass an nhentaiID or bunkr-album link as an argument.
```sh
./archivistgo <id/bunkrlink>
```

In case of an nhentai ID, the program will download the album and compile it into a pdf.
In case of a bunkr link, the program will download the album's images and videos.

For the batch mode, populate the files `hentainumbers.txt`and `bunkrlinks.txt` with the respective IDs and links
and then run the progam without any arguments.
* For bunkr links, please add the album link, and not the direct link to the image or videos.
```sh
./archivistgo
```

## Build

```sh
go build
```

## Dependencies

- github.com/signintech/gopdf
- github.com/anaskhan96/soup