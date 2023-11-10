// Package alma provides a very simple adhoc Alma client.

// API thresholds : 200,000 requests/day and 25 requests/second.
package exl

import (
	"casl/entities"
	"casl/requests"
	"errors"
	"fmt"
	"slices"
)

// AlmaClient is the internal representation of the client
type AlmaClient struct {
	apiKey  string
	baseURL string
	stats   stats
	fetcher requests.Fetcher
}

type stats struct {
	bibs_req  int
	items_req int
}

const almawsURL = "https://api-eu.hosted.exlibrisgroup.com/almaws/v1/"

const (
	bibs_t int = iota
	items_t
)

// NewAlmaClient creates an Alma client with the default http client if none
// is provided.
func NewAlmaClient(apiKey, baseURL string, fetcher requests.Fetcher) (*AlmaClient, error) {
	alma := new(AlmaClient)
	if apiKey == "" {
		return nil, errors.New("NewAlmaClient: no API key provided")
	}
	alma.apiKey = apiKey
	if baseURL == "" {
		alma.baseURL = almawsURL
	} else {
		alma.baseURL = baseURL
	}

	if fetcher == nil {
		return nil, errors.New("NewAlmaClient: no HTTP client provided")
	}
	alma.fetcher = fetcher

	alma.stats = stats{}
	return alma, nil
}

// GetFilteredLocations gets Alma locations of a given PPN, properly filled,
// from the items API. Only the locations regarding the libraries of
// interest, given as a second argument, are provided.
func (a *AlmaClient) GetFilteredLocations(ppn string, lib_codes []string) ([]*entities.AlmaLocation, error) {
	locations, err := a.GetLocations(ppn)
	var filtered []*entities.AlmaLocation
	if err != nil {
		return filtered, err
	}

	for _, location := range locations {
		if slices.Contains(lib_codes, location.Library_code) && entities.ValidLocation(*location) {
			filtered = append(filtered, location)
		}
	}
	return filtered, nil
}

// GetLocations gets all the Alma locations of a given PPN, from the item  API,
// filled with data from client's mappings.
func (a *AlmaClient) GetLocations(ppn string) ([]*entities.AlmaLocation, error) {
	var res []*entities.AlmaLocation
	mms, err := a.getMMSfromPPN(ppn)
	if err != nil {
		return res, err
	}
	if len(mms) == 0 {
		return res, fmt.Errorf("GetAlmaLocation: PPN %s not found", ppn)
	} else if len(mms) > 1 {
		return res, fmt.Errorf("GetAlmaLocation: PPN %s found in %v", ppn, mms)
	}
	items, err := a.getItems(mms[0])
	items_by_mms := make(map[string][]Item)
	for _, item := range items {
		items_by_mms[item.Holding_data.MMS] = append(items_by_mms[item.Holding_data.MMS], item)
	}

	for _, v := range items_by_mms {
		var location entities.AlmaLocation
		var items []*entities.AlmaItem
		for _, item := range v {
			var almaItem entities.AlmaItem
			almaItem.Process_code = item.Details.Process.Code
			almaItem.Process_name = item.Details.Process.Name
			almaItem.Status = item.Details.Status.Code
			items = append(items, &almaItem)
		}
		location.Library_name = v[0].Details.Library.Name
		location.Library_code = v[0].Details.Library.Code
		location.Location_code = v[0].Details.Location.Code
		location.Location_name = v[0].Details.Location.Name
		location.Call_number = v[0].Holding_data.CallNumber
		location.NoDiscovery = v[0].Holding_data.Suppress_from_publishing
		location.Items = items
		res = append(res, &location)
	}

	return res, nil
}

// Stats returns numbers of requests made by the client to the service named
// by the argument ("bibs", "items", "total").
// TODO: provide a better way to select the stat than by string
func (a *AlmaClient) Stats(t string) int {
	switch t {
	case "bibs":
		return a.stats.bibs_req
	case "items":
		return a.stats.items_req
	case "total":
		return a.stats.bibs_req + a.stats.items_req
	default:
		return a.Stats("total")
	}
}

// getItems returns a list of all the items linked to the bibliographic record
// given as a parameter via its MMS. Alma API limits the number of retrieved
// items to 100.
func (a *AlmaClient) getItems(mms string) ([]Item, error) {
	a.stats.items_req += 1
	data, err := a.fetcher.Fetch(a.buildURL(items_t, mms))
	if err != nil {
		return nil, err
	}
	items, err := DecodeItemsXML(data)
	if err != nil {
		return nil, errors.New("alma: getItems: unable to decode XML data")
	}
	return items, nil
}

// getMMSfromPPN returns a list of MMS corresponding to the given PPN.
func (a *AlmaClient) getMMSfromPPN(ppn string) ([]string, error) {
	a.stats.bibs_req += 1
	data, err := a.fetcher.Fetch(a.buildURL(bibs_t, "(PPN)"+ppn))
	if err != nil { // HTTP errors
		return nil, err
	}
	bibs, err := decodeBibsXML(data)
	if err != nil {
		return nil, errors.New("alma: getMMSfromPPN: unable to decode XML data")
	}
	result := []string{}
	for _, bib := range bibs.Bibs {
		if slices.Contains(bib.Network_numbers, ppn) {
			result = append(result, bib.MMS_id)
		}
	}
	return result, nil
}

func (a *AlmaClient) buildURL(urlType int, id string) string {
	switch urlType {
	case bibs_t:
		return a.baseURL + "bibs?view=brief&expand=None&other_system_id=" + id + "&apikey=" + a.apiKey
	case items_t:
		return a.baseURL + "bibs/" + id + "/holdings/ALL/items?limit=100&apikey=" + a.apiKey
	default:
		return a.baseURL + "/" + id
	}
}
