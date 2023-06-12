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
package util

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateObject(client client.Client, updateObject client.Object) error {
	if err := client.Update(context.TODO(), updateObject); err != nil {
		return fmt.Errorf("failed to update %s (%s/%s) %w", updateObject.GetObjectKind(),
			updateObject.GetNamespace(), updateObject.GetName(), err)
	}
	return nil
}

func UpdateObjectStatus(client client.Client, updateObject client.Object) error {
	if err := client.Status().Update(context.TODO(), updateObject); err != nil {
		if apierrors.IsConflict(err) {
			return err
		}
		return fmt.Errorf("failed to update %s (%s/%s) status %w", updateObject.GetObjectKind(),
			updateObject.GetNamespace(), updateObject.GetName(), err)
	}
	return nil
}

func GetNamespacedObject(client client.Client, obj client.Object) error {
	namespacedObject := types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}
	err := client.Get(context.TODO(), namespacedObject, obj)
	if err != nil {
		return err
	}
	return nil
}
