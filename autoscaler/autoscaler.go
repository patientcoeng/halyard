package autoscaler

import (
	"github.com/patientcoeng/halyard/api"
	"math"
)

func Scale(rules []ASRule) []api.ASCommand {
	var result []api.ASCommand

	for i, _ := range rules {
		rules[i].Evaluate()
	}

	for _, rule := range rules {
		if rule.ScalingPolicy == "linear" {
			ratio := int32(math.Ceil(rule.Result / rule.Target))
			if ratio > rule.MaxReplicas {
				ratio = rule.MaxReplicas
			} else if ratio < rule.MinReplicas {
				ratio = rule.MinReplicas
			}

			result = append(result, api.ASCommand{
				Resource:    rule.Resource,
				Cmd:         ratio,
				MinReplicas: rule.MinReplicas,
				MaxReplicas: rule.MaxReplicas,
			})
		}
	}

	return result
}
