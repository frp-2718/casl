// Package alma provides a very simple adhoc Alma client.
//
// API thresholds : 200,000 requests/day and 25 requests/second.
package alma

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Alma struct {
	apiKey      string
	baseURL     string
	client      *http.Client
	fetchClient Fetcher
}

type PPN string
type MMS string

const almawsURL = "https://api-eu.hosted.exlibrisgroup.com/almaws/v1/"

const (
	holdings_t int = iota
	bibs_t
)

// New creates an Alma client with the default http client if none is provided.
func New(client *http.Client, apiKey, baseURL string) (*Alma, error) {
	if apiKey == "" {
		return nil, errors.New("New: illegal argument")
	}
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	alma := new(Alma)
	alma.client = client
	alma.apiKey = apiKey
	if baseURL == "" {
		alma.baseURL = almawsURL
	} else {
		alma.baseURL = baseURL
	}
	alma.fetchClient = &almaFetcher{client: client}
	return alma, nil
}

// SetFetcher allows user to provide a custom Fetcher, for testing
// purposes.
func (a *Alma) SetFetcher(f Fetcher) {
	a.fetchClient = f
}

// GetMMSfromPPN returns a list of MMS corresponding to the given PPN, or
// NotFoundError.
func (a *Alma) GetMMSfromPPN(ppn PPN) ([]MMS, error) {
	data, err := a.Fetch(a.buildURL(bibs_t, string(ppn)))
	if err != nil { // HTTP errors, including NotFoundError
		log.Printf("alma: GetMMSfromPPN: %v", err)
		return nil, err
	}
	bibs, err := decodeBibsXML(data)
	if err != nil {
		log.Printf("alma: GetMMSfromPPN: %v", err)
		return nil, errors.New("alma: GetMMSfromPPN: unable to decode XML data")
	}
	result := []MMS{}
	for _, bib := range bibs.Bibs {
		if ppnMatch(bib.Network_numbers, ppn) {
			result = append(result, bib.MMS_id)
		}
	}
	return result, nil
}

// GetHoldings returns a list of holdings, potentially empty, attached to the
// given bib MMS id.
func (a *Alma) GetHoldings(mms MMS) ([]Holding, error) {
	data, err := a.Fetch(a.buildURL(holdings_t, string(mms)))
	if err != nil {
		log.Printf("alma: GetHoldings: %v", err)
		return nil, err
	}
	holdings, err := decodeHoldingsXML(data)
	if err != nil {
		log.Printf("alma: GetHoldings: %v", err)
		return nil, errors.New("alma: GetHoldings: unable to decode XML data")
	}
	return holdings, nil
}

// GetHoldingsFromPPN returns a list of holding, potentially empty, related to
// the bib record containing the given PPN in its network numbers.
func (a *Alma) GetHoldingsFromPPN(ppn PPN) ([]Holding, error) {
	ids, err := a.GetMMSfromPPN(ppn)
	if err != nil {
		log.Printf("alma: GetHoldingsFromPPN: %v", err)
		return nil, err
	}
	var result []Holding
	for _, mms := range ids {
		h, err := a.GetHoldings(mms)
		if err != nil {
			return nil, err
		}
		result = append(result, h...)
	}
	return result, nil
}

func (a *Alma) buildURL(urlType int, id string) string {
	switch urlType {
	case bibs_t:
		return a.baseURL + "bibs?view=brief&expand=None&other_system_id=" + id + "&apikey=" + a.apiKey
	case holdings_t:
		return a.baseURL + "bibs/" + id + "/holdings?apikey=" + a.apiKey
	default:
		return a.baseURL + "/" + id
	}
}

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type almaFetcher struct {
	client *http.Client
}

func (f *almaFetcher) Fetch(url string) ([]byte, error) {
	if url == "" {
		return nil, &InvalidRequestError{errorMessage: "URL cannot be empty"}
	}
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("alma>fetch: http error: %w", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("alma>fetch: read error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		err := decodeError(data, resp.StatusCode)
		// TODO: optimize number of concurrent requests to minimize the number
		// of 429
		if resp.StatusCode == 429 && err.Error() == "PER_SECOND_THRESHOLD" {
			log.Printf("alma>fetch: HTTP 429 PER_SECOND_THRESHOLD: %s", url)
			time.Sleep(1 * time.Second)
			return f.Fetch(url)
		}
		return nil, err
	} else if notFound(data) {
		return nil, &NotFoundError{id: url, errorMessage: "identifier not found"}
	}
	return data, nil
}

func (a *Alma) Fetch(url string) ([]byte, error) {
	return a.fetchClient.Fetch(url)
}

func ppnMatch(ids []string, ppn PPN) bool {
	for _, id := range ids {
		if id == "(PPN)"+string(ppn) {
			return true
		}
	}
	return false
}
