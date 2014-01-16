package main

import (
  "./crawler"
  "os"
  "strconv"
)

func main () {
  link := os.Args[1]
  wCount, _ := strconv.Atoi(os.Args[2])

  c := crawler.New(link,wCount)
  c.Run()

}
