package entities

import (
	"fmt"
	"slices"
	"strings"
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

func (r BibRecord) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*** PPN: %s\n\n", r.PPN)
	fmt.Fprintf(&sb, "*** MMS: %s\n\n", r.MMS)
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

func (a *AlmaLocation) IsValid(ignored_locations []string) bool {
	if a.NoDiscovery || a.Items == nil || len(a.Items) == 0 || slices.Contains(ignored_locations, a.Location_code) {
		return false
	}
	for _, item := range a.Items {
		// TODO: use a configuration instead of "ACQ" to be able to add status to be
		// ignored
		if item.Process_code != "ACQ" {
			return true
		}
	}
	return false
}
