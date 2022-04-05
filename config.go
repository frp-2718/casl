package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type config struct {
	MappingFilePath string   `json:"alma-rcr_file_path"`
	AlmaAPIKey      string   `json:"alma_api_key"`
	ILNs            []string `json:"iln_to_track"`
	IgnoredAlmaColl []string `json:"ignored_alma_collections"`
	IgnoredSudocRCR []string `json:"ignored_sudoc_rcr"`
}

func loadConfig() *config {
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
