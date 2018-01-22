# iclouddl
Small tool to download photos from an Icloud shared stream. You will need your album ID, if the URL to your shared album is https://www.icloud.com/sharedalbum/#a251EQWOqz9e2jqx, then a251EQWOqz9e2jqx is your album ID.


Installation
---
go get github.com/hummerpaskaa/iclouddl

Example usage
---
  package main

  import (
      client "github.com/hummerpaskaa/iclouddl"
  )

  func main() 
      c, err := client.New("a251EQWOqz9e2jqx")
      if err != nil {
          return
      }

      c.Do("/tmp/pictures")
  }
