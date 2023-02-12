package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//ShieldResponse struct to construct response
type ShieldResponse struct {
	Url             string            `json:"url"`
	QueryStrings    map[string]string `json:"blacklisted_query_strings"`
	MalwareInfected bool              `json:"is_malware_infected"`
}

//MalwareLookup uses the url parameters to check against a in-memory data
//Determines if the url is present in malware database.
func LookUpMalwareDB(malwareDB MalwareDB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		shieldResponse := ShieldResponse{}
		hostPath := params.ByName("hostname_with_port") + params.ByName("original_path")

		var isHostBlacklisted bool
		blacklistedQueryStrings := make(map[string]string)
		if queryStringsInDB, ok := malwareDB[hostPath]; ok {
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

}
