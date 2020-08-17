package main

import (
	"errors"

	yaml "github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RunningOnDifferentNodes struct {
}

func (c *RunningOnDifferentNodes) Check(data []byte) (string, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	if err != nil {
		return "", err
	}
	if deploy.Kind != "Deployment" {
		return "", nil
	}

	if deploy.Spec.Template.Spec.Affinity == nil || deploy.Spec.Template.Spec.Affinity.PodAntiAffinity == nil {

		return "For running on the spot instances, you'd better to distribute the pods among the different nodes. " +
			"'PodAntiAffinity' as the following could be used for that.\n" +
			`
            podAntiAffinity:
              requiredDuringSchedulingIgnoredDuringExecution:
              - labelSelector:
                matchExpressions:
                - key: app

                  operator: In
                  values:
                  - nginx
                  topologyKey: "kubernetes.io/hostname"
			`, nil

	}

	return "", nil
}

func (c *RunningOnDifferentNodes) Correct(org []byte) ([]byte, error) {
	hint, err := c.Check(org)
	if err != nil || hint == "" {
		return nil, errors.New("Can not correct/No need to correct.")
	}
	deploy := &v1.Deployment{}
	yaml.Unmarshal(org, deploy)
	app := deploy.Spec.Template.Labels["app"]
	matchLabels := make(map[string]string)
	matchLabels["app"] = app
	labelReq := metav1.LabelSelectorRequirement{
		Key:      "app",
		Operator: metav1.LabelSelectorOpIn,
		Values:   []string{app},
	}
	sel := &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{labelReq},
	}
	deploy.Spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				corev1.PodAffinityTerm{
					LabelSelector: sel,
					TopologyKey:   "kubernetes.io/hostname",
				},
			},
		},
	}

	corrected, err := yaml.Marshal(deploy)
	if err != nil {
		return nil, err
	}
	return corrected, nil
}
