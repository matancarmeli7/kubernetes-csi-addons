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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetVGRClass(vgr VGRInstance) (*replicationv1alpha1.VolumeGroupReplicationClass, error) {
	if vgr.Instance.Spec.VolumeGroupReplicationClass == "" {
		return nil, fmt.Errorf("VolumeGroupReplicationClass name is empty")
	}
	vgrClassName := vgr.Instance.Spec.VolumeGroupReplicationClass
	vgrClass := &replicationv1alpha1.VolumeGroupReplicationClass{}
	err := vgr.Client.Get(context.TODO(), types.NamespacedName{Name: vgrClassName}, vgrClass)
	if err != nil {
		if apierrors.IsNotFound(err) {
			vgr.Log.Error(err, fmt.Sprintf("(%v) VolumeGroupReplicationClass not found", vgrClassName))
		} else {
			vgr.Log.Error(err, fmt.Sprintf(
				"Got an unexpected error while fetching (%v) VolumeGroupReplicationClass", vgrClassName))
		}

		return nil, err
	}
	return vgrClass, nil
}

func GetSecretDataFromClass(client client.Client, vgrClass *replicationv1alpha1.VolumeGroupReplicationClass,
	logger logr.Logger) (map[string]string, error) {
	secretName, secretNamespace := getSecretFromClass(vgrClass)
	secret := make(map[string]string)
	var err error
	if secretName != "" && secretNamespace != "" {
		secret, err = getSecretData(client, logger, secretName, secretNamespace)
		if err != nil {
			return nil, err
		}
	}
	return secret, nil
}

func getSecretFromClass(vgrClass *replicationv1alpha1.VolumeGroupReplicationClass) (string, string) {
	secretName := vgrClass.Parameters[VGRSecretNameKey]
	secretNamespace := vgrClass.Parameters[VGRSecretNamespaceKey]
	return secretName, secretNamespace
}

func getClassDriver(vgrClass *replicationv1alpha1.VolumeGroupReplicationClass) string {
	if vgrClass.Driver != "" {
		return vgrClass.Driver
	}
	return ""
}

func getVGRDeletionPolicy(vgrClass *replicationv1alpha1.VolumeGroupReplicationClass,
) replicationv1alpha1.VolumeGroupReplicationDeletionPolicy {
	defaultDeletionPolicy := replicationv1alpha1.VolumeGroupReplicationContentDelete
	if vgrClass.VolumeGroupReplicationDeletionPolicy != nil {
		return *vgrClass.VolumeGroupReplicationDeletionPolicy
	}
	return defaultDeletionPolicy
}
