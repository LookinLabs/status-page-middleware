package model

type Request struct {
	Method  string                 `json:"method"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"`
}

type Service struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Type    string   `json:"type"`
	Status  string   `json:"status"`
	Request *Request `json:"request,omitempty"`
}
