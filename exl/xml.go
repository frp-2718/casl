package exl

import (
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

type almaBib struct {
	XMLName         xml.Name `xml:"bib"`
	Network_numbers []string `xml:"network_numbers>network_number"`
	MMS_id          string   `xml:"mms_id"`
}
type bibsResult struct {
	XMLName xml.Name  `xml:"bibs"`
	Bibs    []almaBib `xml:"bib"`
}

type holdingsResult struct {
	XMLName  xml.Name  `xml:"holdings"`
	Holdings []Holding `xml:"holding"`
}

type Holding struct {
	Suppress_from_publishing bool   `xml:"holding_suppress_from_publishing"`
	MMS                      string `xml:"holding_id"`
	CallNumber               string `xml:"call_number"`
}

type Item struct {
	XMLName      xml.Name `xml:"item"`
	Holding_data Holding  `xml:"holding_data"`
	Details      ItemData `xml:"item_data"`
}

type ItemData struct {
	Status   Status   `xml:"base_status"`
	Process  Process  `xml:"process_type"`
	Library  Library  `xml:"library"`
	Location Location `xml:"location"`
}

type Library struct {
	Name string `xml:"desc,attr"`
	Code string `xml:",chardata"`
}

type Location struct {
	Name string `xml:"desc,attr"`
	Code string `xml:",chardata"`
}

type Process struct {
	Name string `xml:"desc,attr"`
	Code string `xml:",chardata"`
}

type Status struct {
	Code   string `xml:"desc,attr"`
	Number string `xml:",chardata"`
}

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

func (i Item) String() string {
	return fmt.Sprintf("%s%s\n", i.Holding_data, i.Details)
}

func (h Holding) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Published: %t\n", h.Suppress_from_publishing)
	fmt.Fprintf(&sb, "Call number: %s\n", h.CallNumber)
	fmt.Fprintf(&sb, "MMS: %s\n", h.MMS)
	return sb.String()
}

func (i ItemData) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Status: %s\n", i.Status)
	fmt.Fprintf(&sb, "Library: %s\n", i.Library)
	fmt.Fprintf(&sb, "Location: %s\n", i.Location)
	fmt.Fprintf(&sb, "Process: %s\n", i.Process)
	return sb.String()
}

func decodeBibsXML(data []byte) (*bibsResult, error) {
	var b bibsResult
	err := xml.Unmarshal(data, &b)
	if err != nil {
		log.Printf("alma: GetMMSfromPPN: %v", err)
		return nil, err
	}
	return &b, nil
}

func decodeError(data []byte, status int) error {
	type almaError struct {
		ErrorCode    string `xml:"errorList>error>errorCode"`
		ErrorMessage string `xml:"errorList>error>errorMessage"`
	}
	var e almaError
	err := xml.Unmarshal(data, &e)
	if err != nil {
		log.Printf("decodeError: %s", err)
		return &FetchError{errorMessage: ""}
	}
	switch status {
	case 400:
		if e.ErrorCode == "GENERAL_ERROR" || e.ErrorCode == "401652" || e.ErrorCode == "402203" {
			return &InvalidRequestError{errorMessage: e.ErrorMessage}
		}
		if e.ErrorCode == "UNAUTHORIZED" {
			return &UnauthorizedError{errorMessage: e.ErrorMessage}
		}
	case 403:
		if e.ErrorCode == "UNAUTHORIZED" || e.ErrorCode == "INVALID_REQUEST" || e.ErrorCode == "FORBIDDEN" {
			return &UnauthorizedError{errorMessage: e.ErrorMessage}
		}
		if e.ErrorCode == "REQUEST_TOO_LARGE" {
			return &InvalidRequestError{errorMessage: e.ErrorMessage}
		}
	case 429:
		if e.ErrorCode == "PER_SECOND_THRESHOLD" || e.ErrorCode == "DAILY_THRESHOLD" {
			return &ThresholdError{errorMessage: e.ErrorCode}
		}
	case 500:
		if e.ErrorCode == "GENERAL_ERROR" {
			return &UnauthorizedError{errorMessage: e.ErrorMessage}
		}
	case 503:
		if e.ErrorCode == "ROUTING_ERROR" {
			return &ServerError{errorMessage: e.ErrorMessage}
		}
	}
	return &FetchError{errorMessage: e.ErrorMessage}
}

func notFound(data []byte) bool {
	type almaResponse struct {
		Count string `xml:"total_record_count,attr"`
	}
	var r almaResponse
	err := xml.Unmarshal(data, &r)
	if err != nil {
		log.Printf("notFound: unmarshal error: %s", err)
		return true
	}
	return r.Count == "0" // item not found
}

func decodeHoldingsXML(data []byte) ([]Holding, error) {
	var h holdingsResult
	h.Holdings = []Holding{} // should be initialized in case of there is 0 holdings in the data
	err := xml.Unmarshal(data, &h)
	if err != nil {
		log.Printf("alma: decodeHoldingsXML: %v", err)
		return nil, err
	}
	return h.Holdings, nil
}

func DecodeItemsXML(data []byte) ([]Item, error) {
	var items Items
	items.Items = []Item{}
	err := xml.Unmarshal(data, &items)
	if err != nil {
		log.Printf("alma: decodeItemsXML: %v", err)
		return nil, err
	}
	return items.Items, nil
}
