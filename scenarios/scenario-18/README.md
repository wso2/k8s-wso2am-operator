## Scenario-18 : Expose API Manager using Ingress (Custom Pattern)

In this scenario we are exposing a single API Manager instance in Kubernetes using Ingress.

#### Deploy the configmap and the custom pattern

The configmap contains the configurations of API Manager node. 

```
 kubectl create -f wso2am-apim-configmap.yaml
 kubectl create -f custom-pattern.yaml
```

#### Access API Manager

To access the API Manager, consider the following steps.

  1. Get the external address of the ingresses using the command,
    ```
    kubectl get ingress
    
    Output:
    NAME                                     HOSTS                ADDRESS           PORTS     AGE
    wso2-am-gateway-px-ingress               wso2apim-gateway     104.154.172.7     80, 443   3m18s
    wso2-am-px-ingress                       wso2apim             104.154.172.7     80, 443   3m18s
    ```
        
  2. Add an **/etc/hosts/** entry as follows.
        
    ```
    /etc/hosts
    ----------
    <EXTERNAL-ADDRESS>       wso2apim wso2apim-gateway wso2apim-analytics    
  
- **API Publisher** : https://wso2apim/publisher 
- **API Devportal** : https://wso2apim/devportal 


## Custom Patterns

You can introduce a new custom pattern with any number of profiles.

There are several types of profiles:

    * api-manager
    * analytics-dashboard
    * analytics-worker


- You can create any no.of profiles with given types.
- You can provide desired configuration values if you wish, if not provided it will take upon the default type specific configuration values.
- Under a profile, everything except name, type, service, and configMaps are optional fields.

#### Available options in a custom pattern can be found here.

```
apiVersion: apim.wso2.com/v1alpha1
kind: APIManager
metadata:
  name: cluster-1
spec:
  pattern: Pattern-X
  service:
    type: NodePort
  useMysql: "<true|false>"
  profiles:
    - name: <NEW_PROFILE_NAME>
      type: api-manager
      service:
        name: <NEW_PROFILE_SERVICE_NAME>
      deployment:
        replicas: 1
        minReadySeconds: 100
        imagePullPolicy: Always
        resources:
          requests:
            memory: 2Gi
            cpu: 2000m
          limits:
            memory: 3Gi
            cpu: 3000m
        livenessProbe:
          initialDelaySeconds: 240
          failureThreshold: 3
          periodSeconds: 10
        readinessProbe:
          initialDelaySeconds: 240
          failureThreshold: 3
          periodSeconds: 10
        configMaps:
          deploymentConfigMap: <NEW_PROFILE_DEPLOYMENT_CONFIGMAP>
```