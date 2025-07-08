package status

type CreateStatusUpdateBody struct {
	Message string `json:"message"`
	Service string `json:"service"`
	Status  string `json:"status"`
}
