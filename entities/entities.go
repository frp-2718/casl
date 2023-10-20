package entities

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type BibRecord struct {
	PPN            string
	MMS            string
	SudocLocations []*SudocLocation
	AlmaLocations  []*AlmaLocation
}

type SudocLocation struct {
	ILN         string
	RCR         string
	Name        string
	Sublocation string
}

type AlmaLocation struct {
	Collection string
	OwnerCode  string
	RCR        []string
}

// TODO: complete the string representation of a BibRecord
func (r BibRecord) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*** PPN: %s\n\n", r.PPN)
	for _, sl := range r.SudocLocations {
		fmt.Fprintf(&sb, "%s----------\n", sl)
	}
	return sb.String()
}

func (s SudocLocation) String() string {
	return fmt.Sprintf("ILN: %s\nRCR: %s\nNAME: %s\nSUBLOCATION: %s\n",
		s.ILN, s.RCR, s.Name, s.Sublocation)
}

func (r BibRecord) toCSV() [][]string {
	var records [][]string
	for _, sudocLoc := range r.SudocLocations {
		records = append(records, []string{r.PPN, sudocLoc.ILN, "", sudocLoc.Name + " - " + sudocLoc.Sublocation, sudocLoc.RCR})
	}
	for _, almaLoc := range r.AlmaLocations {
		records = append(records, []string{r.PPN, "", almaLoc.OwnerCode, "", ""})
	}
	return records
}

func WriteCSV(results []BibRecord) {
	var records [][]string
	records = append(records, []string{"PPN", "ILN", "Bibliothèque Alma",
		"Bibliothèque SUDOC", "RCR"})

	for _, res := range results {
		records = append(records, res.toCSV()...)
	}

	t := time.Now()
	format := fmt.Sprintf("%d%02d%02d-%02d%02d%02d", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	filename := "resultats_" + format + ".csv"
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
