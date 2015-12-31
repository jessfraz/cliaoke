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
	re = regexp.MustCompile(`<!-- Usage of azlyrics.com content by any third-party lyrics provider is prohibited by our licensing agreement. Sorry about that. -->((.|\n|\r)*?)</div>`)
)

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
		return "", fmt.Errorf("[%s] regex parsing failed for body:\n%s\n", query, body)
	}

	return html[0], nil
}
