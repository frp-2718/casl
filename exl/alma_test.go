package exl

import (
	"casl/requests"
	"encoding/xml"
	"os"
	"reflect"
	"testing"
)

type mockHttpFetcher struct{}

const url_bibs = "bibs?view=brief&expand=None&other_system_id=(PPN)"

func (m mockHttpFetcher) Fetch(url string) ([]byte, error) {
	switch url {
	case almawsURL + url_bibs + "ppn_3_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_3_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case almawsURL + url_bibs + "ppn_1_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_1_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case almawsURL + url_bibs + "ppn_0_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_0_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case almawsURL + url_bibs + "ppn_00_mms" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_0_mms.xml")
		if err != nil {
			return nil, err
		}
		return data, nil
	case almawsURL + "bibs/" + "mms_items" + "/holdings/ALL/items?limit=100&apikey=key":
		data, err := os.ReadFile("testdata/mms_items.xml")
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
		{"1 MMS", "ppn_1_mms", []string{"mms2"}},
		{"0 MMS", "ppn_0_mms", []string{}},
		{"0 MMS", "ppn_00_mms", []string{}},
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

func TestGetItems(t *testing.T) {
	var name xml.Name
	items := []Item{
		{
			name,
			Holding{false, "mms_1", "CN_1"},
			ItemData{
				Status{"Item in place", "1"},
				Process{"", ""},
				Library{"Bibliothèque 1", "BIB_1"},
				Location{"Location 1", "LOC_1"},
			},
		},
		{
			name,
			Holding{false, "mms_2", "CN_2"},
			ItemData{
				Status{"Item in place", "1"},
				Process{"", ""},
				Library{"Bibliothèque 2", "BIB_2"},
				Location{"Location 2", "LOC_2"},
			},
		},
	}
	client, _ := NewAlmaClient("key", "", mockHttpFetcher{})
	res, err := client.getItems("mms_1")
	if err != nil {
		t.Errorf("got %v", err)
	}
	if !reflect.DeepEqual(res, items) {
		t.Error()
	}
}
