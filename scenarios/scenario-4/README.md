## Scenario-4 : Deploy A Single API Manager Instance with Analytics (Custom Pattern)

In this scenario we are deploying a single API Manager instance with API analytics (dashboard and worker profiles)in Kubernetes.

#### Deploy MySQL Database

```
 kubectl create -f mysql/
```

#### Deploy the configmap and the custom pattern

The configmap contains the configurations of API Manager node and Analytics

```
 kubectl create -f configmaps/
 kubectl create -f custom-pattern.yaml
```

#### Access API Manager

To access the API Manager, add a host mapping entry to the /etc/hosts file. As we have exposed the API portal service in Node Port type, you can use the IP address of any Kubernetes node.

```
<Any K8s Node IP>  wso2apim
```

- For Docker for Mac use "127.0.0.1" for the K8s node IP
- For Minikube, use minikube ip command to get the K8s node IP
- For GKE
    ```
    (kubectl get nodes -o jsonpath='{ $.items[*].status.addresses[?(@.type=="ExternalIP")].address }')
    ```
    - This will give the external IPs of the nodes available in the cluster. Pick any IP to include in /etc/hosts file.
  
- **API Publisher** : https://wso2apim:32001/publisher 
- **API Devportal** : https://wso2apim:32001/devportal 
- **API Analytics Dashboard**   : https://wso2apim-analytics:32201/analytics-dashboard 


## Custom Patterns

In scenario 2 you can find the pattern-1 deployment of API Manager. By giving the pattern numbers, you can deploy those patterns. In custom patterns, you can come with your own pattern. 

You can introduce a new custom pattern with any number of profiles.

There are several types of profiles:

    * api-manager
    * analytics-dashboard
    * analytics-worker


- You can create any no.of profiles with given types.
- You can provide desired configuration values if you wish, if not provided it will take upon the default type specific configuration values.
- Under a profile, everything except name, type, service, and configMaps are optional fields.
- You need to create relevant configmaps and persistent volume claims and add them under deploymentConfigmap, synapseConfigs and executionPlans. If you do not specify any persistent volume claims, default ones will not be created unlike in Pattern 1. You would not need a NFS deployment in this case.

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
        persistentVolumeClaim:
          synapseConfigs: <NEW_PROFILE_SYNAPSE_CONFIGS>
          executionPlans: <NEW_PROFILE_EXECUTION_PLANS>
```
