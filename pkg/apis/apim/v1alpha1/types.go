/*
 *
 *  * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *  *
 *  * WSO2 Inc. licenses this file to you under the Apache License,
 *  * Version 2.0 (the "License"); you may not use this file except
 *  * in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  * http:www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing,
 *  * software distributed under the License is distributed on an
 *  * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 *  * KIND, either express or implied. See the License for the
 *  * specific language governing permissions and limitations
 *  * under the License.
 *
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Apimanager is a specification for a Apimanager resource
type APIManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   APIManagerSpec   `json:"spec"`
	Status APIManagerStatus `json:"status"`
}

type PVCServer struct {
	SynapseConfigs string `json:"synapseConfigs"`
	ExecutionPlans string `json:"executionPlans"`
}

type Storage struct {
	Memory string `json:"memory"`
	CPU    string `json:"cpu"`
}

type Resources struct {
	Requests Storage `json:"requests"`
	Limits   Storage `json:"limits"`
}

type Probe struct {
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`
	FailureThreshold    int32 `json:"failureThreshold"`
	PeriodSeconds       int32 `json:"periodSeconds"`
}

type NewVolume struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type Configmap struct {
	DeploymentConfigmap string      `json:"deploymentConfigMap"`
	NewConfigmap        []NewVolume `json:"newConfigMap"`
}

type PVC struct {
	SynapseConfigs string      `json:"synapseConfigs"`
	ExecutionPlans string      `json:"executionPlans"`
	NewClaim       []NewVolume `json:"newClaim"`
}

type Deployment struct {
	Resources             Resources `json:"resources"`
	Replicas              *int32    `json:"replicas"`
	MinReadySeconds       int32     `json:"minReadySeconds"`
	Strategy              Strategy  `json:"strategy"`
	Image                 string    `json:"image"`
	ImagePullPolicy       string    `json:"imagePullPolicy"`
	LivenessProbe         Probe     `json:"livenessProbe"`
	ReadinessProbe        Probe     `json:"readinessProbe"`
	Configmaps            Configmap `json:"configMaps"`
	PersistentVolumeClaim PVC       `json:"persistentVolumeClaim"`
	SecurityContext       string    `json:"securityContext"`
}

type Strategy struct {
	Type          string        `json:"type"`
	RollingUpdate RollingUpdate `json:"rollingUpdate"`
}

type RollingUpdate struct {
	MaxSurge   int32 `json:"maxSurge"`
	MaxUnavail int32 `json:"maxUnavailable"`
}

type Ports struct {
	Name string `json:"name"`
	Port string `json:"port"`
}

type Service struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Ports []Ports `json:"ports"`
}
type Profile struct {
	Name       string     `json:"name"`
	Deployment Deployment `json:"deployment"`
	Type       string     `json:"type"`
	Service    Service    `json:"service"`
}

type Artifacts struct {
	Deployment Deployment `json:"deployment"`
}

// APIManagerSpec is the spec for a APIManager resource
type APIManagerSpec struct {
	// Here set UseMysql as string type to set true as default, and since using apiextensions.k8s.io/v1beta1
	// and not using apiextensions.k8s.io/v1
	UseMysql       string    `json:"useMysql"`
	AllowAnalytics string    `json:"allowAnalytics"`
	Pattern        string    `json:"pattern"`
	Replicas       *int32    `json:"replicas"`
	Profiles       []Profile `json:"profiles"`
	Service        Service   `json:"service"`
	Expose         string    `json:"expose"`
}

// APIManagerStatus is the status for a APIManager resource
type APIManagerStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// APIManagerList is a list of APIManager resources
type APIManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []APIManager `json:"items"`
}
