package backend

import (
	"github.com/patientcoeng/halyard/alerting"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"math"
	"testing"
)

var validPromResponse = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"healthy_vault_prop","instance":"k8smon:443","job":"k8smon"},"value":[1500930297.86,"1"]}]}}`
var badNumber = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"healthy_vault_prop","instance":"k8smon:443","job":"k8smon"},"value":[1500930297.86,"abc123"]}]}}`
var invalidPromResponse = `{"status":"success","data":{"resultType":"vector","result":[]}}`
var badJson = `{"status":"success","data":{"resultType":"vector","result":[]}`

var manager = alerting.NewManager()

var backend = &PrometheusBackend{
	Endpoint:     "http://test",
	AlertManager: manager,
}

func TestValidPrometheusResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://test/api/v1/query",
		httpmock.NewStringResponder(200, validPromResponse))

	result := backend.Query("healthy_vault_prop")

	assert.Equal(t, 1.0, result)
}

func TestBadJSON(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://test/api/v1/query",
		httpmock.NewStringResponder(200, badJson))

	result := backend.Query("bad_query")

	assert.True(t, math.IsNaN(result))
}

func TestBadNumber(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://test/api/v1/query",
		httpmock.NewStringResponder(200, badNumber))

	result := backend.Query("bad_query")

	assert.True(t, math.IsNaN(result))
}

func TestInvalidPrometheusResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://test/api/v1/query",
		httpmock.NewStringResponder(200, invalidPromResponse))

	result := backend.Query("bad_query")

	assert.True(t, math.IsNaN(result))
}
