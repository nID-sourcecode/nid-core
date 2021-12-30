#!/usr/bin/env bash
set -e # Fail on error

export NID_CORE_PATH=$PWD

if [[ ! -f "env/e2e-setup.local.env" ]]; then
  echo "could not find env/e2e-setup.local.env. Please copy this file from env/e2e-setup.sample.env (note that .local was added recently)"
  exit 1
fi

source env/e2e-setup.local.env

if [ -z "$NID_SERVICES_PATH" ]; then
  echo "NID_SERVICES_PATH not set, set it in env/e2e-setup.local.env"
  exit 1
fi

if [ -z "$NID_EXAMPLE_PATH" ]; then
  echo "NID_EXAMPLE_PATH not set, set it in env/e2e-setup.local.env"
  exit 1
fi

export DEPLOYMENT_PATH="$NID_EXAMPLE_PATH"

# Check whether creating a cluster is needed
skip_cluster_creation=false
if [ -f "env/e2e-cluster.generated.env" ]; then
  source env/e2e-cluster.generated.env

  cluster=$(kubectl config current-context)
  echo "DOING CHECK FOR $cluster and $CLUSTER_NAME"
  if  [[ $cluster =~ ^.*microk8s.*$CLUSTER_NAME ]]; then
      if timeout 3s kubectl cluster-info; then
        echo "It seems like you're already connected to an e2e cluster: $CLUSTER_NAME. Do you want to use this cluster? (no creates a new cluster instead) [y/n]"
        read
        if [[ $REPLY == "y" ]]; then
          skip_cluster_creation=true
        fi
      fi
  fi
fi;

if [[ $skip_cluster_creation == "false" ]]; then
  # Create a cluster
  if [ -n "$SETUP_CLUSTER_NAME" ]; then
    CLUSTER_NAME="$SETUP_CLUSTER_NAME"
  else
    CLUSTER_UUID="$(uuidgen | tr "[:upper:]" "[:lower:]")"
    CLUSTER_NAME="e2e${CLUSTER_UUID:0:8}" # Generate randomly
  fi

  cd $NID_EXAMPLE_PATH

  echo "creating cluster with name $CLUSTER_NAME"
  weave-cluster create $CLUSTER_NAME

  export CLUSTER_HOST=$(weave-cluster dns)

  cd $NID_CORE_PATH
  echo "export CLUSTER_HOST=$CLUSTER_HOST
  export CLUSTER_NAME=$CLUSTER_NAME
  export CLUSTER_STATE=1" > env/e2e-cluster.generated.env

  export CLUSTER_STATE=1
fi

# Set up execution env (parse from env.yaml)
if [[ $CLUSTER_STATE == 1 ]]; then
  echo "BACKEND_PORT='80'
IS_TLS='false'
NAMESPACE='nid'" > env/e2e-execution.generated.env

  cd $NID_EXAMPLE_PATH

  function unquote(){
    sed -e 's/^"//' -e 's/"$//'
  }

  CLIENT_PASSWORD=$(yq e '.secrets.[] | select(.name == "testing-client-password") | .data.password' env.yaml)

  # FIXME this shows a problem in the way the wallet testing entry is seeded; it should be less flexible and only contain relevant information
  wallet_user_json=$(yq e '.secrets.[] | select(.name == "default-wallet-users") | .data."users.json"' env.yaml | jq '.[] | select(.pseudonym=="Ajmjkuq6JKiCWEevkB1V7SPNCRd8uHAKh0nABF7BSTXZV9a7k2eZ2iE2BxBdVI60")')
  DEVICE_CODE=$(jq -r '.devices[0].code' <<< $wallet_user_json)
  DEVICE_SECRET=$(jq -r '.devices[0].secret' <<< $wallet_user_json)

  DASHBOARD_USER_EMAIL=$(yq e '.secrets.[] | select(.name == "defaultuser") | .data.email' env.yaml)
  DASHBOARD_USER_PASS=$(yq e '.secrets.[] | select(.name == "defaultuser") | .data.pass' env.yaml)

  cd $NID_CORE_PATH
  echo "CLIENT_PASSWORD='$CLIENT_PASSWORD'
DEVICE_CODE='$DEVICE_CODE'
DEVICE_SECRET='$DEVICE_SECRET'
DASHBOARD_USER_EMAIL='$DASHBOARD_USER_EMAIL'
DASHBOARD_USER_PASS='$DASHBOARD_USER_PASS'
BACKEND_URL='$CLUSTER_HOST'" >> env/e2e-execution.generated.env

  export CLUSTER_STATE=2
  sed -i 's/CLUSTER_STATE=1/CLUSTER_STATE=2/g' env/e2e-cluster.generated.env
fi

# Install istio
if [[ $CLUSTER_STATE == 2 ]]; then
  cd $NID_EXAMPLE_PATH
  [ -z "$MONITORING" ] && export MONITORING=false
  echo "installing istio, monitoring set to $MONITORING"
  ./scripts/installistio_dev.sh $MONITORING true

  cd $NID_CORE_PATH
  export CLUSTER_STATE=3
  sed -i 's/CLUSTER_STATE=2/CLUSTER_STATE=3/g' env/e2e-cluster.generated.env
fi

# Install core services
if [[ $CLUSTER_STATE == 3 ]]; then
  cd $NID_CORE_PATH
  [ -z "$SERVICES_TO_DEPLOY" ] && SERVICES_TO_DEPLOY="info,info-manager,info-manager-gql,pseudonymization,autopseudo,scopeverification,auditlog,auth,auth-gql,wallet-rpc,dashboard,documentation,autobsn"
  make deploy_services SERVICES=$SERVICES_TO_DEPLOY

  export CLUSTER_STATE=4
  sed -i 's/CLUSTER_STATE=3/CLUSTER_STATE=4/g' env/e2e-cluster.generated.env
fi

# Install nid-services services
if [[ $CLUSTER_STATE == 4 ]]; then
  cd $NID_SERVICES_PATH
  make deploy_service SERVICE=databron
  make deploy_service SERVICE=other-databron SERVICE_DIRECTORY=databron
  make deploy_service SERVICE=information

  cd $NID_CORE_PATH
  export CLUSTER_STATE=5
  sed -i 's/CLUSTER_STATE=4/CLUSTER_STATE=5/g' env/e2e-cluster.generated.env
fi
