package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
