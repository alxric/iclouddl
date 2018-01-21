package client

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

//Client is the base struct
type Client struct {
	ID         string
	BaseURL    string
	HTTPClient *http.Client
}

//New sets up a new client
func New(id string) (*Client, error) {
	c := &Client{
		ID:         id,
		BaseURL:    fmt.Sprintf("https://sharedstreams.icloud.com/%s/sharedstreams", id),
		HTTPClient: &http.Client{Timeout: time.Millisecond * 10000},
	}
	return c, nil
}

//Do will start the downloader
func (c *Client) Do(path string) {
	errs := make(chan error)
	done := make(chan struct{})
	defer close(done)
	handleErrs(errs, done)
	stream := c.WebStream(done, errs)

	var photos []<-chan string

	for i := 0; i < runtime.NumCPU()*10; i++ {
		photos = append(photos, c.Download(done, c.PhotoStream(done, stream, errs, path), path))
	}

	downloads := merge(done, photos...)
	for download := range downloads {
		if download != "" {
			log.Println("Downloaded", download)
		}
	}
}

func merge(done <-chan struct{}, cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	output := func(c <-chan string) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func handleErrs(errs chan error, done chan struct{}) {
	go func() {
		select {
		case err := <-errs:
			log.Fatal(err)
		}
	}()

}
