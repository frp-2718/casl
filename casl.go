// casl is a tool to compare locations of bibliographic resources between SUDOC
// and Alma.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"time"

	"casl/bib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: casl file1 file2...")
		log.Fatal("casl: called without arguments")
	}

	casl := NewCasl()

	start := time.Now()

	// PPNs to check.
	ppns := make(map[string]bool)

	var records []bib.BibRecord

	for _, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)

		// Discard invalid ppns.
		var ppnPattern = regexp.MustCompile(`[0-9]{8}([0-9]|(x|X))`)
		for scanner.Scan() {
			if line := scanner.Text(); ppnPattern.Match([]byte(line)) {
				ppns[line] = true
				records = append(records, bib.BibRecord{PPN: line})
			} else {
				log.Printf("invalid PPN: %s", line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%d PPN à vérifier...\n", len(ppns))

	fmt.Println(casl.config.FollowedRCR)

	// 	// TODO: concurrent pipeline
	// 	// SUDOC processing
	// 	records = bib.GetSudocLocations(ppns, casl.config.FollowedRCR, casl.httpClient)

	// 	// Alma processing
	// 	almaClient, err := alma.New(nil, casl.config.AlmaAPIKey, "")
	// 	if err != nil {
	// 		log.Fatalf("GetAlmaLocations: unable to initialize the Alma client: %v", err)
	// 	}
	// 	resultats := bib.GetAlmaLocations(almaClient, records, casl.mappings.alma2rcr)

	// 	// End results comparison
	// 	anomalies := bib.Filter(bib.ComparePPN(resultats, casl.config.IgnoredAlmaColl),
	// 		casl.config.MonolithicRCR, casl.httpClient)

	// 	writeCSV(anomalies)

	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)
	// }

	// func writeCSV(results []bib.CRecord) {
	// 	var records [][]string
	// 	records = append(records, []string{"PPN", "ILN", "Bibliothèque Alma",
	// 		"Bibliothèque SUDOC", "RCR"})
	// 	sort.Slice(results, func(i, j int) bool {
	// 		if results[i].ILN != results[j].ILN {
	// 			return results[i].ILN < results[j].ILN
	// 		}
	// 		// if results[i].RCR != results[j].ILN {
	// 		// 	return results[i].RCR < results[j].RCR
	// 		// }
	// 		return results[i].PPN < results[j].PPN
	// 	})
	// 	// for _, res := range results {
	// 	// 	// 	suffix = " - " + res.SUDOCSublocation
	// 	// 	// }
	// 	// 	// record := []string{res.PPN, rcr2iln[res.RCR], alma2string[res.AlmaLibrary],
	// 	// 	// 	res.SUDOCLibrary + suffix, res.RCR}
	// 	// 	// records = append(records, record)
	// 	// }

	// t := time.Now()
	// format := fmt.Sprintf("%d%02d%02d-%02d%02d%02d", t.Year(), t.Month(), t.Day(),
	// 	t.Hour(), t.Minute(), t.Second())
	// filename := "resultats_" + format + ".csv"
	// f, err := os.Create(filename)
	// defer f.Close()

	// if err != nil {
	// 	log.Fatal("failed to open file", err)
	// }

	// w := csv.NewWriter(f)
	// err = w.WriteAll(records)

	// if err != nil {
	// 	log.Fatal(err)
	// }
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
