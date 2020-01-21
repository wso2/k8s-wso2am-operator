## Scenario-6 : Custom pattern

You can introduce a new custom pattern with any no.of profiles.

There are several types a profile can be of, they are:

* api-manager
* analytics-dashboard
* analytics-worker

You can create any no.of profiles with given types.
You can provide desired configuration values if you wish, if not provided it will take upon the default type specific configuration values.

You need to create relavant configmaps and persistent volume claims and add them under deploymentConfigmap, synapseConfigs and executionPlans

Finally, apply the given wso2-apim.yaml file.

```
  kubectl apply -f scenarios/scenario-6/wso2-apim.yaml
```

