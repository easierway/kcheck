# kcheck

## Kubernetes Deployment Configuration Checking Tool

### Introduction
In CI/CD process, we can leverage some checking tools such as, Checkstyle, FindBugs to enforce the best coding practices in the integration phrase. KCheck is the checking tool to enforce the best practices of kubernetes configuration in the deployment phrase. 
Before jumping into the project, let’s list some best practices for your daily Kubernetes deployments.

### Best practice for the Kubernetes deployment
First of the all, a good deployment on Kubernetes, especially running in cloud, should always makes your service prepare well for losing nodes.
In Kubernetes world, containers are the essential elements of a service/application. To improve your service/application’s availability, you should make sure the right things would be done in the each stage of the container lifecycle.

**1 Readiness Probe is required** 
Readiness Probe is to inform the Kubernetes when the pod is ready, and then the pod could be put behind the load balance. The kubelet uses readiness probes to know when a container is ready to start accepting traffic. A Pod is considered ready when all of its containers are ready.
Readiness probe is foundamental to provide an available service. Without readiness probe, the customer’s requests might be dispatched to the unready pods. Even, the very basic kubernetes functions, such as “zero downtime rolling update” and HPA, will be screwed up. For “rolling update”, the strategy settings “maxUnavailable” and “maxSurge” would be ineffective for treating the unready pod as ready.

**2 Liveness Probe is required**
Liveness is to make sure the pod is in the healthy state. Just like the readiness probe, liveness probe is also critical the availability of the service. The same as what I mentioned above about readiness probe, it relates to “rolling update” and HPA. Kubernetes uses liveness probes to know when to restart a container. Restarting a container in the unhealthy state can help to make the application more available despite bugs. It is a kind of recovery oriented solution.
(More details about readiness and liveness, please, see also: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)

**3 Termination needs to be handled carefully**
Terminating the process gracefully is required in many cases. The graceful termination would impact the data’s integrity and consistency. In kubernetes, this can be done by container lifecycle hook. By leveraging “PreStop” hook, your termination process is invoked immediately before a container is terminated due to an API request or management event such as liveness probe failure, preemption, resource contention and others.

**4 Always declare resources requests and limits**
This prevents the pod from being starved of resources while also preventing CPU/Mem from consuming all of the resources on a node. The negative impacts caused by missing resources declarations is not only on the pod itself, but on the whole kubernetes cluster. Without the requests, the Kubernetes scheduler cannot ensure that workloads are spread across your nodes evenly and this may result in an unbalanced cluster with some nodes overcommitted and some nodes underutilized. Also, The resource declarations define the QoS of Pod, for more detail refer to https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/

**5 Control the pods’ distribution**
In some cases, it is better to make our pods be scheduled on the different nodes. This can minimize the impact on the service capacity which is caused by the unexpected terminations of the nodes. For example, we run our services on AWS spot instances. 
In kubernetes, we can use pod affinity to instruct kubernetes to schedule the pods according to our strategy.

There are still other best practices or deployment instructions for your organization. Like coding style, normally, we want to make sure everyone in our organization follow these instructions. So, an automation check is required in CI/CD process. 
This project is the implementation of an extensible kubernetes deployment configuration check tool.



### User Guide

**1 Define the checking rules**
You can define checking rules in a YAML file. Each rule is composed of several check items. Also, you can define the correctors which can correct the configurations according to the check item (Not all check items support auto-correcting, only the one whose implementation realizes the Corrector interface)


**2 Run kcheck**
Command line parameters: 
 -d string
  	the rule definition file

 -f string
  	the kubernetes deployment/configuration file

 -help
  	get the help

 -r string
  	the name of the checking rule

 -c	
​    [Optional] try to correct the files 

Example:
The rule definition 
my_rules.yaml
``` YAML
rules:
- name: spot
  checkItems:
  - RunningOnDifferentNodes
  - WithGracefulTermination
  - WithHealthCheck
  - WithResourceRequestAndLimit
- name: normal
  checkItems:
  - WithHealthCheck
  - WithResourceRequestAndLimit
  - WithReadiness
correctors:
- RunningOnDifferentNodes
```

The deployment configuration 
example_deployment.yaml
``` YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

Run the following command line: 

```
kcheck -d my_rules.yaml -f example_deployment_1.yaml -r normal 
```

For the above example, the following contents will be shown on the console,

```
'LivenessProbe' should be set for container: nginx.
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
      periodSeconds: 5

'Resource requests & limits' should be set for container: nginx.
resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"

It is nice to have 'ReadinessProbe' setting for container: nginx.
spec:
  containers:
    readinessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10 

```

By using -c parameter, if the broken CheckItem being has implements Corrector interface, the configuration would be corrected automatically. The corrected configuration is stored in the file "coorected.yaml"