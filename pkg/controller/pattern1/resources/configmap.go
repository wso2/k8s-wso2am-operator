package resources

import (

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	v1 "k8s.io/api/core/v1"

)

func MakeConfigMap(apimanager *apimv1alpha1.APIManager,configMap *v1.ConfigMap) *corev1.ConfigMap {



	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "control-conf2",
			Namespace: apimanager.Namespace,
			Labels:   map[string]string {"deployment": "wso2am-pattern-1-am",
			"node": "wso2am-pattern-1-am-1",},

			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Data: map[string]string{
			//"am1-am2-deployment.toml":configMap.Data
		},
	}
}
