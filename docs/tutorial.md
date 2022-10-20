<!--
# Copyright IBM Corporation 2022
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
-->

# Tutorial

In this tutorial, we develop, deploy, and invoke a simple Knative actor based on
the actor runtime for Node.js.

## Prerequisites

First install [Go](https://go.dev), [Docker](https://www.docker.com),
[Kind](https://kind.sigs.k8s.io), the [Knative
CLI](https://knative.dev/docs/client/install-kn/), and the [Quickstart
plugin](https://knative.dev/docs/getting-started/quickstart-install/) for the
Knative CLI.

## Install the CLI

Install the Knative Actors CLI from source:
```bash
git clone https://github.com/tardieu/kn-actor.git
cd kn-actor
go install ./...
```

## Setup the Knative cluster

Create a Kind cluster and deploy Knative Serving.
```bash
kn quickstart kind --install-serving
```

Deploy Redis.
```bash
kubectl apply -n knative-serving -f redis.yaml
```

Patch the autoscaler configuration to keep the activator in the request path.
```bash
kubectl patch configmap/config-autoscaler -n knative-serving -p '{"data":{"target-burst-capacity": "-1"}}'
```

Override the default activator image to add support for session and revision
affinity.
```bash
kubectl set image -n knative-serving deployment activator activator=quay.io/tardieu/activator:dev
```

Optionally replicate the activator to test consistency.
```bash
kubectl patch hpa activator -n knative-serving -p '{"spec":{"minReplicas":2,"maxReplicas":20}}'
```

## Create the example actor

Create an actor project using the Node.js actor template.
```bash
mkdir sample
cd sample
kn-actor create --runtime node
```

The template actor is declared in the `index.js` file. It includes three example
methods.
```javascript
class Actor {
  set (v) { this.v = v; return 'OK' }
  get () { return this.v }
  ip () { return require('os').networkInterfaces().eth0[0].address }
}

// DO NOT MODIFY CODE BELOW THIS POINT
// ...
```
The `set` method stores a JSON value in the actor instance. The `get` method
retrieves the stored value if any. The `ip` method returns the ip of the
Kubernetes pod running the actor instance for debugging purposes. Edit the
`index.js` file to alter the behavior of the actor. Methods may be `async`.


## Deploy the example actor

Build and publish the container image for the actor project.
```bash
kn-actor build --image kind.local/sample:dev --push
```

Deploy the actor service to the Knative cluster for instance scaling from 3 to 7
replicas.
```bash
kn service create sample --image kind.local/sample:dev --scale 3..7
```
For now the minimal scale must be no less than `1` and the target burst capacity
should not be changed from the autoscaler `-1` default value due to limitations
of the session affinity implementation.

## Invoke the example actor

Actor instances may be invoked using the CLI.
```bash
kn-actor invoke --service sample --instance instance1 --method set --data 42
```
```bash
kn-actor invoke --service sample --instance instance2 --method set --data '"hello"'
```
```bash
kn-actor invoke --service sample --instance instance1 --method get
```
```bash
kn-actor invoke --service sample --instance instance2 --method get
```
```bash
kn-actor invoke --service sample --instance instance1 --method ip
```
```bash
kn-actor invoke --service sample --instance instance2 --method ip
```
Actor instances are created on the fly when first invoked. The `invoke` command
requires the actor service name (`--service` flag), the instance name
(`--instance` flag), and the method name (`--method` flag). Optional arguments
to the method invocation are specified using `--data` flags. Arguments are
expected to be valid JSON values. In particular, strings have to be enclosed in
quotes. An invocation may return nothing, one JSON value, or one error.

Actor instances may also be invoked from other actor instances. For Node.js
actors, the syntax is:
```JavaScript
await actor.invoke(service, instance, method, ...args)
```

## Cleanup

To cleanup simply delete the Kind cluster:

```bash
kind delete cluster --name knative
```
