package main

import (
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

//MalwareLookup uses the url parameters to check against a in-memory data
//Determines if the url is present in malware database.
func LookUpMalwareDB(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	shieldResponse := ShieldResponse{}
	hostPath := params.ByName("hostname_with_port") + params.ByName("original_path")

	var isHostBlacklisted bool
	blacklistedQueryStrings := make(map[string]string)
	if queryStringsInDB, ok := MalwareInMemoryDB[hostPath]; ok {
		if len(queryStringsInDB) == 0 {
			isHostBlacklisted = true
		} else {
			//Compare queryStrings in request with DB
			for queryStringKeyInRequest, queryStringValuesInRequest := range r.URL.Query() {
				if queryStringValueInDB, ok := queryStringsInDB[queryStringKeyInRequest]; ok {
					for _, queryStringValueInRequest := range queryStringValuesInRequest {
						if queryStringValueInRequest == queryStringValueInDB {
							isHostBlacklisted = true
							blacklistedQueryStrings[queryStringKeyInRequest] = queryStringValueInRequest
						}
					}
				}
			}
		}
	}
	//Construct response to send to client
	shieldResponse.Url = hostPath
	shieldResponse.QueryStrings = blacklistedQueryStrings

	if isHostBlacklisted {
		shieldResponse.MalwareInfected = true
	} else {
		shieldResponse.MalwareInfected = false
	}
	//Construct a response Json to send back to requestor
	responseJSON, err := json.Marshal(shieldResponse)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	fmt.Fprint(w, string(responseJSON))
}

func LookupMalwareEtcD(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	shieldResponse := ShieldResponse{}
	hostPath := params.ByName("hostname_with_port") + params.ByName("original_path")

	var isHostBlacklisted bool
	blacklistedQueryStrings := make(map[string]string)

	//Connect to ETCD
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	resp, err := EtcdCli.Get(ctx, hostPath)
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.Kvs) > 0 {
		for _, ev := range resp.Kvs {
			if string(ev.Key) == hostPath {
				//Unpack the query strings
				queryStringsInDB := make(map[string]string)
				err := json.Unmarshal([]byte(ev.Value), &queryStringsInDB)
				if err != nil {
					log.Printf("ERROR: %s", err.Error())
				}
				if len(queryStringsInDB) == 0 {
					isHostBlacklisted = true
				} else {
					//Compare queryStrings in request with DB
					for queryStringKeyInRequest, queryStringValuesInRequest := range r.URL.Query() {
						if queryStringValueInDB, ok := queryStringsInDB[queryStringKeyInRequest]; ok {
							for _, queryStringValueInRequest := range queryStringValuesInRequest {
								if queryStringValueInRequest == queryStringValueInDB {
									isHostBlacklisted = true
									blacklistedQueryStrings[queryStringKeyInRequest] = queryStringValueInRequest
								}
							}
						}
					}
				}
			}
		}
	}
	//Construct response to send to client
	shieldResponse.Url = hostPath
	shieldResponse.QueryStrings = blacklistedQueryStrings

	if isHostBlacklisted {
		shieldResponse.MalwareInfected = true
	} else {
		shieldResponse.MalwareInfected = false
	}
	//Construct a response Json to send back to requestor
	responseJSON, err := json.Marshal(shieldResponse)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	fmt.Fprint(w, string(responseJSON))
}
