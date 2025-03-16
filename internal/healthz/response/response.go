package response

type HealthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}
