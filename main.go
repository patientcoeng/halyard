package main

import (
	"encoding/json"
	"flag"
	"github.com/patientcoeng/halyard/alerting"
	"github.com/patientcoeng/halyard/api"
	"github.com/patientcoeng/halyard/autoscaler"
	"github.com/patientcoeng/halyard/k8s"
	"github.com/rs/zerolog/log"
	"os"
	"runtime"
	"time"
)

func main() {
	var configPath string
	var config api.ASConfig

	// Parse flags
	flag.StringVar(&configPath, "config", "Config.json", "JSON-encoded configuration")
	flag.Parse()

	// Parse config
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Error().Err(err).Msg("Unable to open config file")
		return
	}
	defer configFile.Close()

	if err = json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Error().Err(err).Msg("Unable to decode config file")
		return
	}

	// Initialization message
	log.Info().Msg("***** HALYARD *****")
	log.Info().Msgf("Go Runtime: %s", runtime.Version())

	client := k8s.NewK8S(config.Namespace)

	alertManager := alerting.NewManager()
	if config.AlertConfig.Slack.WebhookURL != "" && config.AlertConfig.Slack.Channel != "" {
		alertManager.AddAlert(alerting.NewSlack(config.AlertConfig.Slack))
	}

	// Loop indefinitely, at defined period
	for range time.Tick(time.Duration(config.Period) * time.Second) {
		// Check for new pods
		annotations, err := client.GetDeployAnnotations()
		if err != nil {
			log.Error().Msgf("Error getting annotations: %s", err)
			continue
		}

		rules := autoscaler.CreateASRules(annotations, config.BackendEndpoints, alertManager)

		commands := autoscaler.Scale(rules)
		for _, cmd := range commands {
			log.Info().Msgf("%s", cmd.String())
		}

		// Update K8S
		err = client.UpdateReplicas(commands)
		if err != nil {
			log.Error().Msgf("Error updating deployments: %s", err)
			continue
		}
	}
}
