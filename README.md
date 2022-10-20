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

# Knative Actors Client

The Knative Actors CLI enables the development and deployment of Actors using
Knative.

Typically, a developer defines a class, for instance using JavaScript:
```javascript
class Actor {
  set (v) { this.v = v; return 'OK' }
  get () { return this.v }
}
```
Once this actor code is deployed to Knative using the Actors CLI, it becomes
possible to invoke methods on instances of that class _serverlessly_. Actor
instances are created on the fly when first invoked and persisted in memory.
Unused actors are collected after a grace period.

## Limitations of current prototype

Knative Actors necessitate an experimental extension of Knative Serving under
development at https://github.com/tardieu/serving/tree/affinity. The [target
burst
capacity](https://knative.dev/docs/serving/load-balancing/target-burst-capacity/)
must be set to `-1` on Actor services to ensure the  activator is always in the
request path. This extension requires the deployment of a Redis instance.

For now a single Node.js actor runtime is implemented.

Knative Serving offers no durability guarantees.

The current implementation permits multiple method invocations of the same actor
instance to run concurrently.

## Documentation

- [Tutorial](docs/tutorial.md)
