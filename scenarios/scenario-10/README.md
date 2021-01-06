## Scenario-10 : Override default configuration values

You can override the default configuration values of relavant artifacts.

There are different set of profiles based on each patterns. Some of them are:

In pattern-2 we have,

* api-pub-dev-tm-1
* api-pub-dev-tm-2
* analytics-dashboard
* analytics-worker
* api-keymanager
* api-gateway

For the above profiles, you can override the fields such as,

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
* securityContext

You can specify to any or all of the profiles using the given wso2-apim.yaml file. A sample configuration values for one of the profile is given, you can include required profiles as an array. Then apply the command,

#### Deploy Pattern-2 by overriding configurations

Please follow the prerequisites section in scenario 10 to deploy Pattern-2 and execute the following command.

```
  kubectl apply -f wso2-apim.yaml
```
