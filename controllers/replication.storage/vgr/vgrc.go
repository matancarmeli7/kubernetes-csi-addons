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
	vrutils "github.com/csi-addons/kubernetes-csi-addons/controllers/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func GenerateVGRCName(vgUID string) (string, error) {
	if len(vgUID) == 0 {
		return "", fmt.Errorf("Corrupted volumeGroup object, it is missing UID")
	}
	return fmt.Sprintf("%s-%s", VGRCNamePrefix, vgUID), nil
}

func GenerateVGRC(vgr *replicationv1alpha1.VolumeGroupReplication, vgrcName string,
	vgrClass *replicationv1alpha1.VolumeGroupReplicationClass,
) *replicationv1alpha1.VolumeGroupReplicationContent {
	secretName, secretNamespace := getSecretFromClass(vgrClass)
	return &replicationv1alpha1.VolumeGroupReplicationContent{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vgrcName,
			Namespace: vgr.Namespace,
		},
		Spec: generateVGRCSpec(vgr, vgrClass, secretName, secretNamespace),
	}
}

func generateVGRCSpec(vgr *replicationv1alpha1.VolumeGroupReplication,
	vgrClass *replicationv1alpha1.VolumeGroupReplicationClass,
	secretName string, secretNamespace string) replicationv1alpha1.VolumeGroupReplicationContentSpec {
	return replicationv1alpha1.VolumeGroupReplicationContentSpec{
		VolumeGroupReplicationClassName:      &vgrClass.Name,
		VolumeGroupReplicationRef:            vrutils.GenerateObjectReference(vgr),
		Driver:                               getClassDriver(vgrClass),
		VolumeGroupReplicationDeletionPolicy: getVGRDeletionPolicy(vgrClass),
	}
}

func CreateVGRC(vgr VGRInstance, vgrc *replicationv1alpha1.VolumeGroupReplicationContent) error {
	err := vgr.Client.Create(context.TODO(), vgrc)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			vgr.Log.Info(fmt.Sprintf("(%s/%s) VolumeGroupReplicationContent is already exists",
				vgrc.Namespace, vgrc.Name))
			return nil
		}
		vgr.Log.Error(err, fmt.Sprintf("(%s/%s)VolumeGroupReplicationContent creation failed",
			vgrc.Namespace, vgrc.Name))
		return err
	}
	return err
}

func IsVGRCReady(vgr VGRInstance, vgrc *replicationv1alpha1.VolumeGroupReplicationContent) (bool, error) {
	vgrcObj, err := GetVGRC(vgr, vgrc.Name, vgrc.Namespace)
	if err != nil {
		if !errors.IsNotFound(err) {
			return false, err
		}
		return false, nil
	}
	return *vgrcObj.Status.Ready, nil
}

func GetVGRC(vgr VGRInstance, vgrcName, vgrcNamespace string) (*replicationv1alpha1.VolumeGroupReplicationContent, error) {
	if vgrcName == "" {
		group := schema.ParseGroupResource(fmt.Sprintf("%s/%s", VGRCKind, ReplicationGroup))
		return nil, apierrors.NewNotFound(group, "")
	}
	namespacedName := types.NamespacedName{Name: vgrcName, Namespace: vgrcNamespace}

	vgr.Log.Info(fmt.Sprintf("Getting (%s/%s) volumeGroupContent", vgrcNamespace, vgrcName))
	vgrc := &replicationv1alpha1.VolumeGroupReplicationContent{}
	err := vgr.Client.Get(context.TODO(), namespacedName, vgrc)
	if err != nil {
		if apierrors.IsNotFound(err) {
			vgr.Log.Error(err, fmt.Sprintf("(%s/%s) VolumeGroupReplicationContent not found", vgrcNamespace, vgrcName))
		} else {
			vgr.Log.Error(err, fmt.Sprintf(
				"Got an unexpected error while fetching (%s/%s) VolumeGroupReplicationContent", vgrcNamespace, vgrcName))
		}
		return nil, err
	}
	return vgrc, nil
}
