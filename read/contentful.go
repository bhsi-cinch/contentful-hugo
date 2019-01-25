package read

import (
	"fmt"
	"io"
)

const previewURL string = "https://preview.contentful.com"
const contentURL string = "https://cdn.contentful.com"
const contentURLTemplate = "/spaces/%s/environments/%s/content_types?access_token=%s&limit=200&locale=%s"
const itemsURLTemplate = "/spaces/%s/environments/%s/entries?access_token=%s&limit=200&locale=%s&skip=%d"

type Contentful struct {
	Getter     Getter
	ReadConfig ReadConfig
}

// Types will use Contentful's content_types endpoint to retrieve all content types from contentful
func (c *Contentful) Types() (rc io.ReadCloser, err error) {
	url := fmt.Sprintf(contentURLTemplate, c.ReadConfig.SpaceID, c.ReadConfig.Environment, c.ReadConfig.AccessToken, c.ReadConfig.Locale)
	return c.get(url)
}

// Items will use Contentful's entires endpoint to retrieve all 'items' from contetnful
func (c *Contentful) Items(skip int) (rc io.ReadCloser, err error) {
	url := fmt.Sprintf(itemsURLTemplate, c.ReadConfig.SpaceID, c.ReadConfig.Environment, c.ReadConfig.AccessToken, c.ReadConfig.Locale, skip)
	return c.get(url)
}

func (c *Contentful) get(endpoint string) (rc io.ReadCloser, err error) {
	urlBase := contentURL
	if c.ReadConfig.UsePreview {
		urlBase = previewURL
	}

	return c.Getter.Get(urlBase + endpoint)
}
