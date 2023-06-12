/*
Copyright 2023 The Kubernetes-CSI-Addons Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vgr

import (
	replicationv1alpha1 "github.com/csi-addons/kubernetes-csi-addons/apis/replication.storage/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	corev1 "k8s.io/api/core/v1"
)

func generateReplicatedPVC(pvc corev1.PersistentVolumeClaim) replicationv1alpha1.ReplicatedPVC {
	scName := getPVCClass(pvc)
	return replicationv1alpha1.ReplicatedPVC{
		Name:             pvc.Name,
		Namespace:        pvc.Namespace,
		VolumeName:       getPVNameFromPVC(pvc),
		StorageClassName: &scName,
		Labels:           pvc.Labels,
		AccessModes:      pvc.Spec.AccessModes,
		Resources:        pvc.Spec.Resources,
	}
}

func IsReplicatedPVCsEqual(x []replicationv1alpha1.ReplicatedPVC, y []replicationv1alpha1.ReplicatedPVC) bool {
	less := func(a, b corev1.PersistentVolumeClaim) bool { return a.Name < b.Name }
	equalIgnoreOrder := cmp.Diff(x, y, cmpopts.SortSlices(less)) == ""
	return equalIgnoreOrder
}
