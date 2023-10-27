package sudoc

import (
	"casl/requests"
	"fmt"
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
	case DEFAULT_BASE_URL + "ppn" + ".xml":
		data, err := os.ReadFile("testdata/marcxml.xml")
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

	want_client := &SudocClient{rcrs, 0, mockHttpFetcher{}}
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

// TODO: complete this test
func TestGetFilteredLocations(t *testing.T) {
	sc, err := NewSudocClient([]string{"1", "2"}, mockHttpFetcher{})
	if err != nil {
		t.Fatal("NewSudocClient failed")
	}
	loc, err := sc.GetFilteredLocations("ppn", []string{})
	fmt.Println(loc)
	if err != nil {
		t.Error("ERROR")
	}
}

// TODO: test error responses from HTTP and from SUDOC
func TestGetRCRs(t *testing.T) {
	input := []string{"1", "2"}
	sc, _ := NewSudocClient(input, mockHttpFetcher{})

	want := make(map[string]library)
	want["100000001"] = library{"1", "100000001", "UNIV-1.1"}
	want["200000001"] = library{"2", "200000001", "UNIV-2.1"}
	want["200000002"] = library{"2", "200000002", "UNIV-2.2"}
	want["200000003"] = library{"2", "200000003", "UNIV-2.3"}

	got, err := sc.getRCRs(input)
	if err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}
