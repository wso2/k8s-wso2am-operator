## Scenario-9 : Deploy API Manager Pattern-2

#### Installation Prerequisiyes

* [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/).

In this scenario we are deploying API Manager Pattern-2 by using the following simple yaml definition. All API Manager servers are exposed via Ingress Controller.

```yaml
apiVersion: apim.wso2.com/v1alpha1
kind: APIManager
metadata:
  name: cluster-1
spec:
  pattern: Pattern-2
```
-----


#### Deploy pattern 2

```
kubectl create -f wso2-apim.yaml
```

#### Access API Manager

Get the external address of the ingresses using the command,

```

kubectl get ing

Output:
NAME                                  HOSTS                   ADDRESS       PORTS     AGE
wso2-am-ingress                       am.wso2.com             34.67.188.5   80,443    443 6m17s
wso2-am-gw-ingress                    gateway.am.wso2.com     34.67.188.5   80,443    443 6m17s
wso2-am-analytics-dashboard-ingress   analytics.am.wso2.com   34.67.188.5   80,443    443 6m17s
``` 

Add an **/etc/hosts/** entry as follows.

```
/etc/hosts
- - - - - - - - - - 
<EXTERNAL-ADDRESS>       analytics.am.wso2.com gateway.am.wso2.com am.wso2.com
```

- **API Publisher** : https://am.wso2.com/publisher 
- **API Devportal** : https://am.wso2.com/devportal 
- **API Gateway**   : https://gateway.am.wso2.com
- **API Analytics Dashboard**   : https://analytics.am.wso2.com


