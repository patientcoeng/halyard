package autoscaler

import (
	"github.com/patientcoeng/halyard/alerting"
	"github.com/patientcoeng/halyard/api"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"math"
	"testing"
)

var validPromResponse = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"healthy_vault_prop","instance":"k8smon:443","job":"k8smon"},"value":[1500930297.86,"100"]}]}}`

var manager = alerting.NewManager()

var rules = []ASRule{
	ASRule{
		Query:         "test",
		Backend:       "prometheus",
		Resource:      "testPod",
		Result:        0.0,
		Target:        10.0,
		MinReplicas:   0,
		MaxReplicas:   math.MaxInt32,
		ScalingPolicy: "linear",
		Endpoint:      "http://test",
		AlertManager:  manager,
	},
}

var constrainedRules = []ASRule{
	ASRule{
		Query:         "test",
		Backend:       "prometheus",
		Resource:      "testPod",
		Result:        0.0,
		Target:        10.0,
		MinReplicas:   5,
		MaxReplicas:   6,
		ScalingPolicy: "linear",
		Endpoint:      "http://test",
		AlertManager:  manager,
	},
}

var expectedCommands = []api.ASCommand{
	api.ASCommand{
		Resource:    "testPod",
		Cmd:         10.0,
		MinReplicas: 0,
		MaxReplicas: math.MaxInt32,
	},
}

var expectedConstrainedCommands = []api.ASCommand{
	api.ASCommand{
		Resource:    "testPod",
		Cmd:         6.0,
		MinReplicas: 5,
		MaxReplicas: 6,
	},
}

func TestScale(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://test/api/v1/query",
		httpmock.NewStringResponder(200, validPromResponse))

	commands := Scale(rules)

	assert.Equal(t, expectedCommands, commands)
}

func TestOutOfBoundsScale(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://test/api/v1/query",
		httpmock.NewStringResponder(200, validPromResponse))

	commands := Scale(constrainedRules)

	assert.Equal(t, expectedConstrainedCommands, commands)

}
