## Scenario-6 : Custom Pattern

You can introduce a new custom pattern with any no.of profiles.

There are several types a profile can be of, they are:

* api-manager
* analytics-dashboard
* analytics-worker

You can create any no.of profiles with given types.
You can provide desired configuration values if you wish, if not provided it will take upon the default type specific configuration values.
Under a profile, everything except name, type, service, and configMaps are optional fields.

There are two main types of patterns you can implement. Create the necessary profiles for the pattern you want to try out.
* api-portal
  * api-manager
  * mysql
* api-portal-with-analytics
  * api-manager
  * analytics-dashboard
  * analytics-worker
  * mysql

You need to create relavant configmaps and persistent volume claims and add them under deploymentConfigmap, synapseConfigs and executionPlans. If you do not specify any persistent volume claims, default ones will not be created unlike in Pattern 1. You would not need a nfs deployment in this case.

If you choose not to specify any persistent volume claims, apply the given `.yaml` files to deploy a seperate MySQL service.

```
  kubectl apply -f artifacts/install/api-manager-artifacts/pattern-x/mysql
```

Finally, apply the given wso2-apim.yaml file.

```
  kubectl apply -f scenarios/scenario-6/wso2-apim.yaml
```

