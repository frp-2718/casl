package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type config struct {
	MappingFilePath string   `json:"alma-rcr_file_path"`
	AlmaAPIKey      string   `json:"alma_api_key"`
	ILNs            []string `json:"iln_to_track"`
	IgnoredAlmaColl []string `json:"ignored_alma_collections"`
	IgnoredSudocRCR []string `json:"ignored_sudoc_rcr"`
	MonolithicRCR   []string `json:"monolithic_rcr"`
}

func LoadConfig() *config {
	var conf config
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	return &conf
}

func (c *config) Print() {
	fmt.Printf("Mapping file:%s\n", c.MappingFilePath)
	fmt.Printf("Alma API key:%s\n", c.AlmaAPIKey)
	fmt.Printf("ILNs: %s\n", strings.Join(c.ILNs, ", "))
	fmt.Printf("Ignored Alma coll.: %s\n", strings.Join(c.IgnoredAlmaColl, ", "))
	fmt.Printf("Ignored SUDOC RCRs: %s\n", strings.Join(c.IgnoredSudocRCR, ", "))
	fmt.Printf("Monolithic RCRs: %s\n", strings.Join(c.MonolithicRCR, ", "))
}
