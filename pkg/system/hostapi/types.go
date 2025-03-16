package api

type ServiceInfo struct {
	// Service display name
	DisplayName string `json:"DisplayName"`

	// Service start type
	StartType string `json:"StartType"`

	// Service status
	Status string `json:"Status"`
}
