package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

/*WebStream will get the full webstream for the supplied album*/
func (c *Client) WebStream(done <-chan struct{}, errs chan error) <-chan PhotoDownload {
	out := make(chan PhotoDownload)
	streamURL := fmt.Sprintf("%s/webstream", c.BaseURL)
	j, err := json.Marshal(map[string]interface{}{"streamCtag": nil})
	if err != nil {
		errs <- err
		close(out)
		return out
	}
	req, err := http.NewRequest("POST", streamURL, bytes.NewReader(j))
	if err != nil {
		errs <- err
		close(out)
		return out
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		errs <- err
		close(out)
		return out
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		close(out)
		return out
	}
	s := &Stream{}
	err = json.Unmarshal(b, s)
	if err != nil {
		errs <- fmt.Errorf("%s; Invalid album ID supplied?", err.Error())
		close(out)
		return out
	}
	go func() {
		for _, photo := range s.Photos {
			pd := PhotoDownload{
				PhotoGUID: photo.PhotoGUID,
			}
			var fileSize int
			for _, derivative := range photo.Derivatives {
				i, err := strconv.Atoi(derivative.FileSize)
				if err != nil {
					return
				}
				if i > fileSize {
					fileSize = i
					pd.Checksum = derivative.Checksum
				}
			}
			select {
			case out <- pd:
			case <-done:
				return
			}
		}
		close(out)
	}()
	return out

}

/*PhotoDownload is the struct we we will use for downloading*/
type PhotoDownload struct {
	PhotoGUID string
	Checksum  string
}

/*Stream contains the entire stream output*/
type Stream struct {
	StreamCTag    string  `json:"streamCtag"`
	ItemsReturned string  `json:"itemsReturned"`
	UserLastName  string  `json:"userLastName"`
	UserFirstName string  `json:"userFirstName"`
	StreamName    string  `json:"streamName"`
	Photos        []Photo `json:"photos"`
}

/*Photo describes a single photo*/
type Photo struct {
	BatchGUID            string                `json:"batchGuid"`
	Derivatives          map[string]Derivative `json:"derivatives"`
	ContributorLastName  string                `json:"contributorLastName"`
	BatchDateCreated     time.Time             `json:"batchDateCreated"`
	DateCreated          time.Time             `json:"dateCreated"`
	ContributorFirstName string                `json:"contributorFirstName"`
	PhotoGUID            string                `json:"photoGuid"`
	ContributorFullName  string                `json:"contributorFullName"`
	Width                string                `json:"width"`
	Caption              string                `json:"caption"`
	Height               string                `json:"height"`
}

/*Derivative is the struct we will use to determine which size of the picture to download*/
type Derivative struct {
	FileSize string `json:"fileSize"`
	Checksum string `json:"checksum"`
	Width    string `json:"width"`
	Height   string `json:"height"`
}
