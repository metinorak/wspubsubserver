package entity

type WsPayload struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type WsResponse struct {
	Action  string `json:"action"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
