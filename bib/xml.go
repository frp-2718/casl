package bib

import (
	"encoding/xml"
	"log"
)

// structs for Alma XML 'bibs' response
type bibs struct {
	XMLName xml.Name `xml:"bibs"`
	Entries []bib    `xml:"bib"`
}

type bib struct {
	MMS       string `xml:"mms_id"`
	BibRecord record `xml:"record"`
}

type record struct {
	Fields []datafield `xml:"datafield"`
}

type datafield struct {
	Type      string     `xml:"tag,attr"`
	Subfields []subfield `xml:"subfield"`
}

type subfield struct {
	Code string `xml:"code,attr"`
	Data string `xml:",chardata"`
}

// structs matching ABES multiwhere XML format
type response struct {
	Requests []query `xml:"query"`
}

type query struct {
	PPN   string   `xml:"ppn"`
	Items []result `xml:"result"`
}

type result struct {
	Libraries []library `xml:"library"`
}

type library struct {
	RCR  string `xml:"rcr"`
	Name string `xml:"shortname"`
}

// structs matching ABES iln2rcr XML format
type sudoc struct {
	Queries []req `xml:"query"`
}

type req struct {
	ILN     string   `xml:"iln"`
	Results []result `xml:"result"`
}

func decodeAlmaXML(xmldata []byte) ([]almaLocation, error) {
	var result bibs
	var locations []almaLocation
	err := xml.Unmarshal(xmldata, &result)
	if err != nil {
		log.Printf("decodeAlmaXML: %s", err)
		return nil, err
	}
	if len(result.Entries) > 0 { // exists in Alma
		for _, datafield := range result.Entries[0].BibRecord.Fields {
			if datafield.Type == "AVA" { // AVA is a holding record
				var loc almaLocation
				for _, subfield := range datafield.Subfields {
					if subfield.Code == "b" { // library code
						loc.ownerCode = subfield.Data
					}
					if subfield.Code == "j" { // location code
						loc.collection = subfield.Data
					}
				}
				locations = append(locations, loc)
			}
		}
	}
	return locations, nil // possibly nil if there is no location in Alma
}

func decodeLocations(xmldata []byte, rcrs []string) ([]BibRecord, error) {
	var result response
	var records []BibRecord
	err := xml.Unmarshal(xmldata, &result)
	if err != nil {
		log.Printf("decodeLocations: %s", err)
		return nil, err
	}
	for _, query := range result.Requests {
		record := BibRecord{ppn: query.PPN}
		locations := []sudocLocation{}
		for _, item := range query.Items {
			for _, library := range item.Libraries {
				if in(library.RCR, rcrs) {
					locations = append(locations, sudocLocation{rcr: library.RCR, name: library.Name})
				}
			}
		}
		record.sudocLocations = locations
		records = append(records, record)
	}
	return records, nil
}

func decodeRCR(xmldata []byte) ([]string, error) {
	var result sudoc
	var rcrs []string
	err := xml.Unmarshal(xmldata, &result)
	if err != nil {
		log.Printf("decodeRCR: %s", err)
		return nil, err
	}
	for _, query := range result.Queries {
		for _, r := range query.Results {
			for _, l := range r.Libraries {
				rcrs = append(rcrs, l.RCR)
			}
		}
	}
	return rcrs, nil
}
