# Development


## Debug a service
1. Ensure weave cluster is setup correctly and you can deploy services (see [this guide](https://lab.weave.nl/nid/nid-example/-/blob/master/weave-cluster.md))
1. Ensure you have the bleeding edge skaffold version installed https://skaffold.dev/docs/install/
1. Run `make debug_service CLUSTER_HOST="..." SERVICE=...` (similar to `deploy_service` and `dev_service`)
1. Wait for the deployments to stabilizing
1. Connect the debugger of your favorite IDE to port `56268` which is a delve server (more information [here](https://skaffold.dev/docs/workflows/debug/))
1. Have fun debugging like you would for a binary running your local system

### Debugging with intellij
1. Go to run -> edit configurations
1. Add a new remote de configuration with localhost as host and `56268` as port


### Debugging with VSCode

To debug with vscode add the following configuration to your launch.json:
```
{
"name": "Skaffold Debug",
"type": "go",
"request": "attach",
"mode": "remote",
"remotePath": "${workspaceFolder}",
"port": 56268,
"host": "127.0.0.1"
}
```

### Running integration tests
See `integration-tests/README.md`.
