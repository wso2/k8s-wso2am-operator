### Scenario 2

1. Go inside root folder wso2am-k8s-operator

2. Create a new configmap with the name of **wso2am-pattern-1-am-1-conf** using given deployment.toml file using the command,

```
kubectl create configmap wso2am-pattern-1-am-1-conf --from-file=wso2am-k8s-operator/scenarios/scenario-2/deployment.toml
```

2. Follow steps 3,4,5 in the Home page

3. Then apply the given yaml using the command

```
kubectl apply -f kubectl apply -f scenarios/scenario-2/wso2-apim.yaml
```

Now WSO2 API Manager will be exposed via NodePort Service Type successfully.
