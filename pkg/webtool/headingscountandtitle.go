package webtool

import (
	"log"

	"github.com/PuerkitoBio/goquery"
)

// This function gets the title of a web page and counts the number of Headings for each level of headings `H1...H6` in a web page.
func HeadingsCountAndTitle(url string) (string, []int, error) {

	title := ""
	headingsCount := make([]int, 0)
	resp, err := HttpResponse(url)
	if err != nil {
		return title, headingsCount, err
	}
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
