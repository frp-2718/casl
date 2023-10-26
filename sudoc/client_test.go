package sudoc

import (
	"os"
	"reflect"
	"testing"
)

type mockHttpFetcher struct{}

func (m mockHttpFetcher) Fetch(url string) ([]byte, error) {
	data, err := os.ReadFile("testdata/iln2rcr.xml")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func TestGetRCRs(t *testing.T) {
	input := []string{"1", "2"}
	sc, _ := NewSudocClient(input)
	sc.fetcher = mockHttpFetcher{}

	want := make(map[string]library)
	want["100000001"] = library{"1", "100000001", "UNIV-1.1"}
	want["200000001"] = library{"2", "200000001", "UNIV-2.1"}
	want["200000002"] = library{"2", "200000002", "UNIV-2.2"}
	want["200000003"] = library{"2", "200000003", "UNIV-2.3"}

	got, err := sc.getRCRs(input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}
