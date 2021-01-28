package connectorLogger

type ConnectorLogQueryInput struct {
	Connector string `json:"connector,omitempty"`
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
	RawError  string `json:"raw_error,omitempty"`
	User      string `json:"user,omitempty"`
}
