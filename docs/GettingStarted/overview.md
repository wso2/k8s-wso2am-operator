## Overview WSO2AM-K8S-Operator

WSO2 AM K8S Operator provides a fully automated experience for cloud-native API management. A user can deploy WSO2 API Manager different patterns in K8S cluster through a simple command. 

For this, WSO2 AM K8S Operator introduced a new custom resource definition(CRD) called **APIManager**

### Custom resource: APIManager

- `APIManager` holds API-related information. The following YAML shows the simplest custom resource for deploying APIManager in Kubernetes.

    ```
    apiVersion: apim.wso2.com/v1alpha1
    kind: APIManager
    metadata:
      name: cluster-1
    spec:
      pattern: Pattern-1
    ```
- The above CRD has corresponding custom controller. Custom controllers are the “brains” behind the custom resources. 

### Custom Controller: wso2am-controller

- APIManager controller is quite complex compared to other controllers. It has the main task to create Kubernetes artifacts and deploy them into Kubernetes cluster.

- When the wso2am controller is triggered, it starts to create the deployments and services of relevant Kubernetes artifacts based on specified pattern, and finally deploys the entire API Manager in Kubernetes cluster.

- As you can see, wso2am controller has taken out all the complexity from DevOps and automates deployment with all the best practices required to deploy API Manager.
