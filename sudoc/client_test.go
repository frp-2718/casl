package sudoc

import (
	"casl/requests"
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

func TestNewSudocClient(t *testing.T) {
	tests_err := []struct {
		name    string
		input   []string
		fetcher requests.Fetcher
	}{
		{"nil, nil", nil, nil},
		{"empty, nil", []string{}, nil},
		{"empty, ok", []string{}, mockHttpFetcher{}},
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

	sc := SudocClient{rcrs, 0, mockHttpFetcher{}}
	client, err := NewSudocClient([]string{"1", "2"}, nil)
	t.Run("ok, nil", func(t *testing.T) {
		if err != nil {
			t.Error("got error, want properly built SudocClient")
		}
		if !reflect.DeepEqual(client, sc) {
			t.Error("error")
		}
	})
	//{"ok, nil", []string{"1", "2"}, nil},
	// pas d'iln
	// pas de client
	// ni iln ni client
	// arguments ok --> rcr ok (déjà testé par un autre test) ; stats = 0 ; fetcher = fetcher
}

// func TestGetFilteredLocations(t *testing.T) {
// 	sc, _ := NewSudocClient([]string{}, mockHttpFetcher{})
// 	sc.GetFilteredLocations("123456789", []string{})
// }

// func (sc *SudocClient) GetFilteredLocations(ppn string, rcrs []string) ([]*entities.SudocLocation, error) {
// 	locations, err := sc.GetLocations(ppn)
// }

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
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}
