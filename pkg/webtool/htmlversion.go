package webtool

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func CheckHtmlVersion(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	// Making sure we close the writer after reading from it
	defer resp.Body.Close()

	// Reading html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Creating a regex match for Doctype and making it case insensitive
	htmlVersionObj, err := regexp.Compile(`(?i)<!DOCTYPE .*>`)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	htmlVersion := htmlVersionObj.FindStringSubmatch(string(html))

	// Creaitng a regex to match for older version say 4.01
	finalCheck, err := regexp.Compile(`\d\.\d*`)
	version := finalCheck.FindStringSubmatch(htmlVersion[0])
	if len(version) == 0 {
		return "5.0", nil
	}
	return version[0], err
}
