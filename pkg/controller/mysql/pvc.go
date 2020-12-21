/*
 *
 *  * Copyright (c) 2020 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
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
package mysql

import (
	"strings"

	v1 "k8s.io/api/core/v1"
)

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
