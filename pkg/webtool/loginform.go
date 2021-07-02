package webtool

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func CheckLoginForm(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	// Making sure we close the writer after reading from it
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
