package controller

import (
	"casl/entities"
	"casl/exl"
	"casl/requests"
	"casl/sudoc"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"
)

// NewController creates a fully self-configured controller, which is the entry
// point of the process.
func NewController(configFile string, fetcher requests.Fetcher) (Controller, error) {
	var ctrl Controller

	ctrl.loadConfig(configFile)
	ctrl.getMappingsFromCSV(ctrl.Config.MappingFilePath)
	ctrl.getLibs()

	suclient, err := sudoc.NewSudocClient(ctrl.Config.ILNs, fetcher)
	if err != nil {
		return ctrl, err
	}
	ctrl.SUClient = suclient
	ctrl.Config.FollowedRCR = ctrl.SUClient.GetFollowedRCRs()

	ctrl.AlmaClient, err = exl.NewAlmaClient(ctrl.Config.AlmaAPIKey, "", fetcher)
	if err != nil {
		return ctrl, err
	}

	return ctrl, nil
}

func (ctrl *Controller) loadConfig(configFile string) {
	var conf config
	content, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatalf("loadConfig: %s", err)
	}

	ctrl.Config = &conf
}

func (ctrl *Controller) getLibs() {
	libs := make([]string, 0, len(ctrl.Mappings.alma2str))

	for k := range ctrl.Mappings.alma2str {
		libs = append(libs, k)
	}

	ctrl.Config.FolowedLibs = libs
}

// The CSV alma-rcr mapping should be formatted as follows :
// "Library name","Library code",RCR,ILN
func (ctrl *Controller) getMappingsFromCSV(csv_file string) {
	var maps mappings
	maps.alma2rcr = make(map[string][]string)
	maps.rcr2alma = make(map[string][]string)
	maps.rcr2iln = make(map[string]string)
	maps.alma2str = make(map[string]string)
	maps.rcr2str = make(map[string]string)
	f, err := os.Open(csv_file)
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
		maps.rcr2alma[record[2]] = append(maps.rcr2alma[record[2]], record[1])
		maps.rcr2iln[record[2]] = record[3]
		maps.alma2str[record[1]] = record[0]
	}
	ctrl.Mappings = &maps
}

// Summary represents the informations necessary to identify an anomaly, ie a
// record for which alma locations and sudoc locations are not matching.
type Summary struct {
	ILN      string
	RCR      string
	PPN      string
	SudocLib string
	AlmaLib  string
}

// Compare looks for anomalies - ie locations not maching - in the provided
// bib records.
func (ctrl *Controller) Compare(record *entities.BibRecord) []Summary {
	var anomalies []Summary

MAIN_SU_LOOP:
	for _, sloc := range record.SudocLocations {
		almaLibs := ctrl.Mappings.rcr2alma[sloc.RCR]
		for _, aloc := range record.AlmaLocations {
			if slices.Contains(almaLibs, aloc.Library_code) {
				continue MAIN_SU_LOOP
			}
		}
		// TODO: test that
		//library := ctrl.Mappings.rcr2str[sloc.RCR]
		library := sloc.Name
		if slices.Contains(ctrl.Config.MonolithicRCR, sloc.RCR) && sloc.Sublocation != "" {
			library += " - " + sloc.Sublocation
		}
		anomalies = append(anomalies, Summary{ILN: sloc.ILN, RCR: sloc.RCR, PPN: record.PPN, SudocLib: library, AlmaLib: ""})
	}

MAIN_ALMA_LOOP:
	for _, aloc := range record.AlmaLocations {
		rcrs := ctrl.Mappings.alma2rcr[aloc.Library_code]
		for _, sloc := range record.SudocLocations {
			if slices.Contains(rcrs, sloc.RCR) {
				continue MAIN_ALMA_LOOP
			}
		}
		anomalies = append(anomalies, Summary{ILN: ctrl.Mappings.rcr2iln[rcrs[0]], RCR: rcrs[0], PPN: record.PPN, SudocLib: "", AlmaLib: ctrl.Mappings.alma2str[aloc.Library_code]})
	}

	return anomalies
}

func (s Summary) toCSV() []string {
	records := []string{s.PPN, s.ILN, s.AlmaLib, s.SudocLib, s.RCR}
	return records
}

// WriteCSV translates a list of Summaries into a CSV file.
func (ctrl *Controller) WriteCSV(results []Summary) {
	var records [][]string
	records = append(records, []string{"PPN", "ILN", "Bibliothèque Alma",
		"Bibliothèque SUDOC", "RCR"})

	for _, res := range results {
		records = append(records, res.toCSV())
	}

	t := time.Now()
	format := fmt.Sprintf("%d%02d%02d-%02d%02d%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	filename := "resultats_" + format + ".csv"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("failed to open file", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		log.Fatal(err)
	}
}
