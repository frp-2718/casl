package exl

import (
	"casl/requests"
	"os"
	"reflect"
	"testing"
)

type mockHttpFetcher struct{}

func (m mockHttpFetcher) Fetch(url string) ([]byte, error) {
	switch url {
	case almawsURL + "bibs?view=brief&expand=None&other_system_id=(PPN)" + "ppn_3_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_3_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case almawsURL + "bibs?view=brief&expand=None&other_system_id=(PPN)" + "ppn_2_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_2_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case almawsURL + "bibs?view=brief&expand=None&other_system_id=(PPN)" + "ppn_1_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_1_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case almawsURL + "bibs?view=brief&expand=None&other_system_id=(PPN)" + "ppn_0_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_0_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case "not_found":
		data, err := os.ReadFile("testdata/iln2rcr_not_found.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	default:
		return nil, nil
	}
}

func TestNewAlmaClient(t *testing.T) {
	tests_err := []struct {
		name    string
		apiKey  string
		baseURL string
		fetcher requests.Fetcher
	}{
		{"empty, empty, nil", "", "", nil},
		{"empty, url, fetcher", "", "url", mockHttpFetcher{}},
		{"key, url, nil", "key", "url", nil},
	}
	for _, test := range tests_err {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewAlmaClient(test.apiKey, test.baseURL, test.fetcher)
			if err == nil {
				t.Error("should return an error")
			}
		})
	}

	t.Run("no url", func(t *testing.T) {
		client, err := NewAlmaClient("key", "", mockHttpFetcher{})
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
		if client.baseURL != almawsURL {
			t.Errorf("want default base url, got %s", client.baseURL)
		}
	})

	t.Run("all ok", func(t *testing.T) {
		client, err := NewAlmaClient("key", "url", mockHttpFetcher{})
		if err != nil {
			t.Errorf("want no error, got %v", err)
		}
		if client.baseURL != "url" {
			t.Errorf("want baseURL == 'url', got %s", client.baseURL)
		}
		if client.apiKey != "key" {
			t.Errorf("want apiKey == 'key', got %s", client.apiKey)
		}
		if client.stats.bibs_req != 0 || client.stats.items_req != 0 {
			t.Errorf("want zeroed stats, got bibs_req = %d and items_req = %d", client.stats.bibs_req, client.stats.items_req)
		}
	})
}

// func (a *AlmaClient) GetLocations(ppn string) ([]*entities.AlmaLocation, error) {
func TestGetLocations(t *testing.T) {
}

func TestGetMMSFromPPN(t *testing.T) {
	client, _ := NewAlmaClient("key", "", mockHttpFetcher{})

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"3 MMS", "ppn_3_mms", []string{"mms1", "mms2", "mms3"}},
		// {"2 MMS", "ppn_2_mms", []string{"mms1", "mms3"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := client.getMMSfromPPN(test.input)
			if err != nil {
				t.Errorf("want no error, got %v", err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("want %v, got %v", test.want, got)
			}
		})
	}
}

// func (a *AlmaClient) getMMSfromPPN(ppn string) ([]string, error) {
//func TestGetMMSfromPPN(t *testing.T) {
//	alma, _ := New(nil, "apikey", almawsURL)
//	alma.fetchClient = &mockFetcher{}
//	var tests = []struct {
//		ppn      PPN
//		expected []MMS
//	}{
//		{"nonexistent", nil},
//		{"bibsCount1", []MMS{"mms1"}},
//		{"bibsCount2", []MMS{"mms2"}},
//		{"bibsCount3", []MMS{"mms1", "mms3"}},
//		{"bibsCount1noMMS", []MMS{}},
//	}
//	for _, test := range tests {
//		got, _ := alma.GetMMSfromPPN(test.ppn)
//		if !equalMMSslices(got, test.expected) {
//			t.Errorf("%s returned %v, expected %v", test.ppn, got, test.expected)
//		}
//	}
//}
