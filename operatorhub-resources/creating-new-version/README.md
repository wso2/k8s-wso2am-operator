## Creating a new version

#### Prerequisites
[quay.io](<https://quay.io/repository/>) account

#### Steps

1. Create an operator bundle compatible with operator lifecycle manager.
* Bundle editor can be used for creation of the operator bundle.
* Add the CRDs, role, role binding, service account and deployment of the api-manager operator
* Include the description of the operator in metadata
* Include a quadratic image of the logo
* This will create a cluster service version (csv is inside the operator bundle)
* Download the operator bundle

2. You can preview the csv file created in the 1st step using [this](<https://operatorhub.io/preview>)

3. Add the directory with the new version to the api-manager-operator bundle available [here](<https://github.com/wso2-incubator/wso2am-k8s-operator/tree/master/operatorhub-resources/api-manager-operator>) and the bundle should be in the below format.

```bash
   api-manager-operator
   |- 1.0.0
   |  |- api-manager-operator.v1.0.0.clusterserviceversion.yaml
   |  |- apimanagers.apim.wso2.com.crd.yaml
   |- 1.0.1
   |  |- api-manager-operator.v1.0.1.clusterserviceversion.yaml
   |  |- apimanagers.apim.wso2.com.crd.yaml
   |- api-manager-operator.package.yaml
```
   
Note:

* api-operator.package.yaml contains the channels and the latest version available in those channel. If you want to make few versions available (Ex: 1.0.1 and 1.0.2 to be available to download), then you can put the new version under a new channel like below.
packageName: api-operator
```
channels:
  - name: stable
    currentCSV: api-manager-operator.v1.0.0
  - name: alpha
    currentCSV: api-manager-operator.v1.0.1
defaultChannel: stable
```

* If you want the subscribed customers to have a rolling update of the new version, put the new version in the same channel and put the old version in the 'replaces' field of the csv.
* If you are creating a new channel remove the 'replaces' field in the csv

4. Install OLM
```
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.10.0/crds.yaml
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/0.10.0/olm.yaml
```

5. Deploy operator-marketplace. (operator-marketplace operator is needed only for local testing)
```
git clone https://github.com/operator-framework/operator-marketplace.git
kubectl apply -f operator-marketplace/deploy/upstream/
```

6. Install operator courier
```
pip3 install operator-courier
```
6.Verify your bundle
```
thamilini:operatorhub-resources thamilini$ operator-courier --verbose verify  api-manager-operator
INFO:operatorcourier.verified_manifest:The source directory is in nested structure.
INFO:operatorcourier.verified_manifest:Parsing version: 1.0.1
INFO: Validating bundle. []
INFO: Validating custom resource definitions. []
INFO: Evaluating crd apimanagers.apim.wso2.com [1.0.1/apimanagers.apim.wso2.com.crd.yaml]
INFO: Validating cluster service versions. [1.0.1/apimanagers.apim.wso2.com.crd.yaml]
INFO: Evaluating csv api-manager-operator.v1.0.1 [1.0.1/api-manager-operator.v1.0.1.clusterserviceversion.yaml]
INFO: Validating packages. [1.0.1/api-manager-operator.v1.0.1.clusterserviceversion.yaml]
INFO: Evaluating package api-manager-operator [api-manager-operator/api-manager-operator.package.yaml]
INFO: Validating cluster service versions for operatorhub.io UI. [api-manager-operator/api-manager-operator.package.yaml]
INFO: Evaluating csv api-manager-operator.v1.0.1 [api-manager-operator/api-manager-operator.package.yaml]
```

7. Login to quay.io account
```
$ ./operator-courier/scripts/get-quay-token
Username: thamilini
Password:
{"token": "basic abvefjijkl=="}

export QUAY_TOKEN="basic abvefjijkl=="
```

8. Push the bundle to quay.io
```
export OPERATOR_DIR=api-manager-operator/
export QUAY_NAMESPACE=thamiliniram
export PACKAGE_NAME=api-manager-operator
export PACKAGE_VERSION=1.0.1
export TOKEN=$QUAY_TOKEN
operator-courier --verbose push "$OPERATOR_DIR" "$QUAY_NAMESPACE" "$PACKAGE_NAME" "$PACKAGE_VERSION" "$TOKEN"
```
Go to your quay.io account and check if the bundle is available in the applications tab.
#### Note: When you push the bundle, it will be private in your quay.io. Please make sure to change the repository to public

10. Test the operator locally following the instructions [here]()

11. Contribute to [operatorhub](https://github.com/operator-framework/community-operators)
* For Kubernetes operators which are available in operatorhub.io, contribute to [upstream operators](<https://github.com/operator-framework/community-operators/tree/master/upstream-community-operators>)
* For Openshift community operators, contribute to [community operators](<https://github.com/operator-framework/community-operators/tree/master/community-operators>)
