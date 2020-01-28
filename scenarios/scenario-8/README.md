## Scenario 8 : Exposing using Ingresses

1. Go inside root folder _wso2am-k8s-operator_

2. Create a new configmap **<AM-1-DEPLOYMENT-CONFIGMAP>** for API Manager instance 1 using the command,

```
kubectl create configmap <AM-1-DEPLOYMENT-CONFIGMAP> --from-file=wso2am-k8s-operator/scenarios/scenario-2/am-1/deployment.toml
```
3. Similarly, create a new configmap **<AM-2-DEPLOYMENT-CONFIGMAP>** for API Manager instance 1 using the command,
  
```
kubectl create configmap <AM-2-DEPLOYMENT-CONFIGMAP> --from-file=wso2am-k8s-operator/scenarios/scenario-2/am-2/deployment.toml
```
4. Follow steps 3,4,5 in the Home page

5. Then apply the given yaml using the command
```
kubectl apply -f scenarios/scenario-8/wso2-apim.yaml
```

Now WSO2 API Manager will be exposed via Ingresses and ClusterIP Service Type successfully.
