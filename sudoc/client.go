package sudoc

import (
	"casl/entities"
	"casl/marc"
	"casl/requests"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// SudocClient represents the main object to interact with.
type SudocClient struct {
	rcrs    map[string]library
	stats   stats
	fetcher requests.Fetcher
}

// Internal representation of a library in SUDOC's sense.
type library struct {
	iln  string
	rcr  string
	name string
}

type stats struct {
	iln2rcr int
	marcxml int
}

const (
	DEFAULT_BASE_URL = "https://www.sudoc.fr/"
	ILN2RCR_URL      = "https://www.idref.fr/services/iln2rcr/"
)

// NewSudocClient provides a SUDOC client including RCR->library mappings built
// from the iln2rcr API.
func NewSudocClient(ilns []string, fetcher requests.Fetcher) (*SudocClient, error) {
	if ilns == nil || len(ilns) == 0 {
		return nil, errors.New("NewSudocClient: empty or nil list of ILNs")
	}
	if fetcher == nil {
		return nil, errors.New("NewSudocClient: no HTTP client provided")
	}

	var client SudocClient
	client.fetcher = fetcher

	rcrs, err := client.getRCRs(ilns)
	if err != nil {
		return nil, fmt.Errorf("NewSudocClient: mapping failed: %w", err)
	}
	client.rcrs = rcrs
	return &client, nil
}

// GetFilteredLocations gets SUDOC locations of a given PPN, properly filled,
// from the unimarc2marcxml API. Only the locations regarding the RCRs of
// interest, given as a second argument, are provided.
func (sc *SudocClient) GetFilteredLocations(ppn string, rcrs []string) ([]*entities.SudocLocation, error) {
	var filtered []*entities.SudocLocation
	locations, err := sc.GetLocations(ppn)
	if err != nil {
		return filtered, err
	}

	for _, location := range locations {
		if slices.Contains(rcrs, location.RCR) {
			filtered = append(filtered, location)
		}
	}
	return filtered, nil
}

// GetLocations gets all the SUDOC locations of a given PPN, from the
// unimarc2marcxml API, filled with data from client's RCR mappings.
func (sc *SudocClient) GetLocations(ppn string) ([]*entities.SudocLocation, error) {
	var locs []*entities.SudocLocation
	sc.stats.marcxml += 1
	data, err := sc.fetcher.Fetch(DEFAULT_BASE_URL + ppn + ".xml")
	if err != nil {
		return locs, fmt.Errorf("ppn %s: %w\n", ppn, err)
	}
	marcRecord, err := marc.NewRecord(data)
	if err != nil {
		return locs, err
	}

	for _, field := range marcRecord.GetField("930") {
		rcr := field.GetValue("5")
		if len(rcr) != 1 {
			return locs, errors.New("MARC 930$5 does not contain a unique location")
		}

		sublocation := field.GetValue("c")
		if len(sublocation) > 1 {
			return locs, errors.New("MARC 930$c does not contain a unique value")
		}

		var location entities.SudocLocation
		location.RCR = strings.Split(rcr[0], ":")[0]
		if len(sublocation) == 1 {
			location.Sublocation = sublocation[0]
		}

		// Add informations from the RCR mappings
		location.ILN = sc.rcrs[location.RCR].iln
		location.Name = sc.rcrs[location.RCR].name
		locs = append(locs, &location)
	}
	return locs, nil
}

// Stats returns numbers of requests made by the client to the service named
// by the argument ("iln2rcr", "marcxml", "total").
// TODO: provide a better way to select the stat than by string
func (sc *SudocClient) Stats(t string) int {
	switch t {
	case "iln2rcr":
		return sc.stats.iln2rcr
	case "marcxml":
		return sc.stats.marcxml
	case "total":
		return sc.stats.iln2rcr + sc.stats.marcxml
	default:
		return sc.Stats("total")
	}
}

// getRCRs builds the map RCR->Library from the iln2rcr service.
func (sc *SudocClient) getRCRs(rcrs []string) (map[string]library, error) {
	url := ILN2RCR_URL + strings.Join(rcrs, ",")
	sc.stats.iln2rcr += 1
	data, err := sc.fetcher.Fetch(url)
	if err != nil {
		return nil, fmt.Errorf("getRCRs: iln2rcr failed: %w", err)
	}
	result, err := decodeRCR(data)
	if err != nil {
		return nil, fmt.Errorf("getRCRs: decoding XML failed: %w", err)
	}
	return result, nil
}
