package webtool

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// This fuction counts the number of Internal , External and broken Links. In case of broken links ,we are not
// handling the case of url redirection as the redirected code (say 301) can be different from the final response code
// (say 403/404) causing ambiguity in terms of which one to choose as a source of truth.
// Moreover we consider a link to be `not broken` if the final status code returned is 200

func IntExtBrokenLink(url string) (int, int, int, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return 0, 0, 0, err
	}

	// Making sure we close the writer after reading from it
	defer resp.Body.Close()

	// Creeating a Document Object to parse the html data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0, 0, 0, err
	}
	internalLinksCount := 0
	externalLinksCount := 0
	brokenLinksCount := 0

	// Creating url links `str1..str4` to capture the edge cases where internal links can start from http or https
	// and at the same time can have "www" or not.
	regObj, _ := regexp.Compile("^(http|https)://(www.|)(.*)")
	match := regObj.FindStringSubmatch(url)
	str1 := fmt.Sprintf("http://www.%s", match[3])
	str2 := fmt.Sprintf("https://www.%s", match[3])
	str3 := fmt.Sprintf("http://%s", match[3])
	str4 := fmt.Sprintf("https://%s", match[3])

	// Regext to match internal Link as it can start from "/" ,  "#" , "" (nil) or "<somestring>"
	// we also want to make sure it doesn't matches any link starting with http
	// as we are already handling this case with help of `str1..str4`.
	validInternalLink := regexp.MustCompile("^(/.*|#||([^h]|h[^t]|ht[^t]|htt[^p]).*)$")

	// Regex to Match for all external links.
	validExternalLink := regexp.MustCompile("^(http|https)://.*")

	// Traversing each Selecltion node with hyperlink tag `a`
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		//fmt.Println(href)
		brokenUrl := ""

		// Verifying whether the value of `href` Attibute matches to the internal link which can be the regex compiled or
		// a link containing it's own url
		if validInternalLink.MatchString(href) || strings.Contains(href, str1) || strings.Contains(href, str2) ||
			strings.Contains(href, str3) || strings.Contains(href, str4) {
			internalLinksCount += 1
			if validInternalLink.MatchString(href) {
				brokenUrl = url + "/" + href
			} else {
				brokenUrl = href
			}
			fmt.Println("INTERNAL  : ", brokenUrl)
			resp, err := CheckBrokenLink(brokenUrl)
			if err != nil {
				panic(err)
			}
			fmt.Println("STATUS CODE: ", resp.StatusCode)
			if resp.StatusCode != 200 {
				brokenLinksCount += 1
			}
			// Make sure the external link does not matches to any internal url
		} else if validExternalLink.MatchString(href) && href != str1+"/" && href != str2+"/" &&
			href != str3+"/" || href != str4+"/" {
			externalLinksCount += 1
			fmt.Println("EXTERNAL  : ", href)
			resp, err := CheckBrokenLink(href)
			if err != nil {
				panic(err)
			}
			fmt.Println("STATUS CODE: ", resp.StatusCode)
			if resp.StatusCode != 200 {
				brokenLinksCount += 1
			}
		}
		fmt.Println("------------------------------")
	})
	return internalLinksCount, externalLinksCount, brokenLinksCount, nil
}

// This function gets the response object for the url to chceck whether
func CheckBrokenLink(url string) (*http.Response, error) {

	// Create a http client with sensible timeout
	// Setting DisablekeepAlive will only use the connection to the server for a single
	// HTTP request and will avoid EOF error which
	tr := &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// setting close to `true` to prevent the re-use of TCP connections between requests to the same hosts
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// The client must close the response body when finished with it:
	defer resp.Body.Close()
	return resp, err
}
