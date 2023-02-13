package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

//TestLookUpMalwareDB performs unit tests.
func TestLookUpMalwareDB(t *testing.T) {

	testCase := "/urlinfo/1/google.com:8080/search?v=2"
	expectedResponse := "{\"url\":\"google.com:8080/search?v=2\",\"is_malware_infected\":true}\n"

	router := httprouter.New()
	router.Handle("GET", "/urlinfo/1/:hostname_with_port/*original_path", LookUpMalwareDB)
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", testCase, nil)
	router.ServeHTTP(rr, req)

	/*
		t.Log(req.URL.Path)
		t.Log(req.URL)
		t.Log(rr.Code)
	*/

	//Verify Results
	if rr.Code != 200 {
		t.Error("Return Code does not match")
	}
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expectedResponse) {
		t.Log(rr.Body.String())
		t.Log(expectedResponse)
		t.Error("Response Body Does NOT match")
	}
}

func TestLookupMalwareEtcD(t *testing.T) {

	testCase := "/urlinfo/1/google.com:8080/search?v=2"
	expectedResponse := "{\"url\":\"google.com:8080/search?v=2\",\"is_malware_infected\":true}\n"

	router := httprouter.New()
	router.Handle("GET", "/urlinfo/1/:hostname_with_port/*original_path", LookupMalwareEtcD)
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", testCase, nil)
	router.ServeHTTP(rr, req)

	//t.Log(req.URL.Path)
	//t.Log(req.URL)
	//t.Log(rr.Code)

	//Verify Results
	if rr.Code != 200 {
		t.Error("Return Code does not match")
	}
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(expectedResponse) {
		t.Log(rr.Body.String())
		t.Log(expectedResponse)
		t.Error("Response Body Does NOT match")
	}
}
