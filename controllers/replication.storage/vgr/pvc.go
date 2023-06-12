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
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	replicationv1alpha1 "github.com/csi-addons/kubernetes-csi-addons/apis/replication.storage/v1alpha1"
)

func GetMatchingPVCs(vgr VGRInstance) ([]replicationv1alpha1.ReplicatedPVC, error) {
	var matchingPVCs []replicationv1alpha1.ReplicatedPVC
	pvcList, err := ListPVCsByPVCSelector(vgr)
	if err != nil {
		return nil, err
	}
	for _, pvc := range pvcList.Items {
		isPVCShouldBeHandled, err := isPVCNeedToBeHandled(vgr, pvc)
		if err != nil {
			return nil, err
		}
		if isPVCShouldBeHandled && !isPVCInReplicatedPVCs(pvc.Name, matchingPVCs) {
			matchingPVCs = append(matchingPVCs, generateReplicatedPVC(pvc))
		}
	}
	return matchingPVCs, err
}

func ListPVCsByPVCSelector(vgr VGRInstance) (*corev1.PersistentVolumeClaimList, error) {
	pvcLabelSelector := getPVCLabelSelector(vgr)
	pvcSelector, err := metav1.LabelSelectorAsSelector(pvcLabelSelector)
	if err != nil {
		vgr.Log.Error(err, "Error with PVC label selector", "pvcSelector", pvcLabelSelector)
		return nil, fmt.Errorf("error with PVC label selector, %w", err)
	}
	vgr.Log.Info("Fetching PersistentVolumeClaims", "pvcSelector", pvcSelector)

	listOptions := []client.ListOption{
		client.InNamespace(vgr.Instance.Namespace),
		client.MatchingLabelsSelector{
			Selector: pvcSelector,
		},
	}

	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := vgr.Client.List(context.TODO(), pvcList, listOptions...); err != nil {
		vgr.Log.Error(err, "Failed to list PersistentVolumeClaims", "pvcSelector", pvcSelector)
		return nil, fmt.Errorf("failed to list PersistentVolumeClaims, %w", err)
	}

	vgr.Log.Info(fmt.Sprintf("Found %d PVCs using label selector %v", len(pvcList.Items), pvcSelector))
	return pvcList, nil
}

func isPVCNeedToBeHandled(vgr VGRInstance, pvc corev1.PersistentVolumeClaim) (bool, error) {
	if !isPVCInBoundState(pvc) {
		return false, nil
	}
	isPVCHasMatchingDriver, err := isPVCHasMatchingDriver(vgr.Log, vgr.Client, pvc, vgr.Driver)
	if err != nil {
		return false, err
	}
	return isPVCHasMatchingDriver, nil
	//isSCHasVGParam, err := IsPVCInStaticVG(reqLogger, client, pvc)
	//if err != nil {
	//	return false, err
	//}
	//if isSCHasVGParam {
	//	storageClassName, sErr := GetPVCClass(pvc)
	//	if sErr != nil {
	//		return false, sErr
	//	}
	//	msg := fmt.Sprintf(messages.StorageClassHasVGParameter, storageClassName, pvc.Namespace, pvc.Name)
	//	reqLogger.Info(msg)
	//	mErr := fmt.Errorf(msg)
	//	err = HandlePVCErrorMessage(reqLogger, client, pvc, mErr, addingPVC)
	//	if err != nil {
	//		return false, err
	//	}
	//	return false, nil
	//}
}

func isPVCInBoundState(pvc corev1.PersistentVolumeClaim) bool {
	return pvc.Status.Phase == corev1.ClaimBound
}

func isPVCHasMatchingDriver(logger logr.Logger, client client.Client,
	pvc corev1.PersistentVolumeClaim, driver string) (bool, error) {
	storageClassName := getPVCClass(pvc)
	if storageClassName == "" {
		return false, nil
	}
	scProvisioner, err := getStorageClassProvisioner(logger, client, storageClassName)
	if err != nil {
		return false, err
	}
	return scProvisioner == driver, nil
}

func isPVCInReplicatedPVCs(pvcName string, replicatedPVCs []replicationv1alpha1.ReplicatedPVC) bool {
	for _, replicatedPVC := range replicatedPVCs {
		if replicatedPVC.Name == pvcName {
			return true
		}
	}
	return false
}

func getVolumeIds(logger logr.Logger, client client.Client, pvcList []replicationv1alpha1.ReplicatedPVC) ([]string, error) {
	volumeIds := []string{}
	for _, pvc := range pvcList {
		pv, err := GetPVFromPVC(logger, client, pvc)
		if err != nil {
			return nil, err
		}
		if pv != nil {
			volumeIds = append(volumeIds, string(pv.Spec.CSI.VolumeHandle))
		}
	}
	return volumeIds, nil
}
