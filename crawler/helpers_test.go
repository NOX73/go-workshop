package crawler

import (
  . "launchpad.net/gocheck"
  "bytes"
)

func ( s *SuiteT ) TestGetLinks (c *C) {

  host := "example.com"

  body := bytes.NewBufferString(`
    <html>
      <body>

        <a href="/link1"></a>
        <a href="/link2"></a>
        <a href="http://example.com/link1"></a>
        <a href="http://example.com/link1#anchor"></a>
        <a href="http://another.con/link1"></a>

      </body>
    </html>
  `)

  links := getLinksFromBody(body, host)

  c.Assert(links, HasLen, 4)

}

func ( s *SuiteT ) TestNormalizeLink (c *C) {
  var link string
  var ok bool

  link, ok = normilizeLink("/link1", "example.com")

  c.Assert(ok, Equals, true)
  c.Assert(link, Equals, "http://example.com/link1")

  link, ok = normilizeLink("/link1#anchor", "example.com")
  c.Assert(ok, Equals, true)
  c.Assert(link, Equals, "http://example.com/link1")

  link, ok = normilizeLink("http://another.com/link1#anchor", "example.com")
  c.Assert(ok, Equals, false)
}
