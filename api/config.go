package api

type EndpointMap map[string]string

type ASConfig struct {
	Period           int           `json:"period"`
	Namespace        string        `json:"namespace"`
	BackendEndpoints EndpointMap   `json:"backendEndpoints"`
	AlertConfig      ASAlertConfig `json:"alertConfig"`
}

type ASAlertConfig struct {
	Slack SlackConfig `json:"slack"`
}

type SlackConfig struct {
	WebhookURL string `json:"webhookURL"`
	Channel    string `json:"channel"`
}
