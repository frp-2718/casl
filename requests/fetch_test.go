package requests

import (
	"net/http"
	"net/http/httptest"
	"strconv"
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
	max_params := 1
	var ppns []string
	var result [][]byte
	for i := 0; i < 100; i++ {
		ppn := "ppn" + strconv.Itoa(i)
		ppns = append(ppns, ppn)
		result = append(result, []byte(ppn))
	}
	var tests = []struct {
		ppns []string
		want [][]byte
	}{
		{[]string{}, [][]byte{}},
		{[]string{"ppn1"}, [][]byte{[]byte("ppn1")}},
		{[]string{"ppn1", "ppn2"}, [][]byte{[]byte("ppn1"), []byte("ppn2")}},
		{ppns, result},
	}

	for i, test := range tests {
		wanted := wantedLength(len(test.ppns), max_params)
		got := fetchBatch(test.ppns, max_params, mockHttpRequester)
		if len(got) != wanted {
			t.Errorf("[%d] fetchBatch(%v) result has wrong lentgh : got %v", i, test.ppns, got)
		}
		if equalSlicesOfSlicesOfBytes(got, test.want) {
			t.Errorf("[%d] fetchBatch(%v) returned %v : want %v", i, test.ppns, got, test.want)
		}
	}
}

func TestFetchBatchConcurrent(t *testing.T) {
	max_params := 1
	var ppns []string
	var result [][]byte
	for i := 0; i < 100; i++ {
		ppn := "ppn" + strconv.Itoa(i)
		ppns = append(ppns, ppn)
		result = append(result, []byte(ppn))
	}
	var tests = []struct {
		ppns []string
		want [][]byte
	}{
		{[]string{}, [][]byte{}},
		{[]string{"ppn1"}, [][]byte{[]byte("ppn1")}},
		{[]string{"ppn1", "ppn2"}, [][]byte{[]byte("ppn1"), []byte("ppn2")}},
		{ppns, result},
	}

	for i, test := range tests {
		wanted := wantedLength(len(test.ppns), max_params)
		got := fetchBatchConcurrent(test.ppns, max_params, mockHttpRequester)
		if len(got) != wanted {
			t.Errorf("[%d] fetchBatch(%v) result has wrong lentgh : got %v", i, test.ppns, got)
		}
		if equalSlicesOfSlicesOfBytes(got, test.want) {
			t.Errorf("[%d] fetchBatch(%v) returned %v : want %v", i, test.ppns, got, test.want)
		}
	}
}

func BenchmarkFetchBatch(b *testing.B) {
	var ppns []string
	for i := 0; i < 1000; i++ {
		ppns = append(ppns, "ppn"+strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		fetchBatch(ppns, 1, mockHttpRequester)
	}
}

func BenchmarkFetchBatchConcurrent(b *testing.B) {
	var ppns []string
	for i := 0; i < 1000; i++ {
		ppns = append(ppns, "ppn"+strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		fetchBatchConcurrent(ppns, 1, mockHttpRequester)
	}
}

// helpers
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

func wantedLength(n, max int) int {
	if n == 0 {
		return 1
	}
	result := n / max
	if n%max == 0 {
		return result
	}
	return result + 1
}
