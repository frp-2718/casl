// bib handles operations on bibliographic objects.
package bib

import (
	"casl/alma"
	"casl/marc"
	"casl/requests"
	"log"
	"net/http"
	"strings"
)

// GetSudocLocations fetches locations and returns populated BibRecords
// corresponding to the given ppns.
func GetSudocLocations(ppns map[string]bool, rcrs []string, client requests.Fetcher) []BibRecord {
	var records []BibRecord

	// We need a slice from the set of ppns.
	ppnsList := make([]string, 0, len(ppns))

	for key := range ppns {
		ppnsList = append(ppnsList, key)
	}

	for _, data := range client.FetchAll(ppnsList) {
		sudocLoc, err := decodeLocations(data, rcrs)
		if err != nil {
			continue // ignore wrong ppns
		}
		records = append(records, sudocLoc...)
	}
	return records
}

// GetAlmaLocations fetches locations and returns populated BibRecords
// corresponding to the given SUDOC records.
func GetAlmaLocations(bibs []BibRecord, secretParam string, rcrMap map[string]string) []BibRecord {
	a, err := alma.New(&http.Client{}, secretParam, "")
	if err != nil {
		log.Fatalf("GetAlmaLocations: unable to initialize the Alma client: %v", err)
	}
	var result []BibRecord
	for _, record := range bibs {
		locations, err := a.GetHoldingsFromPPN(alma.PPN(record.ppn))
		if err != nil {
			// TODO: handle fetch errors
			continue // ignore errors
		}
		record.almaLocations = convertLocations(locations)

		var newLocations []almaLocation
		// If there are no locations in Alma, len(record.almaLocations) == 0.
		for _, l := range record.almaLocations {
			nl := almaLocation{ownerCode: l.ownerCode, collection: l.collection,
				rcr: rcrMap[l.ownerCode]}
			newLocations = append(newLocations, nl)
		}
		record.almaLocations = newLocations
		result = append(result, record)
	}
	return result
}

// GetRCRs returns a slice of RCRs as strings from a slice of ILNs.
func GetRCRs(ilns []string, client requests.Fetcher) ([]string, error) {
	xmldata := client.FetchRCR(ilns)
	result, err := decodeRCR(xmldata)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ComparePPN builds a list of results from SUDOC and Alma records which don't
// match.
func ComparePPN(records []BibRecord, ignoredCollections []string) []CRecord {
	var results []CRecord
	var result CRecord
	clean := removeIgnored(records, ignoredCollections)
	for _, record := range clean {
		for _, al := range record.almaLocations {
			if !almaInSudoc(al, record.sudocLocations) {
				result = CRecord{PPN: record.ppn, AlmaLibrary: al.ownerCode,
					RCR: al.rcr, InSUDOC: false, InAlma: true}
				results = append(results, result)
			}
		}
		for _, sl := range record.sudocLocations {
			if !sudocInAlma(sl, record.almaLocations) {
				result = CRecord{PPN: record.ppn, RCR: sl.rcr, SUDOCLibrary: sl.name,
					InSUDOC: true, InAlma: false}
				results = append(results, result)
			}
		}
	}
	return results
}

// Filter builds a list of CRecord by filtering the given list according to
// some very specific criterions.
func Filter(records []CRecord, client requests.Fetcher) []CRecord {
	var res []CRecord
	for _, record := range records {
		marcxml := client.FetchMarc(record.PPN)
		marcrecord, err := marc.NewRecord(marcxml)
		if err != nil {
			// ignore error for now
			res = append(res, record)
			continue
		}
		class := marcrecord.GetField("008")[0].GetValue("")[0]
		// Exclusion of electronic resources
		if !strings.HasPrefix(class, "O") {
			res = append(res, record)
		}
	}
	return res
}

func convertLocations(s1 []alma.Holding) []almaLocation {
	var result []almaLocation
	for _, holding := range s1 {
		result = append(result, almaLocation{collection: holding.Location, ownerCode: holding.Library})
	}
	return result
}

func almaInSudoc(al almaLocation, sl []sudocLocation) bool {
	for _, l := range sl {
		if l.rcr == al.rcr {
			return true
		}
	}
	return false
}

func sudocInAlma(sl sudocLocation, al []almaLocation) bool {
	for _, l := range al {
		if l.rcr == sl.rcr {
			return true
		}
	}
	return false
}

func removeIgnored(records []BibRecord, ignored []string) []BibRecord {
	var result []BibRecord
	for _, r := range records {
		r.almaLocations = removeElem(r.almaLocations, ignored)
		result = append(result, r)
	}
	return result
}

func removeElem(locations []almaLocation, ignored []string) []almaLocation {
	var result []almaLocation
	for _, loc := range locations {
		if !in(loc.collection, ignored) {
			result = append(result, loc)
		}
	}
	return result
}
