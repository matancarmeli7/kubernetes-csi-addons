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

package predicate

import (
	"context"
	"reflect"

	replicationv1alpha1 "github.com/csi-addons/kubernetes-csi-addons/apis/replication.storage/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func FinalizerPredicate() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return !reflect.DeepEqual(e.ObjectNew.GetFinalizers(), e.ObjectOld.GetFinalizers())
		},
	}
}

func PVCPredicateFunc() predicate.Funcs {
	pvcPredicate := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			log := ctrl.Log.WithName("pvcmap").WithName("VolumeReplicationGroup")
			log.Info("Create event for PersistentVolumeClaim")
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			log := ctrl.Log.WithName("pvcmap").WithName("VolumeReplicationGroup")
			oldPVC, ok := e.ObjectOld.DeepCopyObject().(*corev1.PersistentVolumeClaim)
			if !ok {
				log.Info("Failed to deep copy older PersistentVolumeClaim")

				return false
			}
			newPVC, ok := e.ObjectNew.DeepCopyObject().(*corev1.PersistentVolumeClaim)
			if !ok {
				log.Info("Failed to deep copy newer PersistentVolumeClaim")

				return false
			}

			log.Info("Update event for PersistentVolumeClaim")
			return isUpdateEventShouldBeProcessed(oldPVC, newPVC, log)
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
	}

	return pvcPredicate
}

func isUpdateEventShouldBeProcessed(oldPVC, newPVC *corev1.PersistentVolumeClaim, log logr.Logger) bool {
	pvcNamespacedName := types.NamespacedName{Name: newPVC.Name, Namespace: newPVC.Namespace}
	predicateLog := log.WithValues("pvc", pvcNamespacedName.String())

	if isLabelsChanged(oldPVC, newPVC) {
		predicateLog.Info("Reconciling due to change in the labels")
		return true
	}

	if isPhaseChanged(oldPVC, newPVC) && newPVC.Status.Phase == corev1.ClaimBound {
		predicateLog.Info("Reconciling due to change in the labels")
		return true
	}

	predicateLog.Info("Not Requeuing", "oldPVC Phase", oldPVC.Status.Phase,
		"newPVC phase", newPVC.Status.Phase)
	return false
}

func isLabelsChanged(oldObject, newObject runtimeclient.Object) bool {
	return !reflect.DeepEqual(oldObject.(*corev1.PersistentVolumeClaim).Labels,
		newObject.(*corev1.PersistentVolumeClaim).Labels)
}

func isPhaseChanged(oldObject, newObject runtimeclient.Object) bool {
	return !reflect.DeepEqual(oldObject.(*corev1.PersistentVolumeClaim).Status.Phase,
		newObject.(*corev1.PersistentVolumeClaim).Status.Phase)
}

func CreateRequests(client runtimeclient.Client) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(
		func(object runtimeclient.Object) []reconcile.Request {
			var vgList replicationv1alpha1.VolumeGroupReplicationList
			if err := client.List(context.TODO(), &vgList); err != nil {
				return []ctrl.Request{}
			}
			// TODO - add a label selector check to the VolumeGroupReplication to filter the list
			// Create a reconcile request for each matching VolumeGroupReplication.
			requests := make([]ctrl.Request, len(vgList.Items))
			for _, vg := range vgList.Items {
				requests = append(requests, ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: vg.Namespace,
						Name:      vg.Name,
					},
				})
			}
			return requests
		})
}
