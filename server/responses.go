package main

//ShieldResponse struct to construct response
type ShieldResponse struct {
	Url             string            `json:"url"`
	QueryStrings    map[string]string `json:"blacklisted_query_strings"`
	MalwareInfected bool              `json:"is_malware_infected"`
}
