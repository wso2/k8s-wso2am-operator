## Scenario 3 : Deploy API Manager Pattern-1 (NodePort Service Type)

1. Go inside root folder _wso2am-k8s-operator_

2. Create a new configmap **wso2am-pattern-1-am-1-conf** for API Manager instance 1 using the command,

```
kubectl create configmap wso2am-pattern-1-am-1-conf --from-file=wso2am-k8s-operator/scenarios/scenario-2/am-1/deployment.toml
```
3. Similarly, create a new configmap **wso2am-pattern-1-am-2-conf** for API Manager instance 1 using the command,
```
kubectl create configmap wso2am-pattern-1-am-2-conf --from-file=wso2am-k8s-operator/scenarios/scenario-2/am-2/deployment.toml
```
4. Follow steps 3,4,5 in the Home page

5. Then apply the given yaml using the command
```
kubectl apply -f scenarios/scenario-2/wso2-apim.yaml
```

- GCP:
    Get the external-ip of one of the nodes in the cluster using the command,
    
    ```
        kubectl get nodes -o wide
        
    ```
    Then add that ip to the **/etc/hosts** file,
    ```
    /etc/hosts
    ----------
    <EXTERNAL-IP-OF-ONE-OF-THE-NODES>                       wso2apim
    <EXTERNAL-IP-OF-ONE-OF-THE-NODES>    wso2apim-analytics-dashboard 
        
    ```
 - Minikube:
      Get the ip address of the Minikube, using the command,
      ```
          minikube ip
      ```
      Then add that ip address to the **/etc/hosts** file,
      ```
      /etc/hosts
      ----------
      <MINIKUBE-IP>      wso2apim
      ```
    
    
    
  Finally, WSO2 API Manager will be exposed via NodePort Service Type successfully.
  
  Access the portals using below urls.
      
   _APIM Publisher_ - https://wso2apim:32001/publisher
   
   _APIM Devportal_ - https://wso2apim:32001/devportal

