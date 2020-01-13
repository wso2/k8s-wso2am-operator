## Scenario 1

Deploying WSO2 API Manager with a simple Custom Resource Definition.
```
apiVersion: apim.wso2.com/v1alpha1
kind: APIManager
metadata:
  name: cluster-1
spec:
  pattern: Pattern-1
```

Apply the above yaml file using the command
```
kubectl apply -f wso2-apim.yaml

Output:
apimanager.apim.wso2.com/cluster-1 created

```

That is it. The WSO2 API Manager has been deployed into K8S Cluster.
