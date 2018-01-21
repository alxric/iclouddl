package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	client "github.com/hummerpaskaa/iclouddl"
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

	c.Do(*pathPtr)
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
