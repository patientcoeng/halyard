package alerting

import (
	"bytes"
	"encoding/json"
	"github.com/patientcoeng/halyard/api"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type SlackClient struct {
	WebhookURL string
	Channel    string
}

type SlackPayload struct {
	Channel     string             `json:"channel"`
	Attachments []SlackAttachments `json:"attachments"`
}

type SlackAttachments struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
}

func NewSlack(slackConfig api.SlackConfig) *SlackClient {
	return &SlackClient{
		WebhookURL: slackConfig.WebhookURL,
		Channel:    slackConfig.Channel,
	}
}

func (s *SlackClient) Trigger(summary string, detail string) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	payload := SlackPayload{
		Channel: s.Channel,
		Attachments: []SlackAttachments{{
			Title: summary,
			Text:  detail,
		},
		},
	}

	payloadStr, err := json.Marshal(payload)
	if err != nil {
		log.Error().Msgf("Unable to create message for Slack: %s - %s. %s", summary, detail, err.Error())
		return
	}

	_, err = client.Post(s.WebhookURL, "application/json", bytes.NewReader(payloadStr))
	if err != nil {
		log.Error().Msgf("Unable to send message to Slack: %s - %s. %s", summary, detail, err.Error())
	}
}
