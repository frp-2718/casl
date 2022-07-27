// casl is a tool to compare locations of bibliographic resources between SUDOC
// and Alma.
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"casl/alma"
	"casl/bib"
	"casl/requests"

	"golang.org/x/exp/slices"
)

// General configuration.
var conf *config

// Mappings Alma/RCR, Alma/Libraries names, RCR/ILN, read from CSV.
var alma2rcr map[string]string
var rcr2iln map[string]string
var alma2string map[string]string

// Filters.
var followedRCR []string
var monolithicRCR []string

var httpFetcher requests.HttpFetch

func init() {
	conf = loadConfig()
	alma2rcr, rcr2iln, alma2string = csvToMap(conf.MappingFilePath)
	httpFetcher := requests.HttpFetch{}
	followed, err := bib.GetRCRs(conf.ILNs, &httpFetcher)
	if err != nil {
		log.Fatalf("casl: unable to fetch RCRs: %s", err)
	}
	followedRCR = filter(followed, conf.IgnoredSudocRCR)
	monolithicRCR = conf.MonolithicRCR
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: casl file1 file2...")
		log.Fatal("casl: called without arguments")
	}

	start := time.Now()

	// PPNs to check.
	ppns := make(map[string]bool)

	for _, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)

		// Assume one valid ppn per line.
		// TODO: add validation test
		for scanner.Scan() {
			ppns[scanner.Text()] = true
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%d PPN à vérifier...\n", len(ppns))

	// TODO: concurrent pipeline
	// SUDOC processing
	records := bib.GetSudocLocations(ppns, followedRCR, &httpFetcher)

	// Alma processing
	almaClient, err := alma.New(nil, conf.AlmaAPIKey, "")
	if err != nil {
		log.Fatalf("GetAlmaLocations: unable to initialize the Alma client: %v", err)
	}
	resultats := bib.GetAlmaLocations(almaClient, records, alma2rcr)

	// End results comparison
	anomalies := bib.Filter(bib.ComparePPN(resultats, conf.IgnoredAlmaColl), monolithicRCR, &httpFetcher)

	writeCSV(anomalies)

	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)
}

func writeCSV(results []bib.CRecord) {
	var records [][]string
	records = append(records, []string{"PPN", "ILN", "Bibliothèque Alma",
		"Bibliothèque SUDOC", "RCR"})
	sort.Slice(results, func(i, j int) bool {
		if results[i].ILN != results[j].ILN {
			return results[i].ILN < results[j].ILN
		}
		if results[i].RCR != results[j].ILN {
			return results[i].RCR < results[j].RCR
		}
		return results[i].PPN < results[j].PPN
	})
	for _, res := range results {
		var suffix string
		if res.SUDOCSublocation != "" {
			suffix = " - " + res.SUDOCSublocation
		}
		record := []string{res.PPN, rcr2iln[res.RCR], alma2string[res.AlmaLibrary],
			res.SUDOCLibrary + suffix, res.RCR}
		records = append(records, record)
	}

	filename := "resultats_" + strconv.FormatInt(time.Now().Unix(), 10) + ".csv"
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		log.Fatal("failed to open file", err)
	}

	w := csv.NewWriter(f)
	err = w.WriteAll(records)

	if err != nil {
		log.Fatal(err)
	}
}

func csvToMap(filename string) (map[string]string, map[string]string, map[string]string) {
	almaRCR := make(map[string]string)
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
		almaRCR[record[1]] = record[2]
		if len(record[2]) > 0 {
			rcrILN[record[2]] = record[3]
		}
		almaSTR[record[1]] = record[0]
	}
	return almaRCR, rcrILN, almaSTR
}

func filter(s, ignored []string) []string {
	var result []string
	for _, elem := range s {
		if !slices.Contains(ignored, elem) {
			result = append(result, elem)
		}
	}
	return result
}
