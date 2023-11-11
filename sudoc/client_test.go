package sudoc

import (
	"casl/entities"
	"casl/requests"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

type mockHttpFetcher struct{}

func (m mockHttpFetcher) Fetch(url string) ([]byte, error) {
	switch url {
	case ILN2RCR_URL + "1,2":
		data, err := os.ReadFile("testdata/iln2rcr.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case ILN2RCR_URL + "not_found":
		data, err := os.ReadFile("testdata/iln2rcr_not_found.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case DEFAULT_BASE_URL + "ppn" + ".xml":
		data, err := os.ReadFile("testdata/marcxml.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case DEFAULT_BASE_URL + "ppn_no_locations" + ".xml":
		data, err := os.ReadFile("testdata/marcxml_no_locations.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, nil
	}
}

func TestNewSudocClient(t *testing.T) {
	tests_err := []struct {
		name    string
		input   []string
		fetcher requests.Fetcher
	}{
		{"nil, nil", nil, nil},
		{"empty, nil", []string{}, nil},
		{"empty, ok", []string{}, mockHttpFetcher{}},
		{"empty, ok", []string{"1", "2"}, nil},
	}

	for _, test := range tests_err {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewSudocClient(test.input, test.fetcher)
			if err == nil {
				t.Error("should return an error")
			}
		})
	}

	rcrs := make(map[string]library)
	rcrs["100000001"] = library{"1", "100000001", "UNIV-1.1"}
	rcrs["200000001"] = library{"2", "200000001", "UNIV-2.1"}
	rcrs["200000002"] = library{"2", "200000002", "UNIV-2.2"}
	rcrs["200000003"] = library{"2", "200000003", "UNIV-2.3"}

	var stats stats
	stats.iln2rcr = 1
	want_client := &SudocClient{rcrs: rcrs, stats: stats, fetcher: mockHttpFetcher{}}
	got_client, err := NewSudocClient([]string{"1", "2"}, mockHttpFetcher{})
	t.Run("ok, nil", func(t *testing.T) {
		if err != nil {
			t.Error("got error, want properly built SudocClient")
		}
		if !reflect.DeepEqual(got_client, want_client) {
			t.Errorf("got %+v, want %+v", got_client, want_client)
		}
	})
}

func TestGetLocations(t *testing.T) {
	sc, err := NewSudocClient([]string{"1", "2"}, mockHttpFetcher{})
	if err != nil {
		t.Fatal("NewSudocClient failed")
	}

	var empty []*entities.SudocLocation
	locations := []*entities.SudocLocation{
		{ILN: "1", RCR: "100000001", Name: "UNIV-1.1", Sublocation: "SUB1"},
		{ILN: "2", RCR: "200000001", Name: "UNIV-2.1"},
		{ILN: "2", RCR: "200000001", Name: "UNIV-2.1"},
		{ILN: "2", RCR: "200000002", Name: "UNIV-2.2"},
	}
	tests := []struct {
		name  string
		input string
		want  []*entities.SudocLocation
	}{
		{"complete locations", "ppn", locations},
		{"no locations", "ppn_no_locations", empty},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := sc.GetLocations(test.input)
			if err != nil {
				t.Error("unexpected error")
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetFilteredLocations(t *testing.T) {
	sc, err := NewSudocClient([]string{"1", "2"}, mockHttpFetcher{})
	if err != nil {
		t.Fatal("NewSudocClient failed")
	}
	var empty []*entities.SudocLocation
	locations := []*entities.SudocLocation{
		{ILN: "1", RCR: "100000001", Name: "UNIV-1.1", Sublocation: "SUB1"},
		{ILN: "2", RCR: "200000001", Name: "UNIV-2.1"},
		{ILN: "2", RCR: "200000001", Name: "UNIV-2.1"},
		{ILN: "2", RCR: "200000002", Name: "UNIV-2.2"},
	}
	tests := []struct {
		name string
		ppn  string
		rcrs []string
		want []*entities.SudocLocation
	}{
		{"empty rcrs", "ppn", []string{}, empty},
		{"no locations", "ppn_no_locations", []string{"100000001"}, empty},
		{"all locations", "ppn", []string{"100000001", "200000001", "200000002"}, locations},
		{"rcr_200000002", "ppn", []string{"200000002"}, []*entities.SudocLocation{
			{ILN: "2", RCR: "200000002", Name: "UNIV-2.2"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := sc.GetFilteredLocations(test.ppn, test.rcrs)
			if err != nil {
				t.Error("unexpected error")
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetRCRs(t *testing.T) {
	input_ok := []string{"1", "2"}
	input_ko := []string{"not_found"}
	sc_ok, err := NewSudocClient(input_ok, mockHttpFetcher{})
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewSudocClient(input_ko, mockHttpFetcher{})
	if err == nil {
		t.Fatalf("null xml should raise an error")
	}

	want := make(map[string]library)
	want["100000001"] = library{"1", "100000001", "UNIV-1.1"}
	want["200000001"] = library{"2", "200000001", "UNIV-2.1"}
	want["200000002"] = library{"2", "200000002", "UNIV-2.2"}
	want["200000003"] = library{"2", "200000003", "UNIV-2.3"}

	got, err := sc_ok.getRCRs(input_ok)
	if err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestStats(t *testing.T) {
	sc, _ := NewSudocClient([]string{"1", "2"}, mockHttpFetcher{})
	iln2rcr := sc.Stats("iln2rcr")
	marcxml := sc.Stats("marcxml")
	total := sc.Stats("total")
	if iln2rcr != 1 || marcxml != 0 || total != 1 {
		t.Errorf("got stats.iln2rcr = %d, stats.marcxml = %d, stats.total = %d, ant iln2rcr = 1, marcxml = 0, total = 1",
			iln2rcr, marcxml, total)
	}

	n := rand.Intn(100) + 10
	for i := 0; i < n; i++ {
		sc.getRCRs([]string{"1", "2"})
		sc.GetLocations("")
	}
	iln2rcr = sc.Stats("iln2rcr")
	marcxml = sc.Stats("marcxml")
	total = sc.Stats("total")
	if iln2rcr != n+1 || marcxml != n || total != n*2+1 {
		t.Errorf("got stats.iln2rcr = %d, stats.marcxml = %d, stats.total = %d, want iln2rcr = %d, marcxml = %d, total = %d",
			iln2rcr, marcxml, total, n+1, n, 2*n+1)
	}
}
