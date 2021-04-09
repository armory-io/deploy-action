# deploy-action

_Warning: This action is in active development and subject to change._

## Description
`deploy-action` enables you to deploy to Kubernetes and AWS EC2 using Armory's Deployment API.

## Requirements
* A running Spinnaker deployment, accessible from your Github Actions runner.
* A `spin` CLI configuration for authentication. See [these instructions](https://spinnaker.io/setup/spin/#configure-spin) for how to set this up.

### Authentication

The `deploy-action` uses the same configuration as the `spin` CLI for authentication. Once you've configured it, we recommend
storing the full config file in a Github secret and injecting it into the action via the `config` property.

## Configuration

### Inputs

| Name            | Type     | Description                                                                        | Required | Default |
|-----------------|----------|------------------------------------------------------------------------------------|----------|---------|
| `application`   | String   | The name of the application you're deploying.                                      | Yes      |         |
| `accounts`      | []String | Comma-separated list of accounts you're deploying to.                              | Yes      |         |
| `cloudProvider` | String   | Name of the cloud provider. Values are `kubernetes`.                               | Yes      |         |
| `wait`          | Boolean  | Whether or not you'd like to wait for the deployment to complete.                  | No       | `false` |
| `baseUrl`       | String   | API URL to Gate.                                                                   | Yes      |         |
| `config`        | String   | Authentication configuration.                                                      | Yes      |         |
| `manifestPath`  | String   | If `cloudProvider` is `kubernetes`, path to where manifest files are stored.       | No       |         |
| `manifest`      | String   | If `cloudProvider` is `kubernetes`, this is the raw manifest you'd like to deploy. | No       |         |

## Examples - Kubernetes

### Deploy manifests that are stored in your repository

```yaml
name: Deploy with Armory
uses: armory-io/deploy-action@main
with:
  application: front50
  accounts: kubernetes-eks-prod
  manifestPath: deploy/manifests
  cloudProvider: kubernetes
  baseUrl: https://spinnaker-api.armory.io
  config: ${{ secrets.DEPLOY_CONFIG }}
  wait: true
```

### Use Kustomize to render manifests and then deploy them

```yaml
# Install Kustomize into the workspace
- name: Setup Kustomize
  uses: imranismail/setup-kustomize@v1
  with:
    kustomize-version: "3.1.0"

# Render manifests stored in deploy/environments/production and write them to a temporary directory
- name: Render manifests with Kustomize
  run: |
    mkdir -p build
    kustomize build deploy/environments/production > build/manifest.yml

# Deploy rendered manifests using the Armory deploy action
- name: Deploy with Armory
  uses: armory-io/deploy-action@main
  with:
    application: front50
    accounts: kubernetes-eks-prod
    manifestPath: build
    cloudProvider: kubernetes
    baseUrl: https://spinnaker-api.armory.io
    config: ${{ secrets.DEPLOY_CONFIG }}
    wait: true
```

## TODO
- [ ] Surface errors from the API instead of just saying there was a failure.