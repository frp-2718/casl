package sudoc

import (
	"os"
	"reflect"
	"testing"
)

func TestDecodeRCR(t *testing.T) {
	data, err := os.ReadFile("testdata/iln2rcr.xml")
	if err != nil {
		t.Fatal(err)
	}

	want := make(map[string]library)
	want["100000001"] = library{"1", "100000001", "UNIV-1.1"}
	want["200000001"] = library{"2", "200000001", "UNIV-2.1"}
	want["200000002"] = library{"2", "200000002", "UNIV-2.2"}
	want["200000003"] = library{"2", "200000003", "UNIV-2.3"}

	got, err := decodeRCR(data)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
}

//func decodeRCR(data []byte) (map[string]library, error) {
