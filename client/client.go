package client

import (
	"fmt"
	"net/http"
	"time"
)

/*Client is the base struct*/
type Client struct {
	ID         string
	BaseURL    string
	HTTPClient *http.Client
}

/*New sets up a new client*/
func New(id string) (*Client, error) {
	c := &Client{
		ID:         id,
		BaseURL:    fmt.Sprintf("https://sharedstreams.icloud.com/%s/sharedstreams", id),
		HTTPClient: &http.Client{Timeout: time.Millisecond * 10000},
	}
	return c, nil
}
