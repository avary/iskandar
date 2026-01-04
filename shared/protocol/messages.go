package protocol

type RegisterTunnelMessage struct {
	Subdomain string `json:"subdomain"`
}

type Message struct {
	Type    string            `json:"type"`
	Id      string            `json:"id"`
	Method  string            `json:"method,omitempty"`
	Path    string            `json:"path,omitempty"`
	Status  int               `json:"status,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    []byte            `json:"body,omitempty"`
	Done    bool              `json:"done,omitempty"`
}
