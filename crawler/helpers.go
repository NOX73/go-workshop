package crawler

import (
  "github.com/opesun/goquery"
  "io"
  "net/url"
)

func getLinksFromBody(body io.Reader, host string) []string {
  tree, _ := goquery.Parse(body)

  links := tree.Find("a").Attrs("href")

  normals := make([]string, 0, len(links))

  for _, link := range links {
    link, ok := normilizeLink(link, host)
    if ok { normals = append(normals, link) }
  }

  return normals
}

func normilizeLink (link, host string) (string, bool) {
  url, err := url.Parse(link)
  if err != nil { return link, false }

  if url.Host == "" {
    url.Host = host
  } else if url.Host != host {
    return link, false
  }

  url.Scheme = "http"
  url.Fragment = ""

  return url.String(), true
}
