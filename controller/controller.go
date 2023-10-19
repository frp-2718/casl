package controller

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
)

func NewController(configFile string) Controller {
	var ctrl Controller

	ctrl.Config = loadConfig(configFile)
	ctrl.Mappings = csvToMap(ctrl.Config.MappingFilePath)
	// casl.httpClient = &requests.HttpFetch{}

	// followed, err := GetRCRs(casl.config.ILNs, casl.httpClient)
	// if err != nil {
	// 	log.Fatalf("casl: unable to fetch RCRs: %s", err)
	// }
	// casl.config.FollowedRCR = filter(followed, casl.config.IgnoredSudocRCR)
	return ctrl
}

func filter(a, b []string) []string {
	return a
}

func loadConfig(configFile string) *config {
	var conf config
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	return &conf
}

// // GetRCRs returns a slice of RCRs as strings from a slice of ILNs.
// func GetRCRs(ilns []string, client requests.Fetcher) ([]string, error) {
// 	xmldata := client.FetchRCR(ilns)
// 	result, err := decodeRCR(xmldata)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// type bibs struct {
// 	XMLName xml.Name `xml:"bibs"`
// 	Entries []bib    `xml:"bib"`
// }

// type bib struct {
// 	MMS       string `xml:"mms_id"`
// 	BibRecord record `xml:"record"`
// }

// type record struct {
// 	Fields []datafield `xml:"datafield"`
// }

// type datafield struct {
// 	Type      string     `xml:"tag,attr"`
// 	Subfields []subfield `xml:"subfield"`
// }

// type subfield struct {
// 	Code string `xml:"code,attr"`
// 	Data string `xml:",chardata"`
// }

// // structs matching ABES multiwhere XML format
// type response struct {
// 	Requests []query `xml:"query"`
// }

// type query struct {
// 	PPN   string   `xml:"ppn"`
// 	Items []result `xml:"result"`
// }

// type result struct {
// 	Libraries []library `xml:"library"`
// }

// type library struct {
// 	RCR  string `xml:"rcr"`
// 	Name string `xml:"shortname"`
// }

// // structs matching ABES iln2rcr XML format
// type sudoc struct {
// 	Queries []req `xml:"query"`
// }

// type req struct {
// 	ILN     string   `xml:"iln"`
// 	Results []result `xml:"result"`
// }

// func decodeRCR(xmldata []byte) ([]string, error) {
// 	var result sudoc
// 	var rcrs []string
// 	err := xml.Unmarshal(xmldata, &result)
// 	if err != nil {
// 		log.Printf("decodeRCR: %s", err)
// 		return nil, err
// 	}
// 	for _, query := range result.Queries {
// 		for _, r := range query.Results {
// 			for _, l := range r.Libraries {
// 				rcrs = append(rcrs, l.RCR)
// 			}
// 		}
// 	}
// 	return rcrs, nil
// }

func csvToMap(filename string) *mappings {
	var maps mappings
	maps.alma2rcr = make(map[string][]string)
	maps.rcr2iln = make(map[string]string)
	maps.alma2str = make(map[string]string)
	maps.rcr2str = make(map[string]string)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		maps.alma2rcr[record[1]] = append(maps.alma2rcr[record[1]], record[2])
		maps.rcr2iln[record[2]] = record[3]
		maps.alma2str[record[1]] = record[0]
		maps.rcr2str[record[2]] = record[4]
	}
	return &maps
}
