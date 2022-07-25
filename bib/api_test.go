package bib

import (
	"casl/marc"
	"errors"
	"strconv"
	"testing"
	"time"

	"casl/alma"
)

// Mocking the HttpFetcher
type mockHttpFetch struct{}

func (f *mockHttpFetch) FetchAll(ppns []string) [][]byte {
	if len(ppns) == 0 {
		return [][]byte{}
	}
	invalid := []string{"invalid1", "invalid2", "invalid3"}
	if in(ppns[0], invalid) && in(ppns[1], invalid) && in(ppns[2], invalid) {
		return [][]byte{[]byte(`<?xml version="1.0" encoding="UTF-8"?>
			<sudoc service="multiwhere">
			<error>Found a null xml in result : values={ppn=PPNTEST}, query=select autorites.MULTIWHERE(#ppn#) from dual </error>
			</sudoc>`)}
	}
	semi := []string{"invalid1", "invalid2"}
	// ppns[0] invalid
	if in(ppns[0], semi) {
		// one of the two ppns must be invalid and the other valid
		if (in(ppns[1], semi) && ppns[2] == "ppn000001") ||
			(ppns[1] == "ppn000001" && in(ppns[2], semi)) {
			return mwsemi
		}
		// ppns[0] valid
	} else if ppns[0] == "ppn000001" {
		// the two remaining ppns must be invalid
		if in(ppns[1], semi) && in(ppns[2], semi) {
			return mwsemi
		}
	}
	return mwxmls
}

func (f *mockHttpFetch) FetchPPN(ppn string, secretParam string) []byte {
	if ppn == "ppn000001" {
		return fakeAlmaRes
	}
	return almaNoLoc
}

func (f *mockHttpFetch) FetchRCR(ilns []string) []byte {
	input := []string{"INVALID", "ILN02"}
	invalid := []string{"INVALID"}
	if len(ilns) == 0 || equalStrings(ilns, invalid) {
		return fakeIln2rcrError
	}
	if equalStrings(ilns, input) {
		return fakeIln2rcr_2
	}
	return fakeIln2rcr
}

func (f *mockHttpFetch) FetchMarc(ppn string) []byte {
	if ppn == "ppn_ok1" || ppn == "ppn_ok2" {
		return marcAax
	} else if ppn == "ppn_err" {
		return marcError
	}
	return marcOax
}

// Tests
func TestGetSudocLocations(t *testing.T) {
	// testdata ppns
	ppns_empty := make(map[string]bool)
	ppns_invalid := make(map[string]bool)
	ppns_invalid["invalid1"] = true
	ppns_invalid["invalid2"] = true
	ppns_invalid["invalid3"] = true
	ppns_semiInvalid := make(map[string]bool)
	ppns_semiInvalid["invalid1"] = true
	ppns_semiInvalid["ppn000001"] = true
	ppns_semiInvalid["invalid2"] = true
	ppns_valid := make(map[string]bool)
	ppns_valid["ppn000001"] = true
	ppns_valid["ppn000002"] = true
	ppns_valid["ppn000003"] = true

	// testdata rcrs
	rcrs_empty := []string{}
	rcrs_invalid := []string{"invalid1", "invalid2", "invalid3"}
	rcrs_semiInvalid := []string{"invalid1", "invalid2", "rcr000001", "rcr000003"}
	rcrs_valid := []string{"rcr000001", "rcr000002", "rcr000003"}

	// expected BibRecord
	expected := makeBibRecords()

	var tests = []struct {
		ppns_input map[string]bool
		rcrs_input []string
		want       []BibRecord
	}{
		{nil, rcrs_valid, []BibRecord{}},
		{ppns_valid, nil, []BibRecord{expected["br_empty1"], expected["br_empty2"], expected["br_empty3"]}},
		{ppns_empty, rcrs_valid, []BibRecord{}},
		{ppns_valid, rcrs_empty, []BibRecord{expected["br_empty1"], expected["br_empty2"], expected["br_empty3"]}},
		{ppns_valid, rcrs_semiInvalid, []BibRecord{expected["br1"], expected["br3"], expected["br4"]}},
		{ppns_valid, rcrs_invalid, []BibRecord{expected["br_empty1"], expected["br_empty2"], expected["br_empty3"]}},
		{ppns_valid, rcrs_valid, []BibRecord{expected["br2"], expected["br3"], expected["br4"]}},
		{ppns_invalid, rcrs_valid, []BibRecord{}},
		{ppns_semiInvalid, rcrs_valid, []BibRecord{expected["br2"]}},
	}

	fetcher := mockHttpFetch{}
	for i, test := range tests {
		got := GetSudocLocations(test.ppns_input, test.rcrs_input, &fetcher)
		if !equalBibRecords(got, test.want) {
			t.Errorf("[%d] GetSudocLocations with %v and %v returned %v ; want %v",
				i, test.ppns_input, test.rcrs_input, got, test.want)
		}
	}
}

// Mocking the AlmaClient
type mockAlmaClient struct{}

func (m mockAlmaClient) GetHoldingsFromPPN(ppn alma.PPN) ([]alma.Holding, error) {
	result := []alma.Holding{}
	switch ppn {
	case "fetcherror":
		return nil, errors.New("Fetch error")
	case "oneloc":
		result = append(result, alma.Holding{Library: "BIB1", Location: "LOC1"})
	case "twoloc":
		result = append(result, alma.Holding{Library: "BIB1", Location: "LOC1"})
		result = append(result, alma.Holding{Library: "BIB2", Location: "LOC2"})
	}
	time.Sleep(10 * time.Millisecond)
	return result, nil
}

func TestGetAlmaLocations(t *testing.T) {
	client := mockAlmaClient{}
	input := []BibRecord{
		{ppn: "noalmaloc"},
		{ppn: "fetcherror"},
		{ppn: "oneloc"},
		{ppn: "twoloc"},
	}
	expected := []BibRecord{
		{ppn: "noalmaloc", almaLocations: []almaLocation{}},
		{ppn: "oneloc", almaLocations: []almaLocation{
			{collection: "LOC1", ownerCode: "BIB1", rcr: "RCR1"},
		}},
		{ppn: "twoloc", almaLocations: []almaLocation{
			{collection: "LOC1", ownerCode: "BIB1", rcr: "RCR1"},
			{collection: "LOC2", ownerCode: "BIB2", rcr: "RCR2"},
		}},
	}

	rcrs := make(map[string]string)
	rcrs["BIB1"] = "RCR1"
	rcrs["BIB2"] = "RCR2"

	got := GetAlmaLocations(client, input, rcrs)
	if !equalBibRecords(got, expected) {
		t.Errorf("GetAlmaLocations with %v returned %v ; want %v", input, got, expected)
	}
}

func TestGetRCRs(t *testing.T) {
	var tests = []struct {
		input []string
		want  []string
	}{
		{[]string{"ILN01", "ILN02"}, []string{"rcr000001", "rcr000002", "rcr000003", "rcr000004"}},
		{[]string{"INVALID", "ILN02"}, []string{"rcr000004"}},
		{[]string{"INVALID"}, []string{}},
		{[]string{}, []string{}},
		{nil, []string{}},
	}

	fetcher := mockHttpFetch{}
	for i, test := range tests {
		got, _ := GetRCRs(test.input, &fetcher)
		if !equalStrings(got, test.want) {
			t.Errorf("[%d] GetRCRs with %v returned %v : want %v", i, test.input, got, test.want)
		}
	}
}

func TestComparePPN(t *testing.T) {
	var cr_alma []CRecord
	cr_alma = append(cr_alma, CRecord{PPN: "ppn_alma", AlmaLibrary: "code01", SUDOCLibrary: "",
		ILN: "", RCR: "rcr01", InAlma: true, InSUDOC: false})
	cr_alma = append(cr_alma, CRecord{PPN: "ppn_alma", AlmaLibrary: "code02", SUDOCLibrary: "",
		ILN: "", RCR: "rcr02", InAlma: true, InSUDOC: false})
	cr_alma = append(cr_alma, CRecord{PPN: "ppn_alma", AlmaLibrary: "code04", SUDOCLibrary: "",
		ILN: "", RCR: "rcr04", InAlma: true, InSUDOC: false})
	var cr_sudoc []CRecord
	cr_sudoc = append(cr_sudoc, CRecord{PPN: "ppn_sudoc", AlmaLibrary: "", SUDOCLibrary: "name01",
		ILN: "", RCR: "rcr01", InAlma: false, InSUDOC: true})
	cr_sudoc = append(cr_sudoc, CRecord{PPN: "ppn_sudoc", AlmaLibrary: "", SUDOCLibrary: "name02",
		ILN: "", RCR: "rcr02", InAlma: false, InSUDOC: true})
	var cr_ignored []CRecord
	cr_ignored = append(cr_ignored, CRecord{PPN: "ppn_ignored", AlmaLibrary: "", SUDOCLibrary: "name01",
		ILN: "", RCR: "rcr01", InAlma: false, InSUDOC: true})
	cr_ignored = append(cr_ignored, CRecord{PPN: "ppn_ignored", AlmaLibrary: "", SUDOCLibrary: "name02",
		ILN: "", RCR: "rcr02", InAlma: false, InSUDOC: true})
	cr_ignored = append(cr_ignored, CRecord{PPN: "ppn_ignored", AlmaLibrary: "", SUDOCLibrary: "name03",
		ILN: "", RCR: "rcr03", InAlma: false, InSUDOC: true})

	var sl []sudocLocation
	sl = append(sl, sudocLocation{iln: "iln01", rcr: "rcr01", name: "name01"})
	sl = append(sl, sudocLocation{iln: "iln02", rcr: "rcr02", name: "name02"})
	sl = append(sl, sudocLocation{iln: "iln03", rcr: "rcr03", name: "name03"})
	sl = append(sl, sudocLocation{iln: "iln04", rcr: "rcr04", name: "name04"})
	var al []almaLocation
	al = append(al, almaLocation{collection: "coll01", rcr: "rcr01", ownerCode: "code01"})
	al = append(al, almaLocation{collection: "coll02", rcr: "rcr02", ownerCode: "code02"})
	al = append(al, almaLocation{collection: "ignored", rcr: "rcr03", ownerCode: "code03"})
	al = append(al, almaLocation{collection: "coll04", rcr: "rcr04", ownerCode: "code04"})
	al = append(al, almaLocation{collection: "ignored2", rcr: "rcr05", ownerCode: "code05"})
	var match []BibRecord
	match = append(match, BibRecord{ppn: "ppn01", sudocLocations: []sudocLocation{sl[0], sl[1]},
		almaLocations: []almaLocation{al[0], al[1], al[2]}})
	match = append(match, BibRecord{ppn: "ppn02", sudocLocations: []sudocLocation{sl[3]},
		almaLocations: []almaLocation{al[2], al[3]}})
	var nomatch []BibRecord
	nomatch = append(nomatch, BibRecord{ppn: "ppn_alma", sudocLocations: []sudocLocation{},
		almaLocations: []almaLocation{al[0], al[1], al[2], al[3]}})
	var onlysu []BibRecord
	onlysu = append(onlysu, BibRecord{ppn: "ppn_sudoc", sudocLocations: []sudocLocation{sl[0], sl[1]},
		almaLocations: []almaLocation{al[2]}})
	var allignored []BibRecord
	allignored = append(allignored, BibRecord{ppn: "ppn_ignored", sudocLocations: []sudocLocation{sl[0], sl[1], sl[2]},
		almaLocations: []almaLocation{al[2], al[4]}})

	var tests = []struct {
		input []BibRecord
		want  []CRecord
	}{
		{[]BibRecord{}, []CRecord{}},
		{match, []CRecord{}},
		{nomatch, cr_alma},
		{onlysu, cr_sudoc},
		{allignored, cr_ignored},
	}
	for i, test := range tests {
		got := ComparePPN(test.input, []string{"ignored", "ignored2"})
		if !equalCRecords(test.want, got) {
			t.Errorf("[%d] ComparePPN with %v returned %v : want %v", i, test.input,
				got, test.want)
		}
	}
}

func TestFilter(t *testing.T) {
	var crecords []CRecord
	crecords = append(crecords, CRecord{PPN: "ppn_ok1"})
	crecords = append(crecords, CRecord{PPN: "ppn_ok2"})
	crecords = append(crecords, CRecord{PPN: "ppn_err"})
	crecords = append(crecords, CRecord{PPN: "ppn_elec"})

	var tests = []struct {
		input []CRecord
		want  []CRecord
	}{
		{[]CRecord{crecords[0], crecords[1]}, []CRecord{crecords[0], crecords[1]}},
		{[]CRecord{crecords[0], crecords[1], crecords[3]}, []CRecord{crecords[0], crecords[1]}},
		{[]CRecord{crecords[0], crecords[2]}, []CRecord{crecords[0], crecords[2]}},
		{[]CRecord{crecords[3], crecords[3]}, []CRecord{}},
		{[]CRecord{}, []CRecord{}},
		{[]CRecord{crecords[2]}, []CRecord{crecords[2]}},
	}

	fetcher := mockHttpFetch{}
	for i, test := range tests {
		got := Filter(test.input, []string{}, &fetcher)
		if !equalCRecords(got, test.want) {
			t.Errorf("[%d] Filter with %v returned %v : want %v", i, test.input, got, test.want)
		}
	}
}

func TestExtractRCR(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"123456789:000111222", "123456789"},
		{"a:", "a"},
		{"", ""},
		{"rcr", "rcr"},
	}

	for _, test := range tests {
		if got := extractRCR(test.input); got != test.want {
			t.Errorf("%s returned %s but want %s", test.input, got, test.want)
		}
	}
}

func TestAddSublocation(t *testing.T) {
	r1, _ := marc.NewRecord(marcOK1)
	r2, _ := marc.NewRecord(marcOK2)
	r3, _ := marc.NewRecord(marcOK3)
	c1 := CRecord{RCR: "rcr-ok", SUDOCSublocation: ""}
	c2 := CRecord{RCR: "rcr-ok", SUDOCSublocation: ""}
	c3 := CRecord{RCR: "rcr-ok", SUDOCSublocation: ""}

	var tests = []struct {
		input1 *CRecord
		input2 *marc.Record
		want   string
	}{
		{&c1, r1, ""},
		{&c2, r2, "sublocation"},
		{&c3, r3, "sublocation1, sublocation2"},
	}

	for _, test := range tests {
		addSublocation(test.input1, test.input2)
		if test.input1.SUDOCSublocation != test.want {
			t.Errorf("got %q want %q", test.input1.SUDOCSublocation, test.want)
		}
	}
}

// benchmarks
type mockHttpClient struct{}

func (f *mockHttpClient) FetchAll(ppns []string) [][]byte {
	return [][]byte{}
}
func (f *mockHttpClient) FetchPPN(ppn string, secretParam string) []byte {
	return []byte{}
}

func (f *mockHttpClient) FetchRCR(ilns []string) []byte {
	return []byte{}
}

func (f *mockHttpClient) FetchMarc(ppn string) []byte {
	time.Sleep(10 * time.Millisecond)
	return []byte{}
}

func BenchmarkFilter(b *testing.B) {
	client := mockHttpClient{}
	var crecords []CRecord
	for i := 0; i < 250; i++ {
		id := "ppn" + strconv.Itoa(i)
		crecords = append(crecords, CRecord{PPN: id})
	}
	for i := 0; i < b.N; i++ {
		Filter(crecords, []string{}, &client)
	}
}

func BenchmarkGetAlmaLocations(b *testing.B) {
	client := mockAlmaClient{}
	rcrs := make(map[string]string)
	rcrs["BIB1"] = "RCR1"
	rcrs["BIB2"] = "RCR2"
	var brecords []BibRecord
	for i := 0; i < 250; i++ {
		brecords = append(brecords, BibRecord{ppn: "twoloc"})
	}
	for i := 0; i < b.N; i++ {
		GetAlmaLocations(client, brecords, rcrs)
	}
}

// helpers
func makeBibRecords() map[string]BibRecord {
	result := make(map[string]BibRecord)

	// BibRecord examples
	sl1 := sudocLocation{rcr: "rcr000001", name: "TEST1"}
	sl2 := sudocLocation{rcr: "rcr000002", name: "TEST2"}
	sl3 := sudocLocation{rcr: "rcr000003", name: "TEST3"}

	result["br_empty1"] = BibRecord{ppn: "ppn000001"}
	result["br_empty2"] = BibRecord{ppn: "ppn000002"}
	result["br_empty3"] = BibRecord{ppn: "ppn000003"}
	result["br1"] = BibRecord{ppn: "ppn000001", sudocLocations: []sudocLocation{sl1, sl3}}
	result["br2"] = BibRecord{ppn: "ppn000001", sudocLocations: []sudocLocation{sl1, sl2, sl3}}
	result["br3"] = BibRecord{ppn: "ppn000002", sudocLocations: []sudocLocation{sl1}}
	result["br4"] = BibRecord{ppn: "ppn000003", sudocLocations: []sudocLocation{sl1}}

	return result
}
