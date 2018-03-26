package k8s

import (
	"github.com/patientcoeng/halyard/api"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

type K8S struct {
	Clientset *kubernetes.Clientset
	Namespace string
}

func NewK8S(namespace string) *K8S {
	// Set up K8S Interface
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error().Msgf("Error creating config: %s", err)
		return nil
	}
	config.Timeout = 30 * time.Second

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Msgf("Error creating clientset: %s", err)
		return nil
	}

	return &K8S{
		Clientset: clientset,
		Namespace: namespace,
	}
}

func (k *K8S) GetDeployAnnotations() (api.AnnotationMap, error) {
	deployments, err := k.Clientset.AppsV1beta1().Deployments(k.Namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make(api.AnnotationMap)
	for _, item := range deployments.Items {
		result[item.Name] = item.Annotations
	}

	return result, nil
}

func (k *K8S) UpdateReplicas(commands []api.ASCommand) error {
	deploymentsClient := k.Clientset.AppsV1beta1().Deployments(k.Namespace)
	deployments, err := deploymentsClient.List(v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, item := range deployments.Items {
		for _, cmd := range commands {
			if cmd.Resource == item.Name {
				newReplicas := cmd.Cmd
				item.Spec.Replicas = &newReplicas

				_, err = deploymentsClient.Update(&item)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
