package read

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

type HttpGetter struct {
	Reqs []string
}

// Get makes an http get request and throws an Error if the response
// statuscode is not 200
func (hg *HttpGetter) Get(url string) (result []byte, err error) {

	resp, err := myClient.Get(url)
	if err != nil {
		hg.Reqs = append(hg.Reqs, url+" err="+err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	status := resp.StatusCode

	if resp.StatusCode == 429 {
		sleep := 1
		if header := resp.Header.Get("X-Contentful-RateLimit-Reset"); header != "" {
			if conv, err := strconv.ParseInt(header, 10, 64); err != nil {
				sleep = int(conv)
			}
		}
		hg.Reqs = append(hg.Reqs, fmt.Sprintf("%s %d", url, status))
		time.Sleep(time.Duration(sleep) * time.Second)
		return hg.Get(url)
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Http request failed: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		hg.Reqs = append(hg.Reqs, fmt.Sprintf("%s %d err=%s", url, status, err.Error()))
		return nil, err
	}

	hg.Reqs = append(hg.Reqs, fmt.Sprintf("%s %d %d", url, status, len(body)))
	return body, err
}

func (hg *HttpGetter) Stats() []string {
	return hg.Reqs
}
