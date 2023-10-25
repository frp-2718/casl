package sudoc

import (
	"casl/entities"
	"casl/marc"
	"casl/requests"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

type SudocClient struct {
	marcxml_url string
	rcrs        map[string]library
	stats       int
}

type library struct {
	iln  string
	rcr  string
	name string
}

func NewSudocClient() *SudocClient {
	var client SudocClient
	client.marcxml_url = "https://www.sudoc.fr/"
	client.makeRCRs("sudoc_rcr.csv")
	return &client
}

func (sc *SudocClient) GetFilteredLocations(ppn string, rcrs []string) ([]*entities.SudocLocation, error) {
	locations, err := sc.GetLocations(ppn)
	var filtered []*entities.SudocLocation
	if err != nil {
		return filtered, err
	}

	for _, location := range locations {
		if slices.Contains(rcrs, location.RCR) {
			location.ILN = sc.rcrs[location.RCR].iln
			location.Name = sc.rcrs[location.RCR].name
			filtered = append(filtered, location)
		}
	}
	return filtered, nil
}

// TODO: abstract the requester
func (sc *SudocClient) GetLocations(ppn string) ([]*entities.SudocLocation, error) {
	var locs []*entities.SudocLocation
	sc.stats += 1
	data, err := requests.FetchMarc(ppn)
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

		locs = append(locs, &location)
	}
	return locs, nil
}

func (sc *SudocClient) Stats() string {
	return fmt.Sprintf("unimarc2marcxml: %d\n", sc.stats)
}

func (sc *SudocClient) makeRCRs(csv_file string) {
	f, err := os.Open(csv_file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	sc.rcrs = make(map[string]library)

	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		lib := library{}
		lib.iln = record[1]
		lib.rcr = record[0]
		lib.name = record[2]
		sc.rcrs[lib.rcr] = lib
	}
}
