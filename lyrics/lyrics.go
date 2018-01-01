package lyrics

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

const (
	queryURI = "http://search.azlyrics.com/search.php"
)

var (
	re     = regexp.MustCompile(`<!-- Usage of azlyrics.com content by any third-party lyrics provider is prohibited by our licensing agreement. Sorry about that. -->((.|\n|\r)*?)</div>`)
	reHTML = regexp.MustCompile(`<.+>`)
)

// Search scrapes azlyrics.com for song lyrics and does regex magic to clean them up.
// Beware, your IP can AND will get blocked while running this, but it is only
// called in `go generate` (see midi/generate.go) so a normal user will never run this.
func Search(query string) (string, error) {
	v := url.Values{
		"q": []string{query},
	}
	uri := fmt.Sprintf("%s?%s", queryURI, v.Encode())

	// start the scrape
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		return "", fmt.Errorf("scraping new document at %s failed: %v", uri, err)
	}

	link, ok := doc.Find("td").First().Find("a").Attr("href")
	if !ok {
		return "", fmt.Errorf("could not find top link at %s", uri)
	}

	// get the lyrics link
	resp, err := http.Get(link)
	if err != nil {
		return "", fmt.Errorf("request to %s failed: %v", link, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading body from %s failed: %v", link, err)
	}

	// get the lyrics from the HTML
	html := re.FindStringSubmatch(string(body))

	if len(html) <= 0 {
		return "", fmt.Errorf("[%s] regex parsing failed for body: %s", query, body)
	}

	// strip html tags from decoded lyrics
	lyrics := reHTML.ReplaceAllString(html[0], "")

	return lyrics, nil
}
