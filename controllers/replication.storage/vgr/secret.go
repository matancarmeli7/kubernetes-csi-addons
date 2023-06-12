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

	"github.com/csi-addons/kubernetes-csi-addons/internal/util"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getSecretData(client client.Client, logger logr.Logger, name, namespace string) (map[string]string, error) {
	namespacedName := types.NamespacedName{Name: name, Namespace: namespace}
	secret := &corev1.Secret{}
	err := client.Get(context.TODO(), namespacedName, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Error(err, "secret not found", "Secret Name", name, "Secret Namespace", namespace)

			return nil, err
		}
		logger.Error(err, "error getting secret", "Secret Name", name, "Secret Namespace", namespace)

		return nil, err
	}

	return util.ConvertMap(secret.Data), nil
}
