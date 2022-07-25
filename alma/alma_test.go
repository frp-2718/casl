package alma

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAlma(t *testing.T) {
	alma, err := New(nil, "apikey", "base/url")
	if err != nil {
		t.Fatalf(`New(nil, "apikey", "base/url") returned error: (%s)`, err)
	}
	if alma == nil {
		t.Fatal("New(nil, 'apikey', 'base/url') returned nil")
	}
	if alma.client == nil {
		t.Error("New: alma.client is nil")
	}
	assertEqualString(t, alma.apiKey, "apikey")
	assertEqualString(t, alma.baseURL, "base/url")

	_, err = New(nil, "", "base/url")
	if err == nil {
		t.Error(`New(nil, "", "base/url") returned no error`)
	}
	_, err = New(nil, "apikey", "")
	if err != nil {
		t.Errorf(`New(nil, "apikey", "") returned error: %v`, err)
	}
}

func TestFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/400_general_error":
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>GENERAL_ERROR</errorCode><errorMessage>general error</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/400_unauthorized":
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>UNAUTHORIZED</errorCode><errorMessage>unauthorized</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/403_unauthorized":
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>UNAUTHORIZED</errorCode><errorMessage>unauthorized</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/403_invalid_request":
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>INVALID_REQUEST</errorCode><errorMessage>ignored</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/403_request_too_large":
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>REQUEST_TOO_LARGE</errorCode><errorMessage>request too large</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/403_forbidden":
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>FORBIDDEN</errorCode><errorMessage>forbidden</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/429_per_second_threshold":
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>PER_SECOND_THRESHOLD</errorCode><errorMessage>per second threshold</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/429_daily_threshold":
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>DAILY_THRESHOLD</errorCode><errorMessage>daily threshold</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/503_routing_error":
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>ROUTING_ERROR</errorCode><errorMessage>routing error</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/500_general_error":
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>GENERAL_ERROR</errorCode><errorMessage>general error</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/400_401652":
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>401652</errorCode><errorMessage>401652</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/400_402203":
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><web_service_result xmlns=""><errorsExist>true</errorsExist><errorList><error><errorCode>402203</errorCode><errorMessage>402203</errorMessage><trackingId></trackingId></error></errorList></web_service_result>`))
			case "/200_count_0":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><bibs total_record_count="0"/>`))
			case "/200_count_1":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><bibs total_record_count="1"/>`))
			case "/200_count_2":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><bibs total_record_count="2"/>`))
			case "/4XX":
				w.WriteHeader(http.StatusBadRequest)
			case "/5XX":
				w.WriteHeader(http.StatusInternalServerError)
			default:
				w.WriteHeader(http.StatusOK)
			}

		}))
	defer server.Close()

	alma, _ := New(nil, "apikey", server.URL)
	alma.fetchClient = &almaFetcher{client: alma.client}

	var errorTests = []struct {
		url           string
		expectedError error
	}{
		{"/400_general_error", &InvalidRequestError{errorMessage: "general error"}},
		{"/400_unauthorized", &UnauthorizedError{errorMessage: "unauthorized"}},
		{"/403_unauthorized", &UnauthorizedError{errorMessage: "unauthorized"}},
		{"/403_invalid_request", &UnauthorizedError{errorMessage: "ignored"}},
		{"/403_request_too_large", &InvalidRequestError{errorMessage: "request too large"}},
		{"/403_forbidden", &UnauthorizedError{errorMessage: "forbidden"}},
		//{"/429_per_second_threshold", &ThresholdError{errorMessage: "per second threshold"}},
		//{"/429_daily_threshold", &ThresholdError{errorMessage: "daily threshold"}},
		{"/503_routing_error", &ServerError{errorMessage: "routing error"}},
		{"/500_general_error", &UnauthorizedError{errorMessage: "general error"}},
		{"/400_401652", &InvalidRequestError{errorMessage: "401652"}},
		{"/400_402203", &InvalidRequestError{errorMessage: "402203"}},
		{"/4XX", &FetchError{errorMessage: ""}},
		{"/5XX", &FetchError{errorMessage: ""}},
		{"", &InvalidRequestError{errorMessage: "URL cannot be empty"}},
	}
	for i, test := range errorTests {
		var err error
		if test.url == "" {
			_, err = alma.Fetch(test.url)
		} else {
			_, err = alma.Fetch(server.URL + test.url)
		}
		if err != nil {
			if reflect.TypeOf(err) != reflect.TypeOf(test.expectedError) {
				t.Errorf("[Test %d] url %q returned %v error ; should be %v",
					i, test.url, reflect.TypeOf(err), reflect.TypeOf(test.expectedError))
			}
			if err.Error() != test.expectedError.Error() {
				t.Errorf("[Test %d] url %q returned %q error message ; should be %q",
					i, test.url, err, test.expectedError)
			}
		}
	}

	// tests for status 200 responses
	e := &NotFoundError{errorMessage: "identifier not found"}
	_, err := alma.Fetch(server.URL + "/200_count_0")
	if err == nil || reflect.TypeOf(err) != reflect.TypeOf(e) {
		t.Errorf("url %q returned %v error ; should be %q", "/200_count_0", reflect.TypeOf(err), reflect.TypeOf(e))
	}
	if err != nil && err.Error() != "identifier not found" {
		t.Errorf("url %q returned %q error message ; should be %q", "/200_count_0", err, "identifier not found")
	}

	_, err = alma.Fetch(server.URL + "/200_count_1")
	if err != nil {
		t.Errorf("url %q returned error", "/200/count_1")
	}
	_, err = alma.Fetch(server.URL + "/200_count_2")
	if err != nil {
		t.Errorf("url %q returned error", "/200/count_2")
	}
}

func TestGetMMSfromPPN(t *testing.T) {
	alma, _ := New(nil, "apikey", almawsURL)
	alma.fetchClient = &mockFetcher{}
	var tests = []struct {
		ppn      PPN
		expected []MMS
	}{
		{"nonexistent", nil},
		{"bibsCount1", []MMS{"mms1"}},
		{"bibsCount2", []MMS{"mms2"}},
		{"bibsCount3", []MMS{"mms1", "mms3"}},
		{"bibsCount1noMMS", []MMS{}},
	}
	for _, test := range tests {
		got, _ := alma.GetMMSfromPPN(test.ppn)
		if !equalMMSslices(got, test.expected) {
			t.Errorf("%s returned %v, expected %v", test.ppn, got, test.expected)
		}
	}
}

func TestGetHoldings(t *testing.T) {
	alma, _ := New(nil, "apikey", almawsURL)
	alma.fetchClient = &mockFetcher{}
	l1 := Holding{Library: "LIBRARY1", Location: "LOCATION1", CallNumber: "CALL1"}
	l2 := Holding{Library: "LIBRARY2", Location: "LOCATION2", CallNumber: "CALL2"}
	r0 := []Holding{}
	r1 := []Holding{l1}
	r2 := []Holding{l1, l2}
	var tests = []struct {
		mms      MMS
		expected []Holding
	}{
		{"mms0", r0},
		{"mms1", r1},
		{"mms2", r2},
	}
	for _, test := range tests {
		got, _ := alma.GetHoldings(test.mms)
		if !equalHoldingSlices(got, test.expected) {
			t.Errorf("%s returned %v, expected %v", test.mms, got, test.expected)
		}
	}
}

func TestGetHoldingsFromPPN(t *testing.T) {
	alma, _ := New(nil, "apikey", almawsURL)
	alma.fetchClient = &mockFetcher{}
	l1 := Holding{Library: "LIBRARY1", Location: "LOCATION1", CallNumber: "CALL1"}
	l2 := Holding{Library: "LIBRARY2", Location: "LOCATION2", CallNumber: "CALL2"}
	r1 := []Holding{l1}
	r2 := []Holding{l1, l2}
	var tests = []struct {
		ppn      PPN
		expected []Holding
	}{
		{"bibsCount1", r1},
		{"bibsCount2", r2},
	}
	for _, test := range tests {
		got, _ := alma.GetHoldingsFromPPN(test.ppn)
		if !equalHoldingSlices(got, test.expected) {
			t.Errorf("%s returned %v, expected %v", test.ppn, got, test.expected)
		}
	}
	_, err := alma.GetHoldingsFromPPN("invalid")
	if err == nil {
		t.Error("invalid ppn returned no error ; expected NotFoundError")
	}
}

func equalMMSslices(s1, s2 []MMS) bool {
	if s1 != nil && s2 != nil && len(s1) == len(s2) {
		for i := range s1 {
			if s1[i] != s2[i] {
				return false
			}
		}
		return true
	}
	return s1 == nil && s2 == nil
}

func equalHoldingSlices(s1, s2 []Holding) bool {
	if s1 != nil && s2 != nil && len(s1) == len(s2) {
		for i := range s1 {
			if s1[i].Library != s2[i].Library ||
				s1[i].Location != s2[i].Location ||
				s1[i].CallNumber != s2[i].CallNumber {
				return false
			}
		}
		return true
	}
	return s1 == nil && s2 == nil
}

func assertEqualString(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Errorf("%q should equal %q", s1, s2)
	}
}

type mockFetcher struct{}

func (f *mockFetcher) Fetch(url string) ([]byte, error) {
	bibs0 := almawsURL + "bibs?view=brief&expand=None&other_system_id=" + "nonexistent" + "&apikey=" + "apikey"
	bibs1 := almawsURL + "bibs?view=brief&expand=None&other_system_id=" + "bibsCount1" + "&apikey=" + "apikey"
	bibs2 := almawsURL + "bibs?view=brief&expand=None&other_system_id=" + "bibsCount2" + "&apikey=" + "apikey"
	bibs3 := almawsURL + "bibs?view=brief&expand=None&other_system_id=" + "bibsCount3" + "&apikey=" + "apikey"
	bibs1noMMS := almawsURL + "bibs?view=brief&expand=None&other_system_id=" + "bibsCount1noMMS" + "&apikey=" + "apikey"
	mms0 := almawsURL + "bibs/" + "mms0" + "/holdings?apikey=" + "apikey"
	mms1 := almawsURL + "bibs/" + "mms1" + "/holdings?apikey=" + "apikey"
	mms2 := almawsURL + "bibs/" + "mms2" + "/holdings?apikey=" + "apikey"
	switch url {
	case bibs0:
		return nil, &NotFoundError{errorMessage: "identifier not found"}
	case bibs1:
		return bibsCount1, nil
	case bibs2:
		return bibsCount2, nil
	case bibs3:
		return bibsCount3, nil
	case bibs1noMMS:
		return bibsCount1noMMS, nil
	case mms0:
		return holdingsCount0, nil
	case mms1:
		return holdingsCount1, nil
	case mms2:
		return holdingsCount2, nil
	default:
		return nil, &NotFoundError{errorMessage: "identifier not found"}
	}
}
