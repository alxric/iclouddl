package main

import (
	"flag"
	"fmt"
	"icloudphotos/client"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	albumPtr *string
	pathPtr  *string
)

func main() {
	c, err := client.New(*albumPtr)
	if err != nil {
		return
	}

	errs := make(chan error)
	done := make(chan struct{})
	defer close(done)
	handleErrs(errs, done)
	stream := c.WebStream(done, errs)

	var photos []<-chan string

	for i := 0; i < runtime.NumCPU()*10; i++ {
		photos = append(photos, c.Download(done, c.PhotoStream(done, stream, errs, *pathPtr), *pathPtr))
	}

	downloads := merge(done, photos...)
	for download := range downloads {
		if download != "" {
			fmt.Println("Downloaded", download)
		}
	}
}

func handleErrs(errs chan error, done chan struct{}) {
	go func() {
		select {
		case err := <-errs:
			log.Fatal(err)
		}
	}()

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

func init() {

	albumPtr = flag.String("a", "", "album id")
	pathPtr = flag.String("p", "", "path to write photos on disk")

	flag.Parse()
	if *albumPtr == "" || *pathPtr == "" {
		fmt.Println("Please supply the following arguments")
		flag.PrintDefaults()
		os.Exit(1)
	}
	*pathPtr = strings.TrimRight(*pathPtr, "/")
	_, err := os.Create(fmt.Sprintf("%s/%s", *pathPtr, time.Now().Format("2006-01-02 15:04:03")))

	if err != nil {
		log.Fatal("Unable to write to supplied path")
	}
	os.Remove(fmt.Sprintf("%s/%s", *pathPtr, time.Now().Format("2006-01-02 15:04:03")))
}
