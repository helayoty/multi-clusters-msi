# Authenticate and Authorize member cluster with hub cluster

### Steps:
- Create member cluster and hub cluster

```sh
  $ az aks create --resource-group demo --name hubCluster --node-count 1  --generate-ssh-keys --enable-aad --enable-azure-rbac

  $ az aks create --resource-group demo --name memberCluster --node-count 1  --generate-ssh-keys --enable-managed-identity
```

- Create a namespace for test purpose in hub cluster

```sh
 $ kubectl apply -f namespace.yaml
```

- Create Role and RoleBinding

```sh
 $ kubectl apply -f hub_roles.yaml
```

- Get the member cluster vmss managed identity' client id (*preq:* Install `jq`. An easy option is to install with `homebrew`)
  
  ```sh
  $ az aks show --resource-group demo --name memberCluster | jq .identityProfile.kubeletidentity.clientId
  ```
  Update var `MemberClusterClientId` in `main.go` with this value.

- Get the member cluster vmss managed identity' object id (*preq:* Install `jq`. An easy option is to install with `homebrew`)

  ```sh
  $ az aks show --resource-group demo --name memberCluster | jq .identityProfile.kubeletidentity.objectId
  ```
  Update field `MEMBER-CLUSTER-MANAGED-IDENTITY-OBJECT-ID>` in `hub_roles.yaml` with this value
- Get `hubCluster` server api address (an easy option is using Azure portal) and update var `HubServerApiAddress` in `main.go` with this value.
- Build your image

```sh
  docker build  . -t <docker_repo>/member-agent:v1
  docker push <docker_repo>/member-agent:v1
```

- Deploy the agent/app to member cluster

```sh
 $ kubectl apply -f member-agent-deployment.yaml
```

### Expectations:
  - The log should show that access to pods is forbidden for namespace `default'.
  - The log should show that there is 0 pod in `member-a` namespace and no error.
  - Try deploy `pod.yaml` in hub cluster. After successful deployment. The log should show that pod `demo-msi` is found in namespace `member-a`