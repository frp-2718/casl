package sudoc

import (
	"os"
	"reflect"
	"testing"
)

func TestDecodeRCR(t *testing.T) {
	data_ok, err := os.ReadFile("testdata/iln2rcr.xml")
	if err != nil {
		t.Fatal(err)
	}
	data_error, err := os.ReadFile("testdata/iln2rcr_notfound.xml")
	if err != nil {
		t.Fatal(err)
	}

	want_ok := make(map[string]library)
	want_ok["100000001"] = library{"1", "100000001", "UNIV-1.1"}
	want_ok["200000001"] = library{"2", "200000001", "UNIV-2.1"}
	want_ok["200000002"] = library{"2", "200000002", "UNIV-2.2"}
	want_ok["200000003"] = library{"2", "200000003", "UNIV-2.3"}

	got_ok, err := decodeRCR(data_ok)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got_ok, want_ok) {
		t.Errorf("want %v, got %v", want_ok, got_ok)
	}

	_, err = decodeRCR(data_error)
	if err == nil {
		t.Error("want error for 'null xml' response")
	}
}
