## Scenario-13 : Deploy API Manager Pattern-4

#### Installation Prerequisiyes

* [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/).

In this scenario we are deploying API Manager Pattern-3 by using the following simple yaml definition. All API Manager servers are exposed via Ingress Controller.

```yaml
apiVersion: apim.wso2.com/v1alpha1
kind: APIManager
metadata:
  name: cluster-1
spec:
  pattern: Pattern-4
```
-----


#### Deploy pattern 4

```
kubectl create -f wso2-apim.yaml
```

#### Access API Manager

Get the external address of the ingresses using the command,

```

kubectl get ing

Output:
NAME                                   HOSTS                    ADDRESS           PORTS      AGE
wso2-am-analytics-dashboard-ingress    analytics.am.wso2.com    35.198.248.85     80, 443    36m
wso2-am-devportal-ingress              devportal.am.wso2.com    35.198.248.85     80, 443    36m
wso2-am-gw-external-ingress            ext.gateway.am.wso2.com  35.198.248.85     80, 443    36m
wso2-am-gw-internal-ingress            int.gateway.am.wso2.com  35.198.248.85     80, 443    36m
```

Add an **/etc/hosts/** entry as follows.

```
/etc/hosts
- - - - - - - - - - 
<EXTERNAL-ADDRESS>       analytics.am.wso2.com gateway.am.wso2.com publisher.am.wso2.com devportal.am.wso2.com
```

- **API Publisher & Devportal** : https://am.wso2.com 
- **API Analytics Dashboard**   : https://analytics.am.wso2.com
- **API Internal Gateway**   : https://int.gateway.am.wso2.com
- **API External Gateway**   : https://ext.gateway.am.wso2.com
