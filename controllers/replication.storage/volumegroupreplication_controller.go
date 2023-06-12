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

package controllers

import (
	"context"
	"fmt"
	"time"

	replicationv1alpha1 "github.com/csi-addons/kubernetes-csi-addons/apis/replication.storage/v1alpha1"
	vrpredicate "github.com/csi-addons/kubernetes-csi-addons/controllers/replication.storage/predicate"
	"github.com/csi-addons/kubernetes-csi-addons/controllers/replication.storage/vgr"
	vrutils "github.com/csi-addons/kubernetes-csi-addons/controllers/utils"
	grpcClient "github.com/csi-addons/kubernetes-csi-addons/internal/client"
	conn "github.com/csi-addons/kubernetes-csi-addons/internal/connection"
	"github.com/go-logr/logr"
	"github.com/google/uuid"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// VolumeReplicationReconciler reconciles a VolumeReplication object.
type VolumeGroupReplicationReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	Log           logr.Logger
	Connpool      *conn.ConnectionPool
	Timeout       time.Duration
	Replication   grpcClient.VolumeReplication
	EventRecorder *vrutils.EventReporter
}

// +kubebuilder:rbac:groups=replication.storage.openshift.io,resources=volumegroupreplications,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=replication.storage.openshift.io,resources=volumegroupreplications/status,verbs=update
// +kubebuilder:rbac:groups=replication.storage.openshift.io,resources=volumegroupreplications/finalizers,verbs=update
// +kubebuilder:rbac:groups=replication.storage.openshift.io,resources=volumegroupreplicationcontents,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=replication.storage.openshift.io,resources=volumegroupreplicationcontents/status,verbs=update
// +kubebuilder:rbac:groups=replication.storage.openshift.io,resources=volumegroupreplicationcontents/finalizers,verbs=update
// +kubebuilder:rbac:groups=replication.storage.openshift.io,resources=volumegroupreplicationclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *VolumeGroupReplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("VolumeGroupReplicationGroup", req.NamespacedName, "rid", uuid.New())
	log.Info("Entering reconcile loop")
	defer log.Info("Exiting reconcile loop")

	v := vgr.VGRInstance{
		EventRecorder:  r.EventRecorder,
		Client:         r.Client,
		Ctx:            ctx,
		Log:            log,
		Instance:       &replicationv1alpha1.VolumeGroupReplication{},
		VolRepPVCs:     []corev1.PersistentVolumeClaim{},
		VolSyncPVCs:    []corev1.PersistentVolumeClaim{},
		ReplClassList:  &replicationv1alpha1.VolumeReplicationClassList{},
		NamespacedName: req.NamespacedName.String(),
	}

	if err := r.Client.Get(ctx, req.NamespacedName, v.Instance); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Resource not found")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get resource")
		return ctrl.Result{}, fmt.Errorf("failed to reconcile VolumeReplicationGroup (%v), %w",
			req.NamespacedName, err)
	}

	vgrClass, err := vgr.GetVGRClass(v)
	if err != nil {
		return ctrl.Result{}, err
	}
	//if err = utils.ValidatePrefixedParameters(vgClass.Parameters); err != nil {
	//	logger.Error(err, "failed to validate parameters of volumegroupClass", "VGClassName", vgClass.Name)
	//	if uErr := utils.UpdateVGStatusError(r.Client, instance, logger, err.Error()); uErr != nil {
	//		return ctrl.Result{}, uErr
	//	}
	//	return ctrl.Result{}, err
	//}
	v.Driver = vgrClass.Driver

	if v.Instance.GetDeletionTimestamp().IsZero() {
		//if err = utils.AddFinalizerToVG(r.Client, logger, instance); err != nil {
		//	return ctrl.Result{}, utils.HandleErrorMessage(logger, r.Client, instance, err, createVG)
		//}
	} else {
		//if commonUtils.Contains(instance.GetFinalizers(), utils.VGFinalizer) && !utils.IsContainOtherFinalizers(instance, logger) {
		if true {
			if err := r.removeInstance(v); err != nil {
				return ctrl.Result{}, err
			}
			log.Info("volumeGroup object is terminated, skipping reconciliation")
		}
		return ctrl.Result{}, nil
	}

	//groupCreationTime := utils.GetCurrentTime()

	//err, isStaticProvisioned := r.handleStaticProvisionedVG(instance, logger, groupCreationTime, vgClass)
	//if isStaticProvisioned {
	//	return ctrl.Result{}, err
	//}

	vgrcName, err := vgr.GenerateVGRCName(string(v.Instance.GetUID()))
	if err != nil {
		return ctrl.Result{}, err
	}
	vgrc := vgr.GenerateVGRC(v.Instance, vgrcName, vgrClass)

	if err = vgr.CreateVGRC(v, vgrc); err != nil {
		return ctrl.Result{}, err
	}
	if isVGCReady, err := vgr.IsVGRCReady(v, vgrc); err != nil {
		return ctrl.Result{}, err
	} else if !isVGCReady {
		return ctrl.Result{Requeue: true}, nil
	}

	err = r.updateItems(v, vgrc.Name)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.updateReplicatedPVCs(v)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *VolumeGroupReplicationReconciler) removeInstance(v vgr.VGRInstance) error {
	var vgrcName, vgrcNamespace string
	if *v.Instance.Spec.VolumeGroupReplicationContentName != "" {
		vgrcNamespace = v.Instance.Namespace
		vgrcName = *v.Instance.Spec.VolumeGroupReplicationContentName
	}
	vgrc, err := vgr.GetVGRC(v, vgrcName, vgrcNamespace)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}

	} else {
		err = r.removeVGRCObject(v, vgrc)
		if err != nil {
			return err
		}
	}
	//if err = vgr.RemoveFinalizerFromVGR(v); err != nil {
	//	return err
	//}
	return nil
}

func (r *VolumeGroupReplicationReconciler) removeVGRCObject(v vgr.VGRInstance,
	vgrc *replicationv1alpha1.VolumeGroupReplicationContent) error {
	if vgrc.Spec.VolumeGroupReplicationDeletionPolicy == replicationv1alpha1.VolumeGroupReplicationContentDelete {
		if err := v.Client.Delete(context.TODO(), vgrc); err != nil {
			v.Log.Error(err, fmt.Sprintf("Failed to delete (%v/%v) volumeGroupReplicationContent",
				vgrc.Namespace, vgrc.Name))
			return err
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VolumeGroupReplicationReconciler) SetupWithManager(mgr ctrl.Manager, ctrlOptions controller.Options) error {
	r.EventRecorder = vrutils.NewEventReporter(mgr.GetEventRecorderFor("controller_VolumeReplicationGroup"))

	generationPred := predicate.GenerationChangedPredicate{}
	pred := predicate.Or(generationPred, vrpredicate.FinalizerPredicate())

	return ctrl.NewControllerManagedBy(mgr).
		For(&replicationv1alpha1.VolumeGroupReplication{}, builder.WithPredicates(pred)).
		Watches(&source.Kind{Type: &corev1.PersistentVolumeClaim{}}, vrpredicate.CreateRequests(r.Client), builder.WithPredicates(vrpredicate.PVCPredicateFunc())).
		Complete(r)
}

func (r *VolumeGroupReplicationReconciler) updateItems(v vgr.VGRInstance, vgcName string) error {
	if err := vgr.UpdateVGRSourceContent(v, vgcName); err != nil {
		return err
	}
	if err := vgr.UpdateVGRStatus(v, vgcName, true); err != nil {
		return err
	}
	return nil
}

func (r *VolumeGroupReplicationReconciler) updateReplicatedPVCs(v vgr.VGRInstance) error {
	vgrObj := v.Instance
	matchingPvcs, err := vgr.GetMatchingPVCs(v)
	if err != nil {
		return err
	}
	if vgr.IsReplicatedPVCsEqual(matchingPvcs, vgrObj.Status.ReplicatedPVCs) {
		return nil
	}
	err = r.ModifyVolumesInVGR(v, r.Replication, matchingPvcs)
	if err != nil {
		return err
	}
	err = vgr.UpdateReplicatedPVCsAndPVs(v, matchingPvcs)
	if err != nil {
		return err
	}
	return nil
}

func (r *VolumeGroupReplicationReconciler) ModifyVolumesInVGR(v vgr.VGRInstance,
	replicationClient grpcClient.VolumeReplication, matchingPvcs []replicationv1alpha1.ReplicatedPVC) error {
	vgrObj := v.Instance
	currentList := make([]replicationv1alpha1.ReplicatedPVC, len(vgrObj.Status.ReplicatedPVCs))
	copy(currentList, vgrObj.Status.ReplicatedPVCs)
	vgrObj.Status.ReplicatedPVCs = matchingPvcs

	err := vgr.ModifyVGR(v, replicationClient)
	if err != nil {
		vgrObj.Status.ReplicatedPVCs = currentList
		return err
	}
	return nil
}
