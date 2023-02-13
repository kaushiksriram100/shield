package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

//TestLookUpMalwareDB performs unit tests for in-memoryDB.
func TestLookUpMalwareDB(t *testing.T) {

	testCase := "/urlinfo/1/google.com:8080/search?v=2"
	expectedResponse := "{\"url\":\"google.com:8080/search?v=2\",\"is_malware_infected\":true}\n"

	router := httprouter.New()
	router.Handle("GET", "/urlinfo/1/:hostname_with_port/*original_path", lookUpMalwareDB)
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
		t.Log("Current Response", rr.Body.String())
		t.Log("Expected Response", expectedResponse)
		t.Error("Response Body Does NOT match")
	}
}

//TestLookupMalwareEtcD performs unit tests for ETCD lookup
func TestLookupMalwareEtcD(t *testing.T) {

	testCase := "/urlinfo/1/google.com:8080/search?v=3"
	expectedResponse := "{\"url\":\"google.com:8080/search?v=3\",\"is_malware_infected\":false}\n"

	router := httprouter.New()
	router.Handle("GET", "/urlinfo/1/:hostname_with_port/*original_path", lookupMalwareEtcD)
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
		t.Log("Current Response: ", rr.Body.String())
		t.Log("Expected Response:", expectedResponse)
		t.Error("Response Body Does NOT match")
	}
}

func TestPutMalwareUrlToEtcD(t *testing.T) {

	testCase := "/urlinfo/1/google.com:8080/search?v=3"

	router := httprouter.New()
	router.Handle("PUT", "/urlinfo/1/:hostname_with_port/*original_path", putMalwareUrlToEtcD)
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", testCase, nil)
	router.ServeHTTP(rr, req)

	//t.Log(req.URL.Path)
	//t.Log(req.URL)
	//t.Log(rr.Code)

	//Verify Results
	if rr.Code != 200 {
		t.Error("Return Code does not match")
	}
}

//TestDeleteMalwareUrlToEtcD tests for DELETE
func TestDeleteMalwareUrlToEtcD(t *testing.T) {

	testCase := "/urlinfo/1/google.com:8080/search?v=3"

	router := httprouter.New()
	router.Handle("DELETE", "/urlinfo/1/:hostname_with_port/*original_path", deleteMalwareUrlInEtcD)
	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", testCase, nil)
	router.ServeHTTP(rr, req)

	//t.Log(req.URL.Path)
	//t.Log(req.URL)
	//t.Log(rr.Code)

	//Verify Results
	if rr.Code != 200 {
		t.Error("Return Code does not match")
	}
}

//BenchmarkShieldServer to benchmark the http server
func BenchmarkShieldServer(b *testing.B) {

	//Create a new router and route handle
	router := httprouter.New()
	router.Handle("GET", "/urlinfo/1/:hostname_with_port/*original_path", lookupMalwareEtcD)

	//Create a http server using the router above
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	//Open listener
	go server.ListenAndServe()
	defer server.Close()
	//start benchmark tests
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := http.Get("http://localhost:8080/urlinfo/1/google.com:8080/search?v=2")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			b.Fatalf("Expected status code 200, but got %d", res.StatusCode)
		}
		_, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			b.Fatalf("Failed to read response: %v", err)
		}
	}
}
