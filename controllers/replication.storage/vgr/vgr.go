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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	replicationv1alpha1 "github.com/csi-addons/kubernetes-csi-addons/apis/replication.storage/v1alpha1"
	vrutils "github.com/csi-addons/kubernetes-csi-addons/controllers/utils"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"
)

func UpdateVGRSourceContent(vgr VGRInstance, vgrcName string) error {
	vgr.Instance.Spec.VolumeGroupReplicationContentName = &vgrcName
	if err := vrutils.UpdateObject(vgr.Client, vgr.Instance); err != nil {
		vgr.Log.Error(err, fmt.Sprintf("failed to update (%s/%s) VolumeGroupReplication source",
			vgr.Instance.Namespace, vgr.Instance.Name))
		return err
	}
	return nil
}

func UpdateVGRStatus(vgr VGRInstance, vgrcName string, ready bool) error {
	vgrObj := vgr.Instance
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		vgrObj.Status.BoundVolumeGroupReplicationContentName = &vgrcName
		vgrObj.Status.ObservedGeneration = vgrObj.Generation
		vgrObj.Status.Ready = &ready
		err := vgrRetryOnConflictFunc(vgr)
		return err
	})
	if err != nil {
		return err
	}

	return updateVGRStatus(vgr)
}

func vgrRetryOnConflictFunc(vgr VGRInstance) error {
	vgrObj := vgr.Instance
	err := updateVGRStatus(vgr)
	if apierrors.IsConflict(err) {
		uErr := vrutils.GetNamespacedObject(vgr.Client, vgrObj)
		if uErr != nil {
			return uErr
		}
		vgr.Log.Info(fmt.Sprintf("Retry update (%s/%s) volumeGroupReplication status due to conflict error",
			vgrObj.Namespace, vgrObj.Name))
	}
	return err
}

func updateVGRStatus(vgr VGRInstance) error {
	vgrObj := vgr.Instance
	vgr.Log.Info(fmt.Sprintf("Updating status of (%s/%s) volumeGroupReplication", vgrObj.Namespace, vgrObj.Name))
	if err := vrutils.UpdateObjectStatus(vgr.Client, vgrObj); err != nil {
		if apierrors.IsConflict(err) {
			return err
		}
		vgr.Log.Error(err, fmt.Sprintf("failed to update (%s/%s) volumeGroupReplication status",
			vgrObj.Namespace, vgrObj.Name))
		return err
	}
	return nil
}

func getPVCLabelSelector(vgr VGRInstance) *metav1.LabelSelector {
	if vgr.Instance.Spec.PVCSelector != nil {
		return vgr.Instance.Spec.PVCSelector
	}
	return &metav1.LabelSelector{}
}

func getVGRId(vgr VGRInstance) (string, error) {
	vgrObj := vgr.Instance
	vgrcName := ""
	if vgrObj.Spec.VolumeGroupReplicationContentName != nil {
		vgrcName = *vgrObj.Spec.VolumeGroupReplicationContentName
	}
	vgrc, err := GetVGRC(vgr, vgrcName, vgrObj.Namespace)
	if err != nil {
		return "", err
	}
	handle := vgrc.Spec.VolumeGroupReplicationHandle
	return *handle, nil
}

func UpdateReplicatedPVCsAndPVs(vgr VGRInstance, matchingPvcs []replicationv1alpha1.ReplicatedPVC) error {
	vgrObj := vgr.Instance
	vgrReplicatedPVCs := make([]replicationv1alpha1.ReplicatedPVC, len(vgrObj.Status.ReplicatedPVCs))
	copy(vgrReplicatedPVCs, vgrObj.Status.ReplicatedPVCs)

	for _, pvc := range vgrReplicatedPVCs {
		if !isPVCInReplicatedPVCs(pvc.Name, matchingPvcs) {
			//err := RemoveVolumeFromPvcListAndPvList(logger, client, driver, pvc, vg)
			//if err != nil {
			//	return err
			//}
		}
	}
	for _, pvc := range matchingPvcs {
		if !isPVCInReplicatedPVCs(pvc.Name, vgrReplicatedPVCs) {
			vgr.Log.Info(fmt.Sprintf("adding %s/%s to ReplicatedPVCs in %s/%s vgr",
				pvc.Namespace, pvc.Name, vgrObj.Namespace, vgrObj.Name))
			//err := AddVolumeToPvcListAndPvList(logger, client, &pvc, vg)
			//if err != nil {
			//	return err
			//}
		}
	}
	return nil
}
