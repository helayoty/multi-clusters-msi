# Authenticate and Authorize member cluster with hub cluster

## Install infrascture

```sh
  $ az aks create --resource-group demo --name hubCluster --node-count 1  --generate-ssh-keys --enable-aad --enable-azure-rbac

  $ az aks create --resource-group demo --name memberCluster --node-count 1  --generate-ssh-keys --enable-managed-identity
```

## Get the managed identity for the member cluster

- Get the member cluster vmss managed identity
  
  ```sh
  $ az aks show --resource-group demo --name memberCluster | jq .identityProfile.kubeletidentity.clientId
  ```
