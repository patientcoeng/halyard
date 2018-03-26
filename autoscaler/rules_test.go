package autoscaler

import (
	"github.com/patientcoeng/halyard/api"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var annotationMap = api.AnnotationMap{
	"testPod": map[string]string{
		QueryAnnotationName:   "testquery",
		BackendAnnotationName: "prometheus",
		TargetAnnotationName:  "1.0",
		ScalingAnnotationName: "linear",
	},
}

var doubleAnnotationMap = api.AnnotationMap{
	"testPod": map[string]string{
		QueryAnnotationName:   "testquery",
		BackendAnnotationName: "prometheus",
		TargetAnnotationName:  "1.0",
		ScalingAnnotationName: "linear",
		MinAnnotationName:     "3",
		MaxAnnotationName:     "15",
	},
	"testPod2": map[string]string{
		QueryAnnotationName:   "testquery",
		BackendAnnotationName: "prometheus",
		TargetAnnotationName:  "1.0",
		ScalingAnnotationName: "linear",
	},
}

var backendMap = api.EndpointMap{
	"prometheus": "http://test/",
}

var expectedRules = []ASRule{
	ASRule{
		Query:         "testquery",
		Backend:       "prometheus",
		Resource:      "testPod",
		Target:        1,
		ScalingPolicy: "linear",
		MinReplicas:   0,
		MaxReplicas:   math.MaxInt32,
		Endpoint:      "http://test/",
		AlertManager:  manager,
	},
}

var doubleExpectedRules = []ASRule{
	ASRule{
		Query:         "testquery",
		Backend:       "prometheus",
		Resource:      "testPod",
		Target:        1,
		ScalingPolicy: "linear",
		MinReplicas:   3,
		MaxReplicas:   15,
		Endpoint:      "http://test/",
		AlertManager:  manager,
	},
	ASRule{
		Query:         "testquery",
		Backend:       "prometheus",
		Resource:      "testPod2",
		Target:        1,
		ScalingPolicy: "linear",
		MinReplicas:   0,
		MaxReplicas:   math.MaxInt32,
		Endpoint:      "http://test/",
		AlertManager:  manager,
	},
}

func TestRules(t *testing.T) {
	rules := CreateASRules(annotationMap, backendMap, manager)
	assert.True(t, equalSlices(expectedRules, rules))
}

func TestMinMaxLeak(t *testing.T) {
	rules := CreateASRules(doubleAnnotationMap, backendMap, manager)
	assert.True(t, equalSlices(doubleExpectedRules, rules))
}

func TestMinMax(t *testing.T) {
	var newAnnotationMap = make(api.AnnotationMap)
	for k, v := range annotationMap {
		newAnnotationMap[k] = v
	}
	newAnnotationMap["testPod"][MinAnnotationName] = "6"
	newAnnotationMap["testPod"][MaxAnnotationName] = "100"

	var newExpectedRules []ASRule
	for _, rule := range expectedRules {
		newExpectedRules = append(newExpectedRules, rule)
	}
	newExpectedRules[0].MinReplicas = 6
	newExpectedRules[0].MaxReplicas = 100

	rules := CreateASRules(newAnnotationMap, backendMap, manager)
	assert.Equal(t, newExpectedRules, rules)
}

func TestRulesFatalAnnotationQuery(t *testing.T) {
	var emptyRules []ASRule
	var missingAnnotationMap = api.AnnotationMap{
		"testPod": map[string]string{
			BackendAnnotationName: "prometheus",
			TargetAnnotationName:  "1.0",
		},
	}

	rules := CreateASRules(missingAnnotationMap, backendMap, manager)
	assert.Equal(t, emptyRules, rules)
}

func TestRulesFatalAnnotationBackend(t *testing.T) {
	var emptyRules []ASRule
	var missingAnnotationMap = api.AnnotationMap{
		"testPod": map[string]string{
			QueryAnnotationName:  "testquery",
			TargetAnnotationName: "1.0",
		},
	}

	rules := CreateASRules(missingAnnotationMap, backendMap, manager)
	assert.Equal(t, emptyRules, rules)
}

func TestRulesFatalAnnotationTarget(t *testing.T) {
	var emptyRules []ASRule
	var missingAnnotationMap = api.AnnotationMap{
		"testPod": map[string]string{
			QueryAnnotationName:   "testquery",
			BackendAnnotationName: "prometheus",
		},
	}

	rules := CreateASRules(missingAnnotationMap, backendMap, manager)
	assert.Equal(t, emptyRules, rules)
}

func TestRulesNonFatalAnnotations(t *testing.T) {
	var newAnnotationMap = make(api.AnnotationMap)
	for k, v := range annotationMap {
		newAnnotationMap[k] = v
	}
	newAnnotationMap["testPod"][MaxAnnotationName] = "abcdef"
	newAnnotationMap["testPod"][MinAnnotationName] = "ghijk"

	rules := CreateASRules(newAnnotationMap, backendMap, manager)
	assert.Equal(t, expectedRules, rules)
}

func TestBadBackend(t *testing.T) {
	var newAnnotationMap = make(api.AnnotationMap)
	for k, v := range annotationMap {
		newAnnotationMap[k] = v
	}
	newAnnotationMap["testPod"][BackendAnnotationName] = "notPrometheus"
	newExpectedRules := expectedRules
	newExpectedRules[0].Backend = "notPrometheus"
	newExpectedRules[0].Endpoint = ""

	rules := CreateASRules(newAnnotationMap, backendMap, manager)
	assert.Equal(t, expectedRules, rules)
	assert.Equal(t, 1, len(rules))
	rules[0].Evaluate()
	assert.Equal(t, rules[0].Target, rules[0].Result)
}

func equalSlices(a []ASRule, b []ASRule) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	allFound := true
	for _, aEle := range a {
		oneFound := false
		for _, bEle := range b {
			if aEle == bEle {
				oneFound = true
				break
			}
		}
		if !oneFound {
			allFound = false
			break
		}
	}

	return allFound
}
