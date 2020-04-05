### Configure Own Custom Pattern

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
