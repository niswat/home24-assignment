package webtool

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func HeadingsCountAndTitle(url string) (string, []int, error) {

	title := ""
	headingsCount := make([]int, 0)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return title, headingsCount, err
	}
	// Making sure we close the writer after reading from it
	defer resp.Body.Close()

	// Creating a Document Object (DOM tree) to parse the html data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
		return title, headingsCount, err
	}

	title = doc.Find("TITLE").Text()
	headingsCount = append(headingsCount, doc.Find("H1").Length(), doc.Find("H2").Length(), doc.Find("H3").Length())
	headingsCount = append(headingsCount, doc.Find("H4").Length(), doc.Find("H5").Length(), doc.Find("H6").Length())

	return title, headingsCount, nil

}
