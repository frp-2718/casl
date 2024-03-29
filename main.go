// casl is a tool to compare locations of bibliographic resources between SUDOC
// and Alma.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"casl/controller"
	"casl/entities"
	"casl/requests"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: casl file1 file2...")
		log.Fatal("casl: called without arguments")
	}

	start := time.Now()

	fetcher := requests.NewHttpFetch(nil)
	ctrl, err := controller.NewController("config.json", fetcher)
	if err != nil {
		log.Fatal(err)
	}

	// PPNs to check.
	ppns := make(map[string]bool)
	var records []entities.BibRecord
	var ppnPattern = regexp.MustCompile(`[0-9]{8}([0-9]|(x|X))`)

	for _, filename := range os.Args[1:] {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)

		// Discard invalid ppns.
		for scanner.Scan() {
			if line := scanner.Text(); ppnPattern.Match([]byte(line)) {
				ppns[line] = true
				records = append(records, entities.BibRecord{PPN: line})
			} else {
				log.Printf("invalid PPN: %s", line)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%d PPN à vérifier...\n", len(ppns))

	var results []entities.BibRecord
	for i, record := range records {
		fmt.Printf("ppn %d/%d...\n", i, len(records))
		sudoc, err := ctrl.SUClient.GetFilteredLocations(record.PPN, ctrl.Config.FollowedRCR)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(sudoc) > 0 {
			record.SudocLocations = sudoc
		}
		alma, err := ctrl.AlmaClient.GetFilteredLocations(record.PPN, ctrl.Config.FolowedLibs, ctrl.Config.IgnoredAlmaColl)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(alma) > 0 {
			record.AlmaLocations = alma
		}
		results = append(results, record)
	}

	var sums []controller.Summary
	for _, res := range results {
		sums = append(sums, ctrl.Compare(&res)...)
	}

	ctrl.WriteCSV(sums)

	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)

	fmt.Println()
	fmt.Println("ALMA STATS")
	fmt.Printf("bibs: %d\n", ctrl.AlmaClient.Stats("bibs"))
	fmt.Printf("items: %d\n", ctrl.AlmaClient.Stats("items"))
	fmt.Printf("total: %d\n", ctrl.AlmaClient.Stats("total"))
	fmt.Println()
	fmt.Println("SUDOC STATS")
	fmt.Printf("iln2rcr: %d\n", ctrl.SUClient.Stats("iln2rcr"))
	fmt.Printf("marcxml: %d\n", ctrl.SUClient.Stats("marcxml"))
	fmt.Printf("total: %d\n", ctrl.SUClient.Stats("total"))
}
