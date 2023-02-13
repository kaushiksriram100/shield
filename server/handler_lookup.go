package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	clientv3 "go.etcd.io/etcd/client/v3"
)

//Global Variable with in memory list of malwareDB
var MalwareInMemoryDB MalwareDB
var EtcdCli *clientv3.Client

//init() runs to load the models one time during startup, instead of loading for every-request
func init() {
	//In-memory malwareData. Load the data from a file
	MalwareInMemoryDB = LoadMalwareDataInMemory("../malwaredata/blacklist.json")
	EtcdCli = ConnectToEtcd()
}

func generateResponse(searchKey string, isHostBlacklisted bool) string {
	response := ShieldResponse{}
	response.Url = searchKey
	response.MalwareInfected = isHostBlacklisted

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(response)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	return buf.String()
}
func generateSearchKey(hostname_with_port, origin_path, raw_query string) string {
	searchKey := hostname_with_port + origin_path
	if len(raw_query) > 0 {
		searchKey = string(searchKey) + "?" + raw_query
	}
	return searchKey
}

//MalwareLookup uses the url parameters to check against a in-memory data
//Determines if the url is present in malware database.
func lookUpMalwareDB(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var isHostBlacklisted bool
	searchKey := generateSearchKey(params.ByName("hostname_with_port"), params.ByName("original_path"), r.URL.RawQuery)
	if _, ok := MalwareInMemoryDB[searchKey]; ok {
		isHostBlacklisted = true
	}
	//Construct response to send to client
	response := generateResponse(searchKey, isHostBlacklisted)
	fmt.Fprint(w, response)
}

func lookupMalwareEtcD(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var isHostBlacklisted bool
	searchKey := generateSearchKey(params.ByName("hostname_with_port"), params.ByName("original_path"), r.URL.RawQuery)
	//Connect to ETCD
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	resp, err := EtcdCli.Get(ctx, searchKey)
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Kvs) > 0 {
		isHostBlacklisted = true
	}
	//Construct response to send to client
	response := generateResponse(searchKey, isHostBlacklisted)
	fmt.Fprint(w, response)
}
