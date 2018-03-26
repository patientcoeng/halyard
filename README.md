# Halyard
[![Build Status](https://travis-ci.org/patientcoeng/halyard.svg?branch=master)](https://travis-ci.org/patientcoeng/halyard)

Halyard is a tool for horizontal autoscaling with Kubernetes (K8S). It leverages a metrics backend to determine the current value of a metric, and linearly scales the number of pods so the metric approaches a given target value. This enables dynamic scaling in response to changes in load.

## Backends
Currently supported backends:

* [Prometheus](https://prometheus.io/)

## Running in Docker
1. docker run -d -v Config.json:/config/Config.json:ro patientcoeng/halyard

## Running in Kubernetes
1. curl -o deployment.yaml https://raw.githubusercontent.com/patientcoeng/halyard/master/manifests/deployment.yaml
2. Edit deployment.yaml to replace the placeholder values with those specific to your setup
3. kubectl apply -f deployment.yaml

## Building from Source
1. Clone the halyard repository
2. Install govendor "go install github.com/kardianos/govendor"
3. Run "govendor sync"
2. Run "go build"
3. Run halyard

## Configuration
Halyard expects a JSON-encoded configuration file to be present in the same directory as the Halyard binary. You can use the config command line option to use a different file. An example file is listed below, for reference:

    {
      "period": 30,
      "namespace": "default",
      "backendEndpoints": {
        "prometheus": "https://my-prometheus-server"
      },
      "alertConfig": {
          "slack": {
            "webhookURL": "https://hooks.slack.com/services/YOUR_SERVICE_ID",
            "channel": "#yourchannel"
          }
      }
    }

## Deployment Annotations
Halyard works by reading annotations on Kubernetes. The following annotations are supported. All values are strings.

* **halyard.patientco.com/query**
  * A query string, recognized by the designated backend
* **halyard.patientco.com/backend**
  * The name of the backend to process the given query
* **halyard.patientco.com/target-value**
  * The target value to set the query result to, as a string.
* **halyard.patientco.com/min-replicas**
  * The minimum number of allowed replicas
* **halyard.patientco.com/max-replicas**
  * The maximum number of allowed replicas
* **halyard.patientco.com/scaling-policy**
  * Currently only "linear" is supported

## CLI Options
    -config string
        Location of JSON-encoded config file (default "Config.json")