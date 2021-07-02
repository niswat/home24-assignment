package main

import (
	"html/template"
	"log"
	"net/http"

	webtool "github.com/niswat/home24-assignment/pkg/webtool"
)

type WebPage struct {
	Title              string
	HtmlVersion        string
	HeadingsCount      []int
	HasLoginForm       string
	InternalLinksCount int
	ExternalLinksCount int
	BrokenLinksCount   int
}

// Calling ParseFiles function once at program initialization to parse all templates into a single *Template (template Caching)
var templates = template.Must(template.ParseFiles("index.html", "parse.html"))

func handler(w http.ResponseWriter, req *http.Request) {
	renderTemplate(w, "index", nil)

}

// This functions helps in avoiding the usage of template.Parsefiles() in each handler
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parse the html template and in case of an error, return an Internal Server Error (500)
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {

	// do whatever you need to do

	myurl := r.FormValue("url")
	// pass this url IntExtBrokenLink
	internalLinksCount, externalLinksCount, brokenLinksCount, err := webtool.IntExtBrokenLink(myurl)
	if err != nil {
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

	myhtmlcontent := WebPage{
		Title:              title,
		HtmlVersion:        version,
		HeadingsCount:      headingsCount,
		HasLoginForm:       loginForm,
		InternalLinksCount: internalLinksCount,
		ExternalLinksCount: externalLinksCount,
		BrokenLinksCount:   brokenLinksCount,
	}
	renderTemplate(w, "parse", myhtmlcontent)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/parse", postHandler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
