Deploying with simple Custom Resource file wso2-apim.yaml
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
```

That is it. The WSO2 API Manager has been deployed into K8S Cluster.
