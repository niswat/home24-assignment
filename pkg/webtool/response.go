package webtool

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

// This function creates a client request with a given url and returns a response object along with err if any.
func HttpResponse(url string) (*http.Response, error) {

	// Setting InsecureSkipVerify to `true` avoids checking for ssl certificates this is analogous to `--skip-ssl-verification` and
	// this should be used only for testing or in combination with VerifyConnection or VerifyPeerCertificate.
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true

	// Create a http client with connection timeout of 10 secs . Also setting  DisablekeepAlivewill to `true` will make sure
	// that the connection is used to the server for a single http request only.
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:      30,
		IdleConnTimeout:   30 * time.Second,
		DisableKeepAlives: true,
		TLSClientConfig:   tlsConfig,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// setting req.close to `true` prevents the re-use of TCP connections between requests to the same hosts.
	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
