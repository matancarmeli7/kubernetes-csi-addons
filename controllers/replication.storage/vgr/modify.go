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
	"fmt"

	"github.com/csi-addons/kubernetes-csi-addons/controllers/replication.storage/replication"
	grpcClient "github.com/csi-addons/kubernetes-csi-addons/internal/client"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ModifyVGR(vgr VGRInstance, replicationClient grpcClient.VolumeReplication) error {
	vgrObj := vgr.Instance
	params, err := generateModifyVGParams(vgr, replicationClient)
	if err != nil {
		return err
	}
	vgr.Log.Info(fmt.Sprintf("Modifying %s volumeGroupReplicationID with %v volumeIDs",
		params.VolumeGroupID, params.VolumeIds))
	volumeGroupRequest := replication.Replication{
		Params: params,
	}
	modifyVGResponse := volumeGroupRequest.Modify()
	responseError := modifyVGResponse.Error
	if responseError != nil {
		vgr.Log.Error(responseError, fmt.Sprintf("Failed to modify %s/%s volumeGroupReplication",
			vgrObj.Namespace, vgrObj.Name))
		return responseError
	}
	vgr.Log.Info(fmt.Sprintf("Successfully modified %s volumeGroupReplicationID", params.VolumeGroupID))
	return nil
}

func generateModifyVGParams(vgr VGRInstance, replicationClient grpcClient.VolumeReplication,
) (replication.CommonRequestParameters, error) {
	vgrId, err := getVGRId(vgr)
	if err != nil {
		return replication.CommonRequestParameters{}, err
	}
	volumeIds, err := getVolumeIds(vgr.Log, vgr.Client, vgr.Instance.Status.ReplicatedPVCs)
	if err != nil {
		return replication.CommonRequestParameters{}, err
	}
	secrets, err := getSecrets(vgr.Log, vgr.Client, vgr)
	if err != nil {
		return replication.CommonRequestParameters{}, err
	}

	return replication.CommonRequestParameters{
		Secrets:       secrets,
		VolumeGroup:   vgClient,
		VolumeGroupID: vgrId,
		VolumeIds:     volumeIds,
	}, nil
}

func getSecrets(logger logr.Logger, client client.Client, vgr VGRInstance) (map[string]string, error) {
	vgc, err := GetVGRClass(vgr)
	if err != nil {
		return nil, err
	}
	secrets, err := GetSecretDataFromClass(client, vgc, logger)
	if err != nil {
		//if uErr := UpdateVGStatusError(client, vg, logger, err.Error()); uErr != nil {
		//	return nil, err
		//}
		return nil, err
	}
	return secrets, nil
}
