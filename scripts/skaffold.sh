#!/usr/bin/env bash
set -e # Fail on error

[ -z "$DEPLOYMENT_PATH" ] && DEPLOYMENT_PATH="../nid-example"

cluster=$(kubectl config current-context)
if  [[ ! $cluster =~ ^.*\microk8s ]]; then
    echo "It seems like you're not connected to a weave dev cluster: $cluster";
    exit 1;
fi


if [[ $(istioctl version -o json | jq .meshVersion) == "" ]]; then
  echo "It seems like istio is not installed yet";
  exit 1;
fi

if [ -z "$CLUSTER_HOST" ]; then
  CLUSTER_HOST=$(cd "$DEPLOYMENT_PATH"; weave-cluster dns)
  export CLUSTER_HOST
fi

if [ -z "$2" ]; then
  echo "Please specify \$2 corresponding to the service name"
  exit 1
fi

export SERVICE=$2

if [ -z "$SERVICE_DIRECTORY" ]; then
  export SERVICE_DIRECTORY=$SERVICE
else
  export SERVICE_DIRECTORY
fi

if [[ -f "svc/$SERVICE_DIRECTORY/proto/api_descriptor.pb" ]]; then
  echo 'API DESCRIPTOR FOUND!'
  API_DESCRIPTOR=$(base64 -i svc/$SERVICE_DIRECTORY/proto/api_descriptor.pb)
  export API_DESCRIPTOR
fi

# The callback function when running debug requires this arg otherwise it will not know which file to use
SKAFFOLD_FILENAME="$DEPLOYMENT_PATH/skaffold.yaml"

export SKAFFOLD_FILENAME

export ISTIO_VERSION="$(istioctl version -o json |  grep -v "no running Istio" | jq .clientVersion.version | sed -e 's/^"//' -e 's/"$//')"
echo "istio version: $ISTIO_VERSION"

# FIXME for now we use a grep to filter the istio logs but there are several skaffold proposals to fix this https://github.com/GoogleContainerTools/skaffold/issues/588
# this does break killing the skaffold cleanup logs

if [ "$1" == "debug" ]; then
  echo "Warning to run debug you need to have the latest master version of skaffold installed"
  additional_args="--port-forward"
  export SERVICE_TYPE=""
else
  SERVICE_TYPE="$(yq e ".type" $DEPLOYMENT_PATH/services/$SERVICE.yaml)";
  export SERVICE_TYPE
fi


skaffold "$1" -f $DEPLOYMENT_PATH/skaffold.yaml --cache-artifacts=false --verbosity=info $port_forward | grep -v "^\[istio-proxy\]\|istio-init"  #--default-repo registry.weave.nl/devops/weave-cluster/nid-$USER-$2
