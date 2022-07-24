// Package alma provides a very simple adhoc Alma client.
//
// API thresholds : 200,000 requests/day and 25 requests/second.
package alma

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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
		// TODO: set default client timeout
		client = http.DefaultClient
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
		// TODO: retry timeout queries
		return nil, errors.New("alma: fetch: http error")
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("alma: fetch: read error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, decodeError(data, resp.StatusCode)
	} else if notFound(data) {
		// TODO: add not found id to error
		return nil, &NotFoundError{errorMessage: "identifier not found"}
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
