### Configure Own Custom Pattern

The user has the flexibility to create a completely new Pattern they wish.

They must provide the type for each of the profiles. And then based on the types, a template for that type of profile will be created.
But the user can specify whatever the fields they want to override without taking the default configuratiosn specific to that type.

For the above profiles, user can override the fields such as,

* Replicas
* MinReadySeconds
* Resources 
  * Requests 
    * Memory 
    * CPU
  * Limits 
    * Memory 
    * CPU
* LivenessProbe
  - InitialDelaySeconds
  - PeriodSeconds
  - FailureTHreshold
* ReadinessProbe
  - InitialDelaySeconds
  - PeriodSeconds
  - FailureTHreshold
* imagePullPolicy


Also they can add the Newly created configmaps and persistent volume claims to override the default ones.

```
  configMaps:
          deploymentConfigMap: <NEW_PROFILE_DEPLOYMENT_CONFIGMAP>
  persistentVolumeClaim:
          synapseConfigs: <NEW_PROFILE_SYNAPSE_CONFIGS_PVC>
          executionPlans: <NEW_PROFILE_EXECUTION_PLANS_PVC>
```
Also, there is an additional functionality to add new configmap and persistent volumes with new MountPaths as well.
```
  configMaps:
          newConfigMap:
              - name: <NEW_CONFIGMAP_1>
                mountapth: <NEW_MOUNTPATH_FOR_CONFIGMAP_1>
  persistentVolumeClaim:
          newClaim:
              - name: <NEW_CLAIM_1>
                mountpath: <NEW_MOUNTPATH_FOR_CLAIM_1>
```
Finally via applying that custom resource YAML file, it will automatically start to create all the artifacts for the newly created pattern.

The scenario for Customizing the new pattern can be seen in [Scenario-6](https://github.com/wso2-incubator/wso2am-k8s-operator/tree/master/scenarios/scenario-6)

