package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bhsi-cinch/contentful-hugo/extract"
	"github.com/bhsi-cinch/contentful-hugo/read"
	"github.com/bhsi-cinch/contentful-hugo/translate"
	"github.com/bhsi-cinch/contentful-hugo/write"
)

func main() {
	space := flag.String("space-id", os.Getenv("CONTENTFUL_API_SPACE"), "The contentful space id to export data from")
	env := flag.String("environment", os.Getenv("CONTENTFUL_API_ENV"), "The contentful environment to export data from")
	key := flag.String("api-key", os.Getenv("CONTENTFUL_API_KEY"), "The contentful delivery API access token")
	config := flag.String("config-file", "extract-config.toml", "Path to the TOML config file to load for export config")
	preview := flag.Bool("p", false, "Use contentful's preview API so that draft content is downloaded")
	flag.Parse()

	if *env == "" {
		*env = "master"
	}

	fmt.Printf("Begin contentful export from space '%s', environment '%s'\n", *space, *env)
	extractor := extract.Extractor{
		ReadConfig: read.ReadConfig{
			UsePreview:  *preview,
			SpaceID:     *space,
			Environment: *env,
			AccessToken: *key,
			Locale:      "en-US",
		},
		Getter: &read.HttpGetter{
			make([]string, 0),
		},
		TransConfig: translate.LoadConfig(*config),
		WStore:      write.FileStore{},
		RStore:      read.FileStore{},
	}
	s, err := extractor.ProcessAll()
	if err != nil {
		fmt.Println("Contentful export failed with the following error: " + err.Error())
		os.Exit(1)
	} else {

		fmt.Println("Contentful API Reqeusts:")
		for _, s := range s.HTTP {
			fmt.Println(s)
		}
		fmt.Println("Directory Level Files Written:")
		for _, s := range s.IndexFiles {
			fmt.Println(s)
		}
		fmt.Println("Entry Files Written:")
		for _, s := range s.Entries {
			fmt.Println(s)
		}
		fmt.Printf("Exported %d entries, and wrote %d index files for %d content types", len(s.Entries), len(s.IndexFiles), s.Types)
	}
}
