package entities

import (
	"fmt"
	"strings"
)

type BibRecord struct {
	PPN            string
	MMS            string
	SudocLocations []*SudocLocation
	AlmaLocations  []*almaLocation
}

type SudocLocation struct {
	ILN         string
	RCR         string
	Name        string
	Sublocation string
}

type almaLocation struct {
	collection string
	ownerCode  string
	rcr        []string
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
