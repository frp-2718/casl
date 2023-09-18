package main

import (
	"casl/bib"
	"casl/requests"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Casl struct {
	config     *Config
	mappings   *Mappings
	httpClient *requests.HttpFetch
}

type Config struct {
	MappingFilePath string   `json:"alma-rcr_file_path"`
	AlmaAPIKey      string   `json:"alma_api_key"`
	ILNs            []string `json:"iln_to_track"`
	IgnoredAlmaColl []string `json:"ignored_alma_collections"`
	IgnoredSudocRCR []string `json:"ignored_sudoc_rcr"`
	MonolithicRCR   []string `json:"monolithic_rcr"`
	FollowedRCR     []string
}

// Mappings Alma/RCR, Alma/Libraries names, RCR/ILN, read from CSV.
type Mappings struct {
	alma2rcr    map[string][]string
	rcr2iln     map[string]string
	alma2string map[string]string
}

func (c *Config) Print() {
	fmt.Printf("Mapping file: %s\n", c.MappingFilePath)
	fmt.Printf("Alma API key: %s\n", c.AlmaAPIKey)
	fmt.Printf("ILNs: %s\n", strings.Join(c.ILNs, ", "))
	fmt.Printf("Ignored Alma coll.: %s\n", strings.Join(c.IgnoredAlmaColl, ", "))
	fmt.Printf("Ignored SUDOC RCRs: %s\n", strings.Join(c.IgnoredSudocRCR, ", "))
	fmt.Printf("Monolithic RCRs: %s\n", strings.Join(c.MonolithicRCR, ", "))
}

func NewCasl() *Casl {
	var casl Casl

	casl.config = LoadConfig()
	casl.mappings = MakeMappings(casl.config)
	casl.httpClient = &requests.HttpFetch{}

	followed, err := bib.GetRCRs(casl.config.ILNs, casl.httpClient)
	if err != nil {
		log.Fatalf("casl: unable to fetch RCRs: %s", err)
	}
	casl.config.FollowedRCR = filter(followed, casl.config.IgnoredSudocRCR)
	return &casl
}

func (m *Mappings) Print() {
	fmt.Println("*** alma2rcr ***")
	for k, v := range m.alma2rcr {
		fmt.Println(k, v)
	}
	fmt.Println("*** rcr2 ***")
	for k, v := range m.rcr2iln {
		fmt.Println(k, v)
	}
	fmt.Println("*** alma2string ***")
	for k, v := range m.alma2string {
		fmt.Println(k, v)
	}
}

func LoadConfig() *Config {
	var conf Config
	content, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	return &conf
}

func MakeMappings(conf *Config) *Mappings {
	var maps Mappings
	maps.alma2rcr, maps.rcr2iln, maps.alma2string = csvToMap(conf.MappingFilePath)
	return &maps
}

func csvToMap(filename string) (map[string][]string, map[string]string, map[string]string) {
	almaRCR := make(map[string][]string)
	rcrILN := make(map[string]string)
	almaSTR := make(map[string]string)
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
		almaRCR[record[1]] = append(almaRCR[record[1]], record[2])
		if len(record[2]) > 0 {
			rcrILN[record[2]] = record[3]
		}
		almaSTR[record[1]] = record[0]
	}
	return almaRCR, rcrILN, almaSTR
}
