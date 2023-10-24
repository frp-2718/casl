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
	Library_name  string
	Library_code  string
	Location_name string
	Location_code string
	Call_number   string
	NoDiscovery   bool
	Items         []*AlmaItem
}

type AlmaItem struct {
	Process_name string
	Process_code string
	Status       string
}

// TODO: complete the string representation of a BibRecord
func (r BibRecord) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*** PPN: %s\n\n", r.PPN)
	for _, sl := range r.SudocLocations {
		fmt.Fprintf(&sb, "%s\n", sl)
	}
	for _, al := range r.AlmaLocations {
		fmt.Fprintf(&sb, "%s\n", al)
	}
	return sb.String()
}

func (s SudocLocation) String() string {
	return fmt.Sprintf("ILN: %s\nRCR: %s\nNAME: %s\nSUBLOCATION: %s\n",
		s.ILN, s.RCR, s.Name, s.Sublocation)
}

func (a AlmaLocation) String() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, "*********************************")
	fmt.Fprintf(&sb, "Library: %s (%s)\n", a.Library_name, a.Library_code)
	fmt.Fprintf(&sb, "Location: %s (%s)\n", a.Location_name, a.Location_code)
	fmt.Fprintf(&sb, "Call number: %s\n", a.Call_number)
	fmt.Fprintf(&sb, "Suppressed from discovery: %t\n", a.NoDiscovery)
	for _, item := range a.Items {
		fmt.Fprintf(&sb, "\tProcess: %s (%s)\n", item.Process_name, item.Process_code)
		fmt.Fprintf(&sb, "\tStatus: %s\n", item.Status)
		fmt.Fprintln(&sb, "\t---------")
	}
	return sb.String()
}

func (r BibRecord) toCSV() [][]string {
	var records [][]string
	for _, sudocLoc := range r.SudocLocations {
		records = append(records, []string{r.PPN, sudocLoc.ILN, "", sudocLoc.Name + " - " + sudocLoc.Sublocation, sudocLoc.RCR})
	}
	for _, almaLoc := range r.AlmaLocations {
		records = append(records, []string{r.PPN, "", almaLoc.Library_name, "", ""})
	}
	return records
}

func (a AlmaLocation) Valid() bool {
	valid := true
	for _, item := range a.Items {
		// TODO: use a slice instead of "ACQ" to be able to add status to be
		// ignored
		valid = valid && item.Process_code != "ACQ"
	}
	return valid && a.NoDiscovery
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

func Filter[T any](collection []T, test func(T) bool) []T {
	var res []T
	for _, e := range collection {
		if test(e) {
			res = append(res, e)
		}
	}
	return res
}
