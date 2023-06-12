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
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getPVCClass(pvc corev1.PersistentVolumeClaim) string {
	if class, found := pvc.Annotations[corev1.BetaStorageClassAnnotation]; found {
		return class
	}

	if pvc.Spec.StorageClassName != nil {
		return *pvc.Spec.StorageClassName
	}
	return ""
}

func getStorageClassProvisioner(logger logr.Logger, client client.Client, scName string) (string, error) {
	sc, err := getStorageClass(logger, client, scName)
	if err != nil {
		return "", err
	}
	return sc.Provisioner, nil
}

func getStorageClass(logger logr.Logger, client client.Client, scName string) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: scName}, sc)
	if err != nil {
		logger.Error(err, fmt.Sprintf("Failed to get (%s) storageClass", scName))
		return nil, err
	}
	return sc, nil
}
