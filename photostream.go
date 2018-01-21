package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//PhotoStream will get more detailed info about a specific photo
func (c *Client) PhotoStream(done <-chan struct{}, in <-chan PhotoDownload, errs chan error, path string) <-chan Downloader {
	out := make(chan Downloader)
	go func() {
		defer close(out)
		for pd := range in {
			downloader := Downloader{}
			if !fileExists(path, pd.PhotoGUID) {
				streamURL := fmt.Sprintf("%s/webasseturls", c.BaseURL)
				j, err := json.Marshal(map[string][]string{
					"photoGuids": []string{pd.PhotoGUID},
				})
				if err != nil {
					errs <- err
					return
				}
				req, err := http.NewRequest("POST", streamURL, bytes.NewReader(j))
				if err != nil {
					errs <- err
					return
				}
				req.Header.Set("Content-Type", "application/json")
				resp, err := c.HTTPClient.Do(req)
				if err != nil {
					errs <- err
					return
				}
				defer resp.Body.Close()
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					errs <- err
					return
				}
				details := PhotoDetails{}
				err = json.Unmarshal(b, &details)
				if err != nil {
					errs <- err
					return
				}
				item := photoURL(details, pd)
				downloader.PhotoGUID = pd.PhotoGUID
				downloader.Item = item
			}
			select {
			case out <- downloader:
			case <-done:
				return
			}
		}
	}()
	return out
}

func fileExists(path string, guid string) bool {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), guid) {
			return true
		}
	}
	return false
}

func photoURL(details PhotoDetails, pd PhotoDownload) Item {
	for _, item := range details.Items {
		if strings.Contains(item.URLPath, pd.Checksum) {
			return item
		}
	}
	return Item{}
}

//Downloader is the struct that will be passed to our download function
type Downloader struct {
	PhotoGUID string
	Item      Item
}

//PhotoDetails contains photo details
type PhotoDetails struct {
	Locations struct {
		CvwsIcloudContentCom struct {
			Scheme string   `json:"scheme"`
			Hosts  []string `json:"hosts"`
		} `json:"cvws.icloud-content.com"`
	} `json:"locations"`
	Items map[string]Item `json:"items"`
}

//Item is URL details for a photo
type Item struct {
	URLExpiry   time.Time `json:"url_expiry"`
	URLLocation string    `json:"url_location"`
	URLPath     string    `json:"url_path"`
}
