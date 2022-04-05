package bib

import (
	"strings"
	"testing"
)

func TestDecodeAlmaXML(t *testing.T) {
	var tests = [][]byte{almaErrMMS, almaErrSys, nil}
	for _, test := range tests {
		_, err := decodeAlmaXML(test)
		if err == nil {
			t.Errorf(`decodeAlmaXML does not return error for %s`, test)
		}
	}
	r1, err := decodeAlmaXML(almaNoLoc)
	if err != nil {
		t.Errorf(`decodeAlmaXML returns error for NoLoc XML : %s`, err)
	}
	if len(r1) > 0 {
		t.Error(`decodeAlmaXML with 0 bibs returns a non zero slice`)
	}

	r2, err := decodeAlmaXML(almaOk)
	if err != nil {
		t.Errorf(`decodeAlmaXML returned error (%s) for almaOK XML`, err)
	}
	if len(r2) != 2 {
		t.Error(`len(almaOK) != 2`)
	}
	var testsOk = []struct {
		input string
		want  string
	}{
		{r2[0].ownerCode, "BIB1"},
		{r2[0].collection, "LOC1"},
		{r2[1].ownerCode, "BIB2"},
		{r2[1].collection, "LOC2"},
	}
	for _, test := range testsOk {
		if test.input != test.want {
			t.Errorf("decodeAlmaXML: got %s, want %s", test.input, test.want)
		}
	}
}

func TestDecodeLocations(t *testing.T) {
	rcrs := []string{"RCR1", "RCR2", "RCR3"}

	_, err := decodeLocations(nil, rcrs)
	if err == nil {
		t.Error("decodeLocations does not return error for nil input")
	}

	empty, err := decodeLocations(sudocUnknown, rcrs)
	if err != nil {
		t.Errorf(`decodeLocations(sudocUnknown) returns error (%s)`, err)
	}
	if len(empty) > 0 {
		t.Errorf(`decodeLocations(sudocUnknown) should return an empty slice, got len(empty) == %d`, len(empty))
	}

	r, err := decodeLocations(sudocOk, rcrs)
	if err != nil {
		t.Errorf(`decodeLocations(sudocOk) returns error (%s)`, err)
	}

	var testsStr = []struct {
		input string
		want  string
	}{
		{r[0].ppn, "PPN1"},
		{r[1].ppn, "PPN2"},
		{r[0].sudocLocations[0].rcr, "RCR1"},
		{r[1].sudocLocations[0].rcr, "RCR2"},
		{r[1].sudocLocations[1].rcr, "RCR3"},
	}
	for _, test := range testsStr {
		if test.input != test.want {
			t.Errorf("decodeLocations: want %s, got %s", test.input, test.want)
		}
	}

	var testsLen = []struct {
		input int
		want  int
	}{
		{len(r), 2},
		{len(r[0].sudocLocations), 1},
		{len(r[1].sudocLocations), 2},
	}
	for _, test := range testsLen {
		if test.input != test.want {
			t.Errorf("decodeLocations: want %d, got %d", test.input, test.want)
		}
	}

	if !strings.HasPrefix(r[0].sudocLocations[0].name, "BIB1") {
		t.Errorf(`decodeLocations: r[0] first name is wrong : want "BIB1",
			got %q`, r[0].sudocLocations[0].name)
	}
	if !strings.HasPrefix(r[1].sudocLocations[0].name, "BIB2") {
		t.Errorf(`decodeLocations: r[1] first name is wrong : want "BIB2",
			got %q`, r[1].sudocLocations[0].name)
	}
	if !strings.HasPrefix(r[1].sudocLocations[1].name, "BIB4") {
		t.Errorf(`decodeLocations: r[1] second name is wrong : want "BIB4",
			got %q`, r[1].sudocLocations[1].name)
	}
}

func TestDecodeRCR(t *testing.T) {
	var tests = []struct {
		input []byte
		want  []string
	}{
		{rcrUnknown, []string{}},
		{rcrOk, []string{"000000001", "000000002", "000000003"}},
	}
	for _, test := range tests {
		got, err := decodeRCR(test.input)
		if err != nil {
			t.Errorf("decodeRCR(%q) returned error: %s", test.input, err)
		}
		if !equalStrings(got, test.want) {
			t.Errorf("decodeRCR: want %v, got %v", test.want, got)
		}
	}
	_, err := decodeRCR(nil)
	if err == nil {
		t.Error("decodeRCR(nil) does not return error")
	}
}
