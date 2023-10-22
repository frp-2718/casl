// Requests package builds and handles HTTP requests.
package requests

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	multiwhere_url = "https://www.sudoc.fr/services/multiwhere/"
	iln2rcr_url    = "https://www.idref.fr/services/iln2rcr/"
	marcxml_url    = "https://www.sudoc.fr/"
)

const MAX_CONCURRENT_REQUESTS = 50

type Fetcher interface {
	FetchAll(ppns []string) [][]byte
	FetchRCR(ilns []string) []byte
	FetchMarc(ppn string) []byte
}

type HttpFetch struct{}
type HttpRequester func(string) ([]byte, error)

// FetchAll returns all xml data corresponding to each ppn in a map, or
// nil if unsuccessful.
// Note that the SUDOC API ignores unknown PPNs when requested with a
// muli-request.
func (f *HttpFetch) FetchAll(ppns []string) [][]byte {
	return fetchBatch(ppns, 20, Fetch)
}

// FetchRCR returns a XML iln2rcr response from a list of ILNs.
func (f *HttpFetch) FetchRCR(ilns []string) []byte {
	data, _ := Fetch(iln2rcr_url + strings.Join(ilns, ","))
	return data
}

func (f *HttpFetch) FetchMarc(ppn string) []byte {
	data, _ := Fetch(marcxml_url + ppn + ".xml")
	return data
}

func FetchMarc(ppn string) ([]byte, error) {
	return Fetch(marcxml_url + ppn + ".xml")
}

// Fetch returns the xml record corresponding to the given URL, or nil if
// unsucessful
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		// if request time out, just ignore
		// TODO: delay and request again
		// TODO: handle other url.errors
		log.Println(err)
		return []byte{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("fetch: HTTP status code = %d", resp.StatusCode)
		return []byte{}, errors.New(strconv.Itoa(resp.StatusCode))
	}
	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return data, nil
}

func fetchBatch(ppns []string, max_params int, request HttpRequester) [][]byte {
	var tokens = make(chan struct{}, MAX_CONCURRENT_REQUESTS)
	urls := buildURLs(ppns, max_params)
	xmlBatch := make([][]byte, len(urls))
	wg := sync.WaitGroup{}

	for index, url := range urls {
		wg.Add(1)
		go func(i int, u string) {
			tokens <- struct{}{}
			defer wg.Done()
			xmlBatch[i], _ = request(u)
			<-tokens
		}(index, url)
	}
	wg.Wait()
	return xmlBatch
}

func buildURLs(params []string, max_params int) []string {
	var urls []string
	for len(params) > max_params {
		newUrl := multiwhere_url + strings.Join(params[:max_params], ",")
		urls = append(urls, newUrl)
		params = params[max_params:]
	}
	urls = append(urls, multiwhere_url+strings.Join(params, ","))
	return urls
}
