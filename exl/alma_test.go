package exl

import (
	"casl/entities"
	"casl/requests"
	"encoding/xml"
	"os"
	"reflect"
	"sort"
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
	case almawsURL + url_bibs + "ppn_get_locations" + "&apikey=key":
		data, err := os.ReadFile("testdata/ppn_get_locations.xml")
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

func TestGetLocations(t *testing.T) {
	location_1 := entities.AlmaLocation{
		Library_name:  "Bibliothèque 1",
		Library_code:  "BIB_1",
		Location_name: "Location 1",
		Location_code: "LOC_1",
		Call_number:   "CN_1",
		NoDiscovery:   false,
		Items: []*entities.AlmaItem{
			{Process_name: "", Process_code: "", Status: "Item in place"},
		},
	}
	location_2 := entities.AlmaLocation{
		Library_name:  "Bibliothèque 2",
		Library_code:  "BIB_2",
		Location_name: "Location 2",
		Location_code: "LOC_2",
		Call_number:   "CN_2",
		NoDiscovery:   false,
		Items: []*entities.AlmaItem{
			{Process_name: "Acquisition", Process_code: "ACQ", Status: "Item in place"},
		},
	}
	locations := []*entities.AlmaLocation{&location_1, &location_2}

	client, _ := NewAlmaClient("key", "", mockHttpFetcher{})
	got, err := client.GetLocations("ppn_get_locations")
	if err != nil {
		t.Errorf("returned error %v", err)
	}
	sort.Slice(got, func(i, j int) bool {
		return got[i].Call_number < got[j].Call_number
	})
	if !reflect.DeepEqual(got, locations) {
		t.Errorf("want %v, got %v", locations, got)
	}
	_, err = client.GetLocations("ppn_3_mms.xml")
	_, err2 := client.GetLocations("ppn_00_mms.xml")
	if err == nil || err2 == nil {
		t.Error("want error, got ok")
	}
}

func TestGetFilteredLocations(t *testing.T) {
	locations := []*entities.AlmaLocation{
		{
			Library_name:  "Bibliothèque 1",
			Library_code:  "BIB_1",
			Location_name: "Location 1",
			Location_code: "LOC_1",
			Call_number:   "CN_1",
			NoDiscovery:   false,
			Items: []*entities.AlmaItem{
				{Process_name: "", Process_code: "", Status: "Item in place"},
			},
		},
	}

	client, _ := NewAlmaClient("key", "", mockHttpFetcher{})
	got, err := client.GetFilteredLocations("ppn_get_locations", []string{"BIB_1"})
	if err != nil {
		t.Errorf("returned error %v", err)
	}
	if !reflect.DeepEqual(got, locations) {
		t.Errorf("want %v, got %v", locations, got)
	}
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
	name.Local = "item"
	items := []Item{
		{
			name,
			Holding{false, "mms_1", "CN_1"},
			ItemData{
				Status{"Item in place", "1"},
				Process{},
				Library{"Bibliothèque 1", "BIB_1"},
				Location{"Location 1", "LOC_1"},
			},
		},
		{
			name,
			Holding{false, "mms_2", "CN_2"},
			ItemData{
				Status{"Item in place", "1"},
				Process{"Acquisition", "ACQ"},
				Library{"Bibliothèque 2", "BIB_2"},
				Location{"Location 2", "LOC_2"},
			},
		},
	}
	client, _ := NewAlmaClient("key", "", mockHttpFetcher{})
	got, err := client.getItems("mms_items")
	if err != nil {
		t.Errorf("got %v", err)
	}
	if !reflect.DeepEqual(got, items) {
		t.Errorf("want %v, got %v", items, got)
	}
}

func TestStats(t *testing.T) {
	client, _ := NewAlmaClient("key", "", mockHttpFetcher{})
	bibs, items, total := getStats(client)
	if bibs != 0 || items != 0 || total != 0 {
		t.Errorf("want 0 0 0, got %d %d %d", bibs, items, total)
	}
	client.getItems("mms_1")
	bibs, items, total = getStats(client)
	if bibs != 0 || items != 1 || total != 1 {
		t.Errorf("want 0 1 1, got %d %d %d", bibs, items, total)
	}
	client.GetLocations("ppn_get_locations")
	bibs, items, total = getStats(client)
	if bibs != 1 || items != 2 || total != 3 {
		t.Errorf("want 1 2 3, got %d %d %d", bibs, items, total)
	}
	client.GetFilteredLocations("ppn_get_locations", []string{"mms_items"})
	bibs, items, total = getStats(client)
	if bibs != 2 || items != 3 || total != 5 {
		t.Errorf("want 2 3 5, got %d %d %d", bibs, items, total)
	}
}

func getStats(client *AlmaClient) (int, int, int) {
	return client.Stats("bibs"), client.Stats("items"), client.Stats("total")
}
