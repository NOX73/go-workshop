package crawler

import (
  "net/url"
  "net/http"
  "log"
  "errors"
)

func New (link string, workersCount int) *crawler {
  url, err := url.Parse(link)
  if err != nil { panic(err) }

  host := url.Host

  return &crawler{host, workersCount, []string{link}, make(map[string]bool), make(chan []string), make(chan string)}
}

type crawler struct {
  Host string
  WorkersCount int
  Jobs []string
  visitedLinks map[string]bool
  In chan []string
  Out chan string
}

func (c *crawler) Run () {
  var link string
  link, _ = c.popJob()

  for i := c.WorkersCount; i > 0 ; i-- {
    c.runWorker(c.Out, c.In)
  }

  for {
    c.visitedLinks[link] = true

    select {
    case links := <- c.In:
      c.pushLinks(links)
    case c.Out <- link:
      link, _ = c.popJob()
    }

  }
}

func (c *crawler) pushLinks (links []string)  {
  for _, link := range links {
    if !c.isVisited(link) {
      c.pushJob(link)
    }
  }
}

func (c *crawler) processResp (resp *http.Response) []string {
  log.Println(resp.Status, resp.Request.URL.String())
  var links []string

  if resp.StatusCode == 200 {
    body := resp.Body
    links = getLinksFromBody(body, c.Host)
  }

  return links
}

func (c *crawler) isVisited (link string) bool {
  _, ok := c.visitedLinks[link]
  return ok
}

func (c *crawler) popJob () (string, error) {
  if len(c.Jobs) == 0 {return "", errors.New("") }

  var link string
  link, c.Jobs = c.Jobs[len(c.Jobs)-1], c.Jobs[:len(c.Jobs)-1]
  return link, nil
}

func (c *crawler) pushJob (link string) {
  c.Jobs = append(c.Jobs, link)
}

func (c *crawler) runWorker (in chan string, out chan []string) {
  go func(){
    for {
      link := <- in

      resp, err := http.Get(link)
      if err != nil {continue}

      defer resp.Body.Close()
      links := c.processResp(resp)
      out <- links

    }
  }()
}
