package main

//ShieldResponse struct to construct response
type ShieldResponse struct {
	Url             string `json:"url"`
	MalwareInfected bool   `json:"is_malware_infected"`
}

//ShieldAdminResponse struct to construct response
type ShieldAdminResponse struct {
	Url    string `json:"url"`
	Status string `json:"status"`
}
