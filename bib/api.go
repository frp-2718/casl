// bib handles operations on bibliographic objects.
package bib

import (
	"casl/alma"
	"casl/marc"
	"casl/requests"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

const MAX_CONCURRENT_REQUESTS = 50

type AlmaClient interface {
	GetHoldingsFromPPN(ppn alma.PPN) ([]alma.Holding, error)
}

// GetSudocLocations fetches locations and returns populated BibRecords
// corresponding to the given ppns.
func GetSudocLocations(ppns map[string]bool, rcrs []string, client requests.Fetcher) []BibRecord {
	var records []BibRecord

	ppnsList := mapKeys(ppns)

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
func GetAlmaLocations(a AlmaClient, bibs []BibRecord, rcrMap map[string]string) []BibRecord {
	var tokens = make(chan struct{}, MAX_CONCURRENT_REQUESTS)
	wg := sync.WaitGroup{}
	var mu = &sync.Mutex{}
	var result []BibRecord

	for _, record := range bibs {
		wg.Add(1)
		go func(record BibRecord) {
			tokens <- struct{}{}
			defer wg.Done()
			locations, err := a.GetHoldingsFromPPN(alma.PPN(record.ppn))
			if err != nil {
				// TODO: handle fetch errors
				<-tokens
				return // ignore errors
			}
			record.almaLocations = convertLocations(locations)

			var newLocations []almaLocation
			// If there are no locations in Alma, len(record.almaLocations) == 0.
			for _, l := range record.almaLocations {
				nl := almaLocation{ownerCode: l.ownerCode, collection: l.collection,
					rcr: rcrMap[l.ownerCode]}
				newLocations = append(newLocations, nl)
			}
			if len(newLocations) > 0 {
				record.almaLocations = newLocations
			}
			mu.Lock()
			result = append(result, record)
			mu.Unlock()
			<-tokens
		}(record)
	}
	wg.Wait()
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

func Filter(records []CRecord, monoRCRs []string, client requests.Fetcher) []CRecord {
	var tokens = make(chan struct{}, MAX_CONCURRENT_REQUESTS)
	var res []CRecord
	wg := sync.WaitGroup{}
	var mu = &sync.Mutex{}

	for _, record := range records {
		wg.Add(1)
		go func(r CRecord) {
			tokens <- struct{}{}
			defer wg.Done()
			marcxml := client.FetchMarc(r.PPN)
			marcrecord, err := marc.NewRecord(marcxml)
			if err != nil {
				//log.Printf("bib.Filter: unable to create MARC record from PPN %s", r.PPN)
				// this is always safe to ignore errors and pretend that the record
				// should not be filtered
				mu.Lock()
				res = append(res, r)
				mu.Unlock()
				<-tokens
				return
			}
			class := marcrecord.GetField("008")[0].GetValue("")[0]
			// Exclusion of electronic resources
			if !strings.HasPrefix(class, "O") {
				// Add sublocations for some monolithic RCRs
				if r.SUDOCLibrary != "" && slices.Contains(monoRCRs, r.RCR) {
					addSublocation(&r, marcrecord)
				}
				mu.Lock()
				res = append(res, r)
				mu.Unlock()
			}
			<-tokens
		}(record)
	}
	wg.Wait()
	return res
}

func addSublocation(r *CRecord, m *marc.Record) {
	fields := m.GetField("930")
	sep := ""
	var sublocations []string
	for _, f := range fields {
		if extractRCR(f.GetValue("5")[0]) == r.RCR {
			if r.SUDOCSublocation != "" {
				sep = ", "
			}
			if sublocation := f.GetValue("c"); sublocation != nil {
				if !slices.Contains(sublocations, sublocation[0]) {
					sublocations = append(sublocations, sublocation[0])
					r.SUDOCSublocation = r.SUDOCSublocation + sep + sublocation[0]
				}
			}
		}
	}
}

func extractRCR(from string) string {
	return strings.Split(from, ":")[0]
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
		if !slices.Contains(ignored, loc.collection) {
			result = append(result, loc)
		}
	}
	return result
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
