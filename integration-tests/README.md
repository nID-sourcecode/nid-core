# Integration tests

## Prerequisites
- Make sure your local `nid-example` and `nid-services` are up to date.
- In your `nid-example` dir, make sure your `env.yaml` matches `env_sample.yaml` and has all values populated.

## Short instructions
To run the integration tests:
```bash
cp env/e2e-setup.sample.env env/e2e-setup.local.env # change the values if needed
make setup_e2e_cluster
make test.e2e
```

## Detailed instructions
To run the integration tests:
### Setup env 
Copy env/e2e_setup.sample.env to env/e2e_setup.env, change the values if needed.

| Variable             | Description                            | Default                      |
| -------------------- | -------------------------------------- | ---------------------------- |
| SETUP_CLUSTER_NAME   | The name the new cluster will be given | `e2e[random 8-char string]` |
| NID_EXAMPLE_PATH     | The local path to nid-example          |                              |
| NID_SERVICES_PATH    | The local path to nid-services         |                              |

### Setup cluster
Run `make setup_e2e_cluster`.
This spins up a new cluster, installs istio, and installs all services needed to run the e2e tests.

It will also create two new environment files:

`env/e2e-cluster.generated.env`:
```.env
export CLUSTER_HOST=e2e92661b5d.wouter.dev.weave.nl
export CLUSTER_NAME=e2e92661b5d
export CLUSTER_STATE=5
```
Which represents the state of your cluster and is used by `make setup_e2e_cluster` to continue if something went wrong,
and to check whether you already have a cluster running.

`env/e2e-execution.generated.env`:
```.env
BACKEND_PORT='80'
IS_TLS='false'
NAMESPACE='nid'
CLIENT_PASSWORD='8]B-;wfb\7R3sJ2p'
DEVICE_CODE='testingdevice'
DEVICE_SECRET='y\wU;3UyH89P_a_9'
DASHBOARD_USER_EMAIL='wouter@weave.nl'
DASHBOARD_USER_PASS='4485ub9XQBiJ@kdf'
BACKEND_URL='e2e92661b5d.wouter.dev.weave.nl'
```
This information is pulled from `$NID_EXAMPLE_PATH/env.yaml` and is used for the test execution.

### Run the tests

Run `make test.e2e`. This will source `env/e2e-execution.generated.env` and run the integration tests. 
