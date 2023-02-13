package main

//ShieldResponse struct to construct response
type ShieldResponse struct {
	Url             string `json:"url"`
	MalwareInfected bool   `json:"is_malware_infected"`
}
