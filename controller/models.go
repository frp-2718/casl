package controller

import (
	"casl/entities"
	"fmt"
	"strings"
)

type suClient interface {
	GetLocations(ppn string) ([]*entities.SudocLocation, error)
	GetFilteredLocations(ppn string, rcrs []string) ([]*entities.SudocLocation, error)
}

type almaClient interface {
	GetLocations(ppn string) ([]*entities.AlmaLocation, error)
	GetFilteredLocations(ppn string, lib_codes []string) ([]*entities.AlmaLocation, error)
}

type Controller struct {
	Config     *config
	Mappings   *mappings
	SUClient   suClient
	AlmaClient almaClient
}

type config struct {
	MappingFilePath string   `json:"alma-rcr_file_path"`
	AlmaAPIKey      string   `json:"alma_api_key"`
	ILNs            []string `json:"iln_to_track"`
	IgnoredAlmaColl []string `json:"ignored_alma_collections"`
	IgnoredSudocRCR []string `json:"ignored_sudoc_rcr"`
	MonolithicRCR   []string `json:"monolithic_rcr"`
	FollowedRCR     []string
	FolowedLibs     []string
}

// Mappings Alma/RCR, Alma/Libraries names, RCR/ILN, RCR/label, read from CSV.
type mappings struct {
	alma2rcr map[string][]string
	rcr2iln  map[string]string
	alma2str map[string]string
	rcr2str  map[string]string
	rcr2alma map[string][]string
}

func (c Controller) String() string {
	var sb strings.Builder
	fmt.Fprintln(&sb, c.Config)
	fmt.Fprintln(&sb, c.Mappings)
	fmt.Fprintf(&sb, "*** CLIENTS\n\n")
	fmt.Fprintf(&sb, "SUDOC client: %s\n", c.SUClient)
	fmt.Fprintf(&sb, "Alma client: %p\n", c.AlmaClient)
	return sb.String()
}

func (c config) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*** CONFIG\n\n")
	fmt.Fprintf(&sb, "Mapping file path: %s\n", c.MappingFilePath)
	fmt.Fprintf(&sb, "Alma API key: %s\n", c.AlmaAPIKey)
	fmt.Fprintf(&sb, "Concerned ILN: %v\n", c.ILNs)
	fmt.Fprintf(&sb, "Alma collections to ignore: %v\n", c.IgnoredAlmaColl)
	fmt.Fprintf(&sb, "RCR to ignore: %v\n", c.IgnoredSudocRCR)
	fmt.Fprintf(&sb, "RCR with sublocations: %v\n", c.MonolithicRCR)
	fmt.Fprintf(&sb, "RCR to inspect: %v\n", c.FollowedRCR)
	return sb.String()
}

func (m mappings) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "*** MAPPINGS\n\n")
	fmt.Fprintln(&sb, "-- Alma -> RCR:")
	for k, v := range m.alma2rcr {
		fmt.Fprintf(&sb, "%s -> %v\n", k, v)
	}
	fmt.Fprintln(&sb, "\n-- RCR -> ILN:")
	for k, v := range m.rcr2iln {
		fmt.Fprintf(&sb, "%s -> %v\n", k, v)
	}
	fmt.Fprintln(&sb, "\n-- Alma code -> Alma name:")
	for k, v := range m.alma2str {
		fmt.Fprintf(&sb, "%s -> %v\n", k, v)
	}
	fmt.Fprintln(&sb, "\n-- RCR -> SUDOC name:")
	for k, v := range m.rcr2str {
		fmt.Fprintf(&sb, "%s -> %s\n", k, v)
	}
	return sb.String()
}
