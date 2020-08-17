// normal_check.go
package main

import (
	yaml "github.com/ghodss/yaml"
	v1 "k8s.io/api/apps/v1"
)

type WithGracefulTermination struct {
}

func (c *WithGracefulTermination) Check(data []byte) (string, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	if err != nil {
		return "", err
	}
	if deploy.Kind != "Deployment" {
		return "", nil
	}
	hints := ""
	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].Lifecycle == nil ||
				deploy.Spec.Template.Spec.Containers[i].Lifecycle.PreStop == nil {

				hints = "'preStop' should be set for a graceful termination [container: " +
					deploy.Spec.Template.Spec.Containers[i].Name + "]." +

					`
spec:
  containers:
  - name: lifecycle-demo-container
    image: nginx
    lifecycle:
      preStop:
        exec:
          command: ["/bin/sh","-c","nginx -s quit"]
			` + "\n"

			}

		}
	}
	if hints == "" {
		return "", nil
	}
	return hints, nil
}

type WithHealthCheck struct {
}

func (c *WithHealthCheck) Check(data []byte) (string, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	if err != nil {
		return "", err
	}
	if deploy.Kind != "Deployment" {
		return "", nil
	}
	hints := ""
	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].LivenessProbe == nil {

				hints = "'LivenessProbe' should be set for container: " +
					deploy.Spec.Template.Spec.Containers[i].Name + "." +

					`
spec:
  containers:
  - name: lifecycle-demo-container
    image: nginx
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5` + "\n"

			}

		}
	}
	if hints == "" {
		return "", nil
	}
	return hints, nil
}

type WithReadiness struct {
}

func (c *WithReadiness) Check(data []byte) (string, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	if err != nil {
		return "", err
	}
	if deploy.Kind != "Deployment" {
		return "", nil
	}
	hints := ""
	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].ReadinessProbe == nil {

				hints = "It is nice to have 'ReadinessProbe' setting for container: " +
					deploy.Spec.Template.Spec.Containers[i].Name + "." +

					`
spec:
  containers:
    readinessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10 `

			}

		}
	}
	if hints == "" {
		return "", nil
	}
	return hints, nil
}

type WithResourceRequestAndLimit struct {
}

func (c *WithResourceRequestAndLimit) Check(data []byte) (string, error) {
	deploy := &v1.Deployment{}
	err := yaml.Unmarshal(data, deploy)
	if err != nil {
		return "", err
	}
	if deploy.Kind != "Deployment" {
		return "", nil
	}
	hints := ""
	if deploy.Spec.Template.Spec.Containers != nil &&
		len(deploy.Spec.Template.Spec.Containers) > 0 {
		for i := 0; i < len(deploy.Spec.Template.Spec.Containers); i++ {
			if deploy.Spec.Template.Spec.Containers[i].Resources.Requests == nil ||
				deploy.Spec.Template.Spec.Containers[i].Resources.Limits == nil {

				hints = "'Resource requests & limits' should be set for container: " +
					deploy.Spec.Template.Spec.Containers[i].Name + "." +

					`
resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"` + "\n"

			}

		}
	}
	if hints == "" {
		return "", nil
	}
	return hints, nil
}
