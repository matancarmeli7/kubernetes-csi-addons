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

	replicationv1alpha1 "github.com/csi-addons/kubernetes-csi-addons/apis/replication.storage/v1alpha1"
	vrutils "github.com/csi-addons/kubernetes-csi-addons/controllers/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
)

type VGRInstance struct {
	EventRecorder        *vrutils.EventReporter
	Client               client.Client
	Ctx                  context.Context
	Log                  logr.Logger
	Instance             *replicationv1alpha1.VolumeGroupReplication
	SavedInstanceStatus  replicationv1alpha1.VolumeGroupReplicationStatus
	VolRepPVCs           []corev1.PersistentVolumeClaim
	VolSyncPVCs          []corev1.PersistentVolumeClaim
	ReplClassList        *replicationv1alpha1.VolumeReplicationClassList
	VrgObjectProtected   *metav1.Condition
	KubeObjectsProtected *metav1.Condition
	VrcUpdated           bool
	NamespacedName       string
	Driver               string
}
