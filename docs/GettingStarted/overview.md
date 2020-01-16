## Overview WSO2AM-K8S-Operator

WSO2 AM K8S Operator provides a fully automated experience for cloud-native API management. A user can deploy the entire WSO2 API Manager based of different patterns into K8S cluster through a simple command. 

For this, WSO2 AM K8S Operator introduced a new custom resource definition(CRD) called **APIManager**

### Custom resource: APIManager
`APIManager` holds API-related information. You can see the APIManager definition and data structure for APIManager here, with cluster name and pattern name. The following YAML shows sample payload for APIManager.

```
apiVersion: apim.wso2.com/v1alpha1
kind: APIManager
metadata:
  name: cluster-1
spec:
  pattern: Pattern-1
```
The above CRD has corresponding custom controller. Custom controllers are the “brains” behind the custom resources. 

### Custom Controller: wso2am-controller

APIManager controller is quite complex compared to other controllers. It has the main task to Create Kubernetes artifacts and deploy them into Kubernetes cluster.

When the API custom controller is triggered, it will start to create the  deployements and services of relevant kubernetes artifacts based on specified pattern, and finally deploys the entire APIManager into K8S cluster.

As you can see, API Controller has taken out all the complexity from DevOps and automates deployment with all the best practices required to deploy APIManager.
