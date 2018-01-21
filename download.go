package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

//Download will download photos to disk
func (c *Client) Download(done <-chan struct{}, in <-chan Downloader, path string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for p := range in {
			var output string
			if p.PhotoGUID != "" && p.Item.URLPath != "" {
				dlURL := fmt.Sprintf("https://%s%s", p.Item.URLLocation, p.Item.URLPath)
				path := fmt.Sprintf("%s/%s.%s", path, p.PhotoGUID, fileName(dlURL))
				dl(dlURL, path)
				output = path
			}
			select {
			case out <- output:
			case <-done:
				return
			}
		}
	}()
	return out
}

func dl(url string, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func fileName(url string) string {
	u := strings.Split(url, "/")
	if len(u) < 6 {
		return ""
	}
	u = strings.Split(u[5], "?")
	u = strings.Split(u[0], ".")
	return u[1]
}
