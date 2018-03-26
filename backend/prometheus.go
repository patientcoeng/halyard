package backend

import (
	"encoding/json"
	"fmt"
	"github.com/patientcoeng/halyard/alerting"
	"github.com/patientcoeng/halyard/api"
	"github.com/rs/zerolog/log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const PROMETHEUS_ERROR = "Halyard Prometheus Error"

type PrometheusBackend struct {
	Endpoint     string
	AlertManager *alerting.Manager
}

func (b *PrometheusBackend) Query(query string) float64 {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	urlValues := url.Values{}
	urlValues.Add("query", query)
	urlValues.Add("timeout", "30")

	queryUrl, _ := url.Parse(b.Endpoint)
	queryUrl.Path = "/api/v1/query"
	queryUrl.RawQuery = urlValues.Encode()

	resp, err := client.Get(queryUrl.String())
	if err != nil {
		msg := "Unable to query Prometheus; short circuiting evaluation"
		log.Error().Msg(msg)
		b.AlertManager.Trigger(PROMETHEUS_ERROR, msg)
		return math.NaN()
	}
	defer resp.Body.Close()

	var promResp api.PrometheusResponse
	err = json.NewDecoder(resp.Body).Decode(&promResp)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse Prometheus response for Query %s; short circuiting evaluation", query)
		log.Error().Msg(msg)
		b.AlertManager.Trigger(PROMETHEUS_ERROR, msg)
		return math.NaN()
	}

	var result float64
	if len(promResp.Data.Result) > 0 && len(promResp.Data.Result[0].Value) > 0 {
		stringResult := promResp.Data.Result[0].Value[1].(string)
		result, err = strconv.ParseFloat(stringResult, 64)
		if err != nil {
			msg := fmt.Sprintf("Unable to parse Prometheus response for Query %s; short circuiting evaluation", query)
			log.Error().Msg(msg)
			b.AlertManager.Trigger(PROMETHEUS_ERROR, msg)
			return math.NaN()
		}
	} else {
		msg := fmt.Sprintf("Unable to parse Prometheus response for Query %s; short circuiting evaluation", query)
		log.Error().Msg(msg)
		b.AlertManager.Trigger(PROMETHEUS_ERROR, msg)

		return math.NaN()
	}

	return result
}
