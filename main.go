package main

import (
	"flag"
	"os"

	"github.com/bhsi-cinch/contentful-hugo/extract"
	"github.com/bhsi-cinch/contentful-hugo/read"
	"github.com/bhsi-cinch/contentful-hugo/translate"
	"github.com/bhsi-cinch/contentful-hugo/write"
)

func main() {
	extractor := extract.Extractor{
		read.ReadConfig{
			"https://cdn.contentful.com",
			*flag.String("space-id", os.Getenv("CONTENTFUL_API_SPACE"), "The contentful space id to export data from"),
			*flag.String("api-key", os.Getenv("CONTENTFUL_API_KEY"), "The contentful delivery API access token"),
			"en-US",
		},
		read.HttpGetter{},
		translate.LoadConfig(*flag.String("config-file", "extract-config.toml", "Path to the TOML config file to load for export config")),
		write.FileStore{},
	}

	extractor.ProcessAll()
}
