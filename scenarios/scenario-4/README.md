## Scenario-4 : Override ConfigMaps and PersistentVolumeClaims

You can override the default deployemnt configmap present for each profile.

1. Create a new deployment configmap using the deployment.toml/ deployment.yaml files. And apply the command.

```
  kubectl create configmap <NEW-PROFILE-CONFIGMAP> --from-file=<PATH_TO_DEPLOYMENT.TOML>
```

2. Then in the given wso2-apim.yaml file replace, the **deploymentConfigmap** name to the above created name.

Similarly, you can create the persistent volume claims for synapse configs and execution plans. 
And add the newly created pvc names in relevant places of **synapseConfigs** and **executionPlans**.

3. Finally apply the command,

```
  kubectl apply -f scenarios/scenario-4/wso2-apim.yaml
```

