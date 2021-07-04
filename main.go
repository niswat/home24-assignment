package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	webtool "github.com/niswat/home24-assignment/pkg/webtool"
)

type WebPage struct {
	Exists             bool
	Title              string
	HtmlVersion        string
	HeadingsCount      []int
	HasLoginForm       string
	InternalLinksCount int
	ExternalLinksCount int
	BrokenLinksCount   int
	MsgExists          bool
	Msg                string
}

// This function evaluates whether a given url is a valid url or not.
func Validate(url string) (string, bool) {
	match := true
	resp, err := webtool.HttpResponse(url)
	if err != nil || resp.StatusCode >= 400 {
		match = false
	}
	if !match {
		return "domain does not exist", match
	}
	defer resp.Body.Close()
	return "", match
}

// Calling ParseFiles function once at program initialization to parse all templates (if present) into a
// single *Template (template Caching)
var templates = template.Must(template.ParseFiles("index.html"))

// This functions helps in avoiding the usage of template.Parsefiles() in each handler if there are any
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {

	// Parse the html template and in case of an error, return an Internal Server Error (500)
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func scrapHandler(w http.ResponseWriter, req *http.Request) {

	// This is to handle the case if application is accessed at some other endpoint than `\`
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	if req.Method == "GET" {
		myhtmlcontent := WebPage{
			Exists:    false,
			MsgExists: false,
		}

		renderTemplate(w, "index", myhtmlcontent)
	} else if req.Method == "POST" {

		myurl := req.FormValue("url")
		msg, isThere := Validate(myurl)
		myhtmlcontent := WebPage{
			Exists:    false,
			MsgExists: true,
			Msg:       msg,
		}
		if !isThere {
			renderTemplate(w, "index", myhtmlcontent)
			return
		}

		internalLinksCount, externalLinksCount, brokenLinksCount, err := webtool.IntExtBrokenLink(myurl)
		if err != nil {
			fmt.Println("This is err")
			log.Fatal(err)
			return
		}

		title, headingsCount, err := webtool.HeadingsCountAndTitle(myurl)
		if err != nil {
			log.Fatal(err)
			return
		}

		version, _ := webtool.CheckHtmlVersion(myurl)
		if err != nil {
			log.Fatal(err)
			return
		}

		loginForm, err := webtool.CheckLoginForm(myurl)
		if err != nil {
			log.Fatal(err)
			return
		}
		exists := true

		myhtmlcontentf := WebPage{
			Exists:             exists,
			Title:              title,
			HtmlVersion:        version,
			HeadingsCount:      headingsCount,
			HasLoginForm:       loginForm,
			InternalLinksCount: internalLinksCount,
			ExternalLinksCount: externalLinksCount,
			BrokenLinksCount:   brokenLinksCount,
			MsgExists:          false,
		}
		renderTemplate(w, "index", myhtmlcontentf)
	}

}

func main() {
	http.HandleFunc("/", scrapHandler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
