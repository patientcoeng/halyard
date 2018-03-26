package autoscaler

import (
	"github.com/patientcoeng/halyard/alerting"
	"github.com/patientcoeng/halyard/api"
	"github.com/patientcoeng/halyard/backend"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
)

type ASRule struct {
	Query         string
	Backend       string
	Resource      string
	Result        float64
	Target        float64
	MinReplicas   int32
	MaxReplicas   int32
	ScalingPolicy string
	Endpoint      string
	AlertManager  *alerting.Manager
}

type BackendAdapter interface {
	Query(query string) float64
}

const (
	QueryAnnotationName   = "halyard.patientco.com/query"
	BackendAnnotationName = "halyard.patientco.com/backend"
	TargetAnnotationName  = "halyard.patientco.com/target-value"
	MinAnnotationName     = "halyard.patientco.com/min-replicas"
	MaxAnnotationName     = "halyard.patientco.com/max-replicas"
	ScalingAnnotationName = "halyard.patientco.com/scaling-policy"
)

func CreateASRules(annotationMap api.AnnotationMap, backendMap api.EndpointMap, manager *alerting.Manager) []ASRule {
	var result []ASRule

	for resource, annotations := range annotationMap {
		var minValue int32 = -1
		var maxValue int32 = -1
		var scalingPolicy string = "linear"

		queryAnnotation, ok := annotations[QueryAnnotationName]
		if !ok {
			log.Debug().Msgf("Skipping resource: %s", resource)
			continue
		}

		backendAnnotation, ok := annotations[BackendAnnotationName]
		if !ok {
			log.Debug().Msgf("Skipping resource: %s", resource)
			continue
		}

		targetAnnotation, ok := annotations[TargetAnnotationName]
		if !ok {
			log.Debug().Msgf("Skipping resource: %s", resource)
			continue
		}

		minAnnotation, ok := annotations[MinAnnotationName]
		if ok {
			minValue64, err := strconv.ParseInt(minAnnotation, 10, 32)
			if err != nil {
				minValue = -1
			} else {
				minValue = int32(minValue64)
			}
		}

		maxAnnotation, ok := annotations[MaxAnnotationName]
		if ok {
			maxValue64, err := strconv.ParseInt(maxAnnotation, 10, 32)
			if err != nil {
				maxValue = -1
			} else {
				maxValue = int32(maxValue64)
			}
		}

		scalingAnnotation, ok := annotations[ScalingAnnotationName]
		if ok {
			scalingPolicy = scalingAnnotation
		}

		targetValue, err := strconv.ParseFloat(targetAnnotation, 64)
		if err != nil {
			log.Error().Msgf("Error converting %s to a float.", targetAnnotation)
			continue
		}

		backendEndpoint, ok := backendMap[backendAnnotation]
		if !ok {
			log.Error().Msgf("Unable to find endpoint for %s", backendAnnotation)
		}

		if maxValue < 0 {
			maxValue = math.MaxInt32
		}

		if minValue < 0 {
			minValue = 0
		}

		log.Debug().Msgf("Creating rule for %s", resource)

		rule := ASRule{
			Query:         queryAnnotation,
			Backend:       backendAnnotation,
			Resource:      resource,
			Target:        targetValue,
			ScalingPolicy: scalingPolicy,
			MinReplicas:   minValue,
			MaxReplicas:   maxValue,
			Endpoint:      backendEndpoint,
			AlertManager:  manager,
		}

		result = append(result, rule)
	}

	return result
}

func (a *ASRule) Evaluate() {
	var backendInterface BackendAdapter

	if a.Backend == "prometheus" {
		backendInterface = &backend.PrometheusBackend{
			Endpoint:     a.Endpoint,
			AlertManager: a.AlertManager,
		}
	} else {
		log.Debug().Msgf("Invalid or unsupported backend")
		a.Result = a.Target
		return
	}

	result := backendInterface.Query(a.Query)

	if math.IsNaN(result) {
		a.Result = a.Target
	} else {
		a.Result = result
	}
}
