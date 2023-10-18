package bib

import (
	"fmt"
	"strings"
)

type BibRecord struct {
	PPN            string
	mms            string
	SudocLocations []*sudocLocation
	almaLocations  []almaLocation
}

type sudocLocation struct {
	iln         string
	rcr         string
	name        string
	sublocation string
}

type almaLocation struct {
	collection string
	ownerCode  string
	rcr        []string
}

// CRecord is the fusion of a SUDOC record and an Alma record.
type CRecord struct {
	PPN              string
	AlmaLibrary      string
	SUDOCLibrary     string
	SUDOCSublocation string
	ILN              string
	RCR              []string
	InAlma           bool
	InSUDOC          bool
}

func (r BibRecord) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*** START ***\n")
	fmt.Fprintf(&sb, "PPN: %s\n", r.PPN)
	for _, sl := range r.SudocLocations {
		fmt.Fprintf(&sb, "%s\n----------\n", sl)
	}
	fmt.Fprintf(&sb, "*** END ***\n")
	return sb.String()
}

func (s sudocLocation) String() string {
	return fmt.Sprintf("ILN: %s\nRCR: %s\nNAME: %s\nSUBLOCATION: %s\n",
		s.iln, s.rcr, s.name, s.sublocation)
}
