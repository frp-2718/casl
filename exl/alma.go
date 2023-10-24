// Package alma provides a very simple adhoc Alma client.
//
// API thresholds : 200,000 requests/day and 25 requests/second.
package exl

import (
	"casl/entities"
	"casl/requests"
	"errors"
	"fmt"
	"log"
)

type AlmaClient struct {
	apiKey  string
	baseURL string
}

const almawsURL = "https://api-eu.hosted.exlibrisgroup.com/almaws/v1/"

const (
	holdings_t int = iota
	bibs_t
	items_t
)

// New creates an Alma client with the default http client if none is provided.
func NewAlmaClient(apiKey, baseURL string) *AlmaClient {
	alma := new(AlmaClient)
	alma.apiKey = apiKey
	if baseURL == "" {
		alma.baseURL = almawsURL
	} else {
		alma.baseURL = baseURL
	}
	return alma
}

// GetMMSfromPPN returns a list of MMS corresponding to the given PPN, or
// NotFoundError.
func (a *AlmaClient) GetMMSfromPPN(ppn string) ([]string, error) {
	data, err := requests.Fetch(a.buildURL(bibs_t, "(PPN)"+ppn))
	if err != nil { // HTTP errors, including NotFoundError
		// log.Printf("alma: GetMMSfromPPN: %v", err)
		return nil, err
	}
	bibs, err := decodeBibsXML(data)
	if err != nil {
		log.Printf("alma: GetMMSfromPPN: %v", err)
		return nil, errors.New("alma: GetMMSfromPPN: unable to decode XML data")
	}
	result := []string{}
	for _, bib := range bibs.Bibs {
		// if ppnMatch(bib.Network_numbers, ppn) {
		// 	result = append(result, bib.MMS_id)
		// }
		result = append(result, bib.MMS_id)
	}
	return result, nil
}

// max 100 items
func (a *AlmaClient) GetItems(mms string) ([]Item, error) {
	data, err := requests.Fetch(a.buildURL(items_t, mms))
	if err != nil {
		// log.Printf("alma: GetHoldings: %v", err)
		return nil, err
	}
	items, err := DecodeItemsXML(data)
	if err != nil {
		log.Printf("alma: GetItems: %v", err)
		return nil, errors.New("alma: GetItems: unable to decode XML data")
	}
	return items, nil
}

func (a *AlmaClient) DecodeItemsXML(data []byte) ([]Item, error) {
	return DecodeItemsXML(data)
}

// GetHoldings returns a list of holdings, potentially empty, attached to the
// given bib MMS id.
func (a *AlmaClient) GetHoldings(mms string) ([]Holding, error) {
	data, err := requests.Fetch(a.buildURL(holdings_t, string(mms)))
	if err != nil {
		// log.Printf("alma: GetHoldings: %v", err)
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
func (a *AlmaClient) GetHoldingsFromPPN(ppn string) ([]Holding, error) {
	ids, err := a.GetMMSfromPPN(ppn)
	if err != nil {
		// log.Printf("alma: GetHoldingsFromPPN: %v", err)
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

func (a *AlmaClient) GetAlmaLocation(ppn string) ([]*entities.AlmaLocation, error) {
	var res []*entities.AlmaLocation
	mms, err := a.GetMMSfromPPN(ppn)
	if err != nil {
		return res, err
	}
	if len(mms) == 0 {
		return res, fmt.Errorf("GetAlmaLocation: PPN %s not found", ppn)
	}
	// TODO: manage the multi-MMS case
	items, err := a.GetItems(mms[0])
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
		location.Published = v[0].Holding_data.Suppress_from_publishing
		location.Items = items
		res = append(res, &location)
	}

	return res, nil
}

func (a *AlmaClient) buildURL(urlType int, id string) string {
	switch urlType {
	case bibs_t:
		return a.baseURL + "bibs?view=brief&expand=None&other_system_id=" + id + "&apikey=" + a.apiKey
	case holdings_t:
		return a.baseURL + "bibs/" + id + "/holdings?apikey=" + a.apiKey
	case items_t:
		return a.baseURL + "bibs/" + id + "/holdings/ALL/items?limit=100&apikey=" + a.apiKey
	default:
		return a.baseURL + "/" + id
	}
}

func ppnMatch(ids []string, ppn string) bool {
	for _, id := range ids {
		if id == "(PPN)"+string(ppn) {
			return true
		}
	}
	return false
}
