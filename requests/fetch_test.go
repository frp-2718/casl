package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/404":
			w.WriteHeader(http.StatusNotFound)
		case "/400":
			w.WriteHeader(http.StatusBadRequest)
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
		case "/empty":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte{})
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
	}))
	defer server.Close()

	var tests = []struct {
		url  string
		want []byte
	}{
		{"/404", []byte{}},
		{"/400", []byte{}},
		{"/500", []byte{}},
		{"/empty", []byte{}},
		{"/200", []byte("ok")},
	}

	for i, test := range tests {
		got := fetch(server.URL + test.url)
		if string(got) != string(test.want) {
			t.Errorf("[%d] fetch(%q) returned %q : want %q", i, test.url, got, test.want)
		}
	}
}

func mockHttpRequester(url string) []byte {
	time.Sleep(10 * time.Millisecond)
	return []byte(url)
}

func TestFetchBatch(t *testing.T) {
	var tests = []struct {
		ppns []string
		want [][]byte
	}{
		{[]string{}, [][]byte{}},
		{[]string{"ppn1"}, [][]byte{[]byte("ppn1")}},
		{[]string{"ppn1", "ppn2"}, [][]byte{[]byte("ppn1"), []byte("ppn2")}},
	}

	for i, test := range tests {
		got := fetchBatch(test.ppns, mockHttpRequester)
		if len(got) != len(test.ppns)/max_params+1 {
			t.Errorf("[%d] fetchBatch(%v) result has wrong lentgh : got %v", i, test.ppns, got)
		}
		if equalSlicesOfSlicesOfBytes(got, test.want) {
			t.Errorf("[%d] fetchBatch(%v) returned %v : want %v", i, test.ppns, got, test.want)
		}
	}
}

func TestFetchBatchConcurrent(t *testing.T) {
	var tests = []struct {
		ppns []string
		want [][]byte
	}{
		{[]string{}, [][]byte{}},
		{[]string{"ppn1"}, [][]byte{[]byte("ppn1")}},
		{[]string{"ppn1", "ppn2"}, [][]byte{[]byte("ppn1"), []byte("ppn2")}},
	}

	for i, test := range tests {
		got := fetchBatchConcurrent(test.ppns, mockHttpRequester)
		if len(got) != len(test.ppns)/max_params+1 {
			t.Errorf("[%d] fetchBatch(%v) result has wrong lentgh : got %v", i, test.ppns, got)
		}
		if equalSlicesOfSlicesOfBytes(got, test.want) {
			t.Errorf("[%d] fetchBatch(%v) returned %v : want %v", i, test.ppns, got, test.want)
		}
	}
}

func BenchmarkFetchBatch(b *testing.B) {
	ppns := []string{"111111111", "222222222", "333333333"}
	for i := 0; i < b.N; i++ {
		fetchBatch(ppns, mockHttpRequester)
	}
}

func BenchmarkFetchBatchConcurrent(b *testing.B) {
	ppns := []string{"111111111", "222222222", "333333333"}
	for i := 0; i < b.N; i++ {
		fetchBatchConcurrent(ppns, mockHttpRequester)
	}
}

func equalSlicesOfSlicesOfBytes(s1, s2 [][]byte) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, s := range s1 {
		if string(s) != string(s2[i]) {
			return false
		}
	}
	return true
}
