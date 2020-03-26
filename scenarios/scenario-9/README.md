## Scenario 9 : Exposing using Ingresses

1. Go inside root folder _wso2am-k8s-operator_
2. Install the [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/).

3. Create a new configmaps for  APIM instances and Analytics Dashboard instance using the following commands.

```
kubectl apply -f scenarios/scenario-9/am-1/wso2am-conf.yaml
kubectl apply -f scenarios/scenario-9/am-2/wso2am-conf.yaml
kubectl apply -f scenarios/scenario-9/dash/wso2am-dash-conf.yaml
```

4. Follow steps 3,4,5 in the Home page

5. Then apply the given yaml using the command
```
kubectl apply -f scenarios/scenario-9/wso2-apim.yaml
```

Get the external address of the ingresses using the command,
```
kubectl get ingress

Output:
NAME                                     HOSTS                          ADDRESS         PORTS     AGE
wso2-am-analytics-dashboard-p1-ingress   wso2apim-analytics-dashboard   34.93.244.141   80, 443   24m
wso2-am-gateway-p1-ingress               wso2apim-gateway               34.93.244.141   80, 443   24m
wso2-am-p1-ingress                       wso2apim                       34.93.244.141   80, 443   24m
```
        
Then add add those ingresses with the Host Names and Addresses obtained in **/etc/hosts/**,
    
```
/etc/hosts
----------
<EXTERNAL-ADDRESS>       wso2apim-analytics-dashboard              
<EXTERNAL-ADDRESS>       wso2apim-gateway
<EXTERNAL-ADDRESS>       wso2apim 
```
        

Now WSO2 API Manager will be exposed via Ingresses and ClusterIP Service Type successfully.

Access the portals using below urls.

   _APIM Publisher_ - https://wso2apim:443/publisher

   _APIM Devportal_ - https://wso2apim:443/devportal

