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
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: casl file1 file2...")
		log.Fatal("casl: called without arguments")
	}

	start := time.Now()

	ctrl := controller.NewController("config.json")

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
	for _, record := range records {
		locs, err := ctrl.SUClient.GetFilteredLocations(record.PPN, ctrl.Config.FollowedRCR)
		if err != nil {
			continue
		}
		if len(locs) > 0 {
			record.SudocLocations = locs
			results = append(results, record)
		}
	}

	entities.WriteCSV(results)

	elapsed := time.Since(start)
	fmt.Printf("Elapsed time: %s\n", elapsed)
}
