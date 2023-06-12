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

	replicationv1alpha1 "github.com/csi-addons/kubernetes-csi-addons/apis/replication.storage/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetPVFromPVC(logger logr.Logger, client client.Client, pvc replicationv1alpha1.ReplicatedPVC) (*corev1.PersistentVolume, error) {
	logger.Info(fmt.Sprintf("Get matching persistentVolume from %s/%s persistentVolumeClaim", pvc.Namespace, pvc.Name))
	pvName := pvc.VolumeName
	if pvName == "" {
		logger.Info("PersistentVolumeClaim does not Have persistentVolume")
		return nil, nil
	}

	pv, err := getPV(logger, client, pvName)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("%s persistentVolume does not exist", pvName)
		}
		return nil, err
	}
	return pv, nil
}

func getPV(logger logr.Logger, client client.Client, pvName string) (*corev1.PersistentVolume, error) {
	logger.Info(fmt.Sprintf("Getting %s persistentVolume", pvName))
	pv := &corev1.PersistentVolume{}
	namespacedPV := types.NamespacedName{Name: pvName}
	err := client.Get(context.TODO(), namespacedPV, pv)
	if err != nil {
		logger.Error(err, fmt.Sprintf("Failed to get %s persistentVolume", pvName))
		return nil, err
	}
	return pv, nil
}

func getPVNameFromPVC(pvc corev1.PersistentVolumeClaim) string {
	if pvc.Spec.VolumeName != "" {
		return pvc.Spec.VolumeName
	}
	return ""
}
