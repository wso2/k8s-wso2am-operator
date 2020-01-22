package mysql

import (
	apimv1alpha1 "github.com/wso2-incubator/wso2am-k8s-operator/pkg/apis/apim/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "strings"
    "fmt"
)
 
type PvcConfig struct {
	 Name string
	 AccessModes        []v1.PersistentVolumeAccessMode
	 ClaimSize          string
	 StorageClass       *string
	 PVCConfigmap       string
 }

 // GetAccessModesFromString returns an array of AccessModes from a string created by GetAccessModesAsString
func GetAccessModesFromString(modes string) []v1.PersistentVolumeAccessMode {
	strmodes := strings.Split(modes, ",")
	accessModes := []v1.PersistentVolumeAccessMode{}
	for _, s := range strmodes {
		s = strings.Trim(s, " ")
		switch {
		case s == "ReadWriteOnce":
			accessModes = append(accessModes, v1.ReadWriteOnce)
		case s == "ReadOnlyMany":
			accessModes = append(accessModes, v1.ReadOnlyMany)
		case s == "ReadWriteMany":
			accessModes = append(accessModes, v1.ReadWriteMany)
		}
	}
	return accessModes
}

func AssignConfigMapValuesForMysqlPvc(apimanager *apimv1alpha1.APIManager, pvcConfWso2 *v1.ConfigMap) *PvcConfig {

    ControlConfigData := pvcConfWso2.Data
    
    fmt.Println(ControlConfigData)
    name, ok := ControlConfigData["wso2amMysqlPvcName"]
    if !ok {
        fmt.Println("sql error")
    }
	accessModes, _ := ControlConfigData["wso2amPvcAccessmode"]
	claimSize, _ := ControlConfigData["wso2amPvcMysqlStorage"]
	storageClass, _ := ControlConfigData["storageClassName"]

	pvcConfigMap := "pvc-config"

	pvc3cmvalues := &PvcConfig{
		Name:               name,
		AccessModes:        GetAccessModesFromString(accessModes),
		ClaimSize:          claimSize,
		StorageClass:       &storageClass,
		PVCConfigmap:       pvcConfigMap,
	}
	return pvc3cmvalues
}

func MakeMysqlPvc(apimanager *apimv1alpha1.APIManager, sqlconf *PvcConfig) *corev1.PersistentVolumeClaim {

	if len(sqlconf.AccessModes) == 0 {
		sqlconf.AccessModes = append(sqlconf.AccessModes, v1.ReadWriteOnce)
	}

	if len(sqlconf.ClaimSize) == 0 {
		sqlconf.ClaimSize = "20Gi"
	}

	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: sqlconf.Name,
			Namespace: apimanager.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(apimanager, apimv1alpha1.SchemeGroupVersion.WithKind("Apimanager")),
			},
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: sqlconf.AccessModes,
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					v1.ResourceStorage: resource.MustParse(sqlconf.ClaimSize),
				},
			},
			StorageClassName: sqlconf.StorageClass,
		},
	}
}

