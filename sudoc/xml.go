package sudoc

import (
	"encoding/xml"
	"errors"
	"log"
)

type iln2rcr_response struct {
	XMLName xml.Name        `xml:"sudoc"`
	Queries []iln2rcr_query `xml:"query"`
}

type iln2rcr_query struct {
	XMLName   xml.Name          `xml:"query"`
	ILN       string            `xml:"iln"`
	Libraries []iln2rcr_library `xml:"result>library"`
}

type iln2rcr_library struct {
	XMLName xml.Name `xml:"library"`
	RCR     string   `xml:"rcr"`
	Name    string   `xml:"shortname"`
}

func decodeRCR(data []byte) (map[string]library, error) {
	mapping := make(map[string]library)
	var result iln2rcr_response
	err := xml.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
	}
	// iln not found
	if len(result.Queries) == 0 {
		return nil, errors.New("null xml")
	}
	for _, q := range result.Queries {
		for _, lib := range q.Libraries {
			mapping[lib.RCR] = library{q.ILN, lib.RCR, lib.Name}
		}
	}
	return mapping, nil
}
