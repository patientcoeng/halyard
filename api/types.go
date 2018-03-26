package api

import "fmt"

// Map from resource name to annotations
type AnnotationMap map[string]map[string]string

type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type ASCommand struct {
	Resource    string
	Cmd         int32
	MinReplicas int32
	MaxReplicas int32
}

func (a ASCommand) String() string {
	return fmt.Sprintf("Avast me mateys! Scaling %s to %d replica(s) in range [%d, %d]",
		a.Resource, a.Cmd, a.MinReplicas, a.MaxReplicas)
}
