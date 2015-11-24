package kiteroot

import (
	"net/http"
	"strings"
	"sync"
	"testing"
)

func TestParse(t *testing.T) {
	html := `
<html>
<head>
<title>   KiteRoot </title>
<meta name="user" content="phynalle"/>
<meta name="profile" content="nothing"/>
</head>
<body>
omg
</body>
</html>
`

	r := strings.NewReader(html)
	_, err := Parse(r)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseWebPages(t *testing.T) {
	list := []string{
		"http://golang.org",
		"http://google.com",
		"http://facebook.com",
		"http://apple.com",
		"http://github.com",
	}

	var wg sync.WaitGroup
	wg.Add(len(list))
	for _, url := range list {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				t.Skipf("[%s] %s", url, err)
				wg.Done()
				return
			}
			defer resp.Body.Close()

			_, err = Parse(resp.Body)
			if err != nil {
				t.Errorf("[%s] %s", url, err)
			}
			wg.Done()
		}(url)
	}
	wg.Wait()
}
