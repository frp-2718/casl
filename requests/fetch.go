// Requests package builds and handles HTTP requests.
package requests

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	multiwhere_url = "https://www.sudoc.fr/services/multiwhere/"
	iln2rcr_url    = "https://www.idref.fr/services/iln2rcr/"
	marcxml_url    = "https://www.sudoc.fr/"
)

const max_params = 20 // must be > 0

type Fetcher interface {
	FetchAll(ppns []string) [][]byte
	FetchRCR(ilns []string) []byte
	FetchMarc(ppn string) []byte
}

type HttpFetch struct{}

// FetchAll returns all xml data corresponding to each ppn in a map, or
// nil if unsuccessful.
// Note that the SUDOC API ignores unknown PPNs when requested with a
// muli-request.
func (f *HttpFetch) FetchAll(ppns []string) [][]byte {
	var xmlBatch [][]byte
	for _, url := range buildURLs(ppns) {
		xmlBatch = append(xmlBatch, fetch(url))
	}
	return xmlBatch
}

// FetchRCR returns a XML iln2rcr response from a list of ILNs.
func (f *HttpFetch) FetchRCR(ilns []string) []byte {
	return fetch(iln2rcr_url + strings.Join(ilns, ","))
}

func (f *HttpFetch) FetchMarc(ppn string) []byte {
	return fetch(marcxml_url + ppn + ".xml")
}

// Fetch returns the xml record corresponding to the given URL, or nil if
// unsucessful
func fetch(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		// if request time out, just ignore
		// TODO: delay and request again
		// TODO: handle other url.errors
		log.Println(err)
		return []byte{}
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("fetch: HTTP status code = %d", resp.StatusCode)
		return []byte{}
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return data
}

func buildURLs(params []string) []string {
	var urls []string
	for len(params) > max_params {
		newUrl := multiwhere_url + strings.Join(params[:max_params], ",")
		urls = append(urls, newUrl)
		params = params[max_params:]
	}
	urls = append(urls, multiwhere_url+strings.Join(params, ","))
	return urls
}
