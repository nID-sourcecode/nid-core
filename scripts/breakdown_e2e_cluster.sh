#!/usr/bin/env bash
set -e # Fail on error

export NID_CORE_PATH=$PWD

source env/e2e-setup.local.env
if [ -f env/e2e-cluster.generated.env ]; then
    source env/e2e-cluster.generated.env
fi

if [ -z "$CLUSTER_NAME" ]; then
  echo "CLUSTER_NAME not set, nothing to break down"
  exit 1
fi

if [ -z "$NID_EXAMPLE_PATH" ]; then
  echo "NID_EXAMPLE_PATH not set, set it in env/e2e-setup.local.env"
  exit 1
fi

cd $NID_EXAMPLE_PATH
weave-cluster delete $CLUSTER_NAME

cd $NID_CORE_PATH
rm env/e2e-cluster.generated.env
