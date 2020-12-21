## Scenario 8 : Expose API Manager using Ingress

You can expose API Manager servers via Ingress Controller. 

1. Install the [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/deploy/).

2. Deploy configuration updated configmaps.

    ```
    kubectl apply -f configmaps/
    ```

3. Deploy API Manager

    ```
    kubectl apply -f wso2-apim.yaml
    ```

4. Access API Manager

    Get the external address of the ingresses using the command,
    ```
    kubectl get ingress
    
    Output:
    NAME                                     HOSTS                ADDRESS           PORTS     AGE
    wso2-am-analytics-dashboard-p1-ingress   wso2apim-analytics   104.154.172.7     80, 443   3m18s
    wso2-am-gateway-p1-ingress               wso2apim-gateway     104.154.172.7     80, 443   3m18s
    wso2-am-p1-ingress                       wso2apim             104.154.172.7     80, 443   3m18s
    ```
        
    Add an **/etc/hosts/** entry as follows.
        
    ```
    /etc/hosts
    ----------
    <EXTERNAL-ADDRESS>       wso2apim wso2apim-gateway wso2apim-analytics              
    ```
        

- **API Publisher** : https://wso2apim/publisher 
- **API Devportal** : https://wso2apim/devportal 
- **API Analytics Dashboard**   : https://wso2apim-analytics/analytics-dashboard 
