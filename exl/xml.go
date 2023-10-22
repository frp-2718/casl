package exl

import (
	"encoding/xml"
	"fmt"
	"log"
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
	Library    string `xml:"library"`
	Location   string `xml:"location"`
	CallNumber string `xml:"call_number"`
}

func (h Holding) String() string {
	return fmt.Sprintf("---\nLibrary: %s\nLocation: %s\nCallNumber: %s\n",
		h.Library, h.Location, h.CallNumber)
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
