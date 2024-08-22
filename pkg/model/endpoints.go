package model

type Request struct {
	Method  string                 `json:"method"`
	Headers map[string]string      `json:"headers"`
	Body    map[string]interface{} `json:"body"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Service struct {
	Name      string     `json:"name"`
	URL       string     `json:"url"`
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	Error     string     `json:"error,omitempty"`
	Request   *Request   `json:"request,omitempty"`
	BasicAuth *BasicAuth `json:"basic_auth,omitempty"`
}
