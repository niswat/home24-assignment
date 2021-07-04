package webtool

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// This function counts the number of Internal, External and Broken Links. To assess whether a link is broken or not, we
// are following default redirection in go and marking it as broken if the final status code is not equal to `200`.
func IntExtBrokenLink(url string) (int, int, int, error) {

	resp, err := HttpResponse(url)
	if err != nil {
		return 0, 0, 0, err
	}
	// The client must close the response body when finished with it:
	defer resp.Body.Close()

	// Creating a Document Object to parse the html data
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
	if strings.Contains(match[3], "/") {
		match[3] = strings.Split(match[3], "/")[0]
	}
	str1 := fmt.Sprintf("http://www.%s", match[3])
	str2 := fmt.Sprintf("https://www.%s", match[3])
	str3 := fmt.Sprintf("http://%s", match[3])
	str4 := fmt.Sprintf("https://%s", match[3])

	// Regex to match internal Link as it can start from "/" ,  "#" , "" (nil) or "<somestring>" we also want to make sure
	// it doesn't matches any link starting with http as we are already handling this case with help of `str1..str4`.
	validInternalLink := regexp.MustCompile("^(/.*|#||([^h]|h[^t]|ht[^t]|htt[^p]).*)$")

	// Regex to Match for all external links.
	validExternalLink := regexp.MustCompile("^(http|https)://.*")

	// Traversing each Selection node having hyperlink tag `a`
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		fmt.Println("URL  : ", url)
		fmt.Println("HREF :", href)
		brokenUrl := ""

		// Verifying whether the value of `href` Attibute matches to the internal link which can be the regex compiled or
		// a link containing it's own url
		if validInternalLink.MatchString(href) || strings.Contains(href, str1) || strings.Contains(href, str2) ||
			strings.Contains(href, str3) || strings.Contains(href, str4) {
			internalLinksCount += 1
			if validInternalLink.MatchString(href) {
				brokenUrl = str4 + "/" + href
			} else {
				brokenUrl = href
			}
			fmt.Println("INTERNAL   : ", brokenUrl)
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
			fmt.Println("EXTERNAL   : ", href)
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

// This function gets the response object for the url to chceck whether it is broken or not.
func CheckBrokenLink(url string) (*http.Response, error) {
	brkResponse := &http.Response{}
	resp, err := HttpResponse(url)
	if err != nil {
		// handle the case where an internal/external url is pointing to localhost only and if so , return a custom status code
		if strings.Contains(err.Error(), "localhost") {
			brkResponse.StatusCode = 460
			return brkResponse, nil
			// if no response is returned in 10 sec , we assume it to be timed out.
		} else if strings.Contains(err.Error(), "timeout") {
			brkResponse.StatusCode = 408
			return brkResponse, nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	return resp, err
}
