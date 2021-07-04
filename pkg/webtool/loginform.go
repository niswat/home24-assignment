package webtool

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

// This functions checks whether a given web page has a login form or not
func CheckLoginForm(url string) (string, error) {

	resp, err := HttpResponse(url)
	if err != nil {
		return "", err
	}
	// The client must close the response body when finished with it
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	loginFromCount := 0
	doc.Find("input").Each(func(index int, item *goquery.Selection) {
		// The idea is to look for value `password` for each attribute of kind `type`.
		// Also here we are capturing both login form (count --> 1) and signup form (count --> 2)
		typePassword, _ := item.Attr("type")
		if typePassword == "password" {
			loginFromCount += 1
		}
	})
	if loginFromCount != 0 {
		return "Yes", nil
	}
	return "No", nil
}
