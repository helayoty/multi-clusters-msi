# Authenticate and Authorize member cluster with hub cluster

- Install infrascture

```sh
  $ az aks create --resource-group demo --name hubCluster --node-count 1  --generate-ssh-keys --enable-aad --enable-azure-rbac

  $ az aks create --resource-group demo --name memberCluster --node-count 1  --generate-ssh-keys --enable-managed-identity
```

- Create namespace in the member cluster, aka `member-a`.
- Create Role and RoleBinding

```sh
 $ kubectl apply -f hub_roles.yaml
```

- Get the member cluster vmss managed identity
  
  ```sh
  $ az aks show --resource-group demo --name memberCluster | jq .identityProfile.kubeletidentity.clientId
  ```

- Create new managed identity Credential using the member cluster clientId
- Use AKS scope `6dae42f8-4368-4678-94ff-3960e28e3630` to create TokenRequestOptions.
- Get the token from the created managed identity.
- Create rest config using the following:
  - BearerToken --> token
  - Host -->   by running `kubectl config view -o jsonpath='{.clusters[0].cluster.server}'`
  - set TLS as insecure = `true`
- Create clientSet using this config
- Do any operation to confirm.

- Build your image

```sh
  docker build  . -t <docker_repo>/member-agent:v1
  docker push <docker_repo>/member-agent:v1
```

- Deploy the agent/app to member cluster

```sh
 $ kubectl apply -f member-agent-deployment.yaml
```

