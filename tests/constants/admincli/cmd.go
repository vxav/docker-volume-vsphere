// Copyright 2017 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A home to hold test constants related with vmdkops_admin cli.

package admincli

const (
	// location of the vmdkops binary
	vmdkopsAdmin = "/usr/lib/vmware/vmdkops/bin/vmdkops_admin.py "

	// vmdkops_admin volume
	vmdkopsAdminVolume = vmdkopsAdmin + "volume "

	// ListVolumes referring to vmdkops_admin volume ls
	ListVolumes = vmdkopsAdminVolume + "ls "

	// CreatePolicy Create a policy
	CreatePolicy = vmdkopsAdmin + " policy create "

	// SetVolumeAccess set volume access
	SetVolumeAccess = vmdkopsAdminVolume + " set "

	// CreateVMgroup referring to create vmgroup
	// where --name will be name of the vmgroup
	CreateVMgroup = vmdkopsAdmin + "vmgroup create --name="

	// RemoveVMgroup referring to remove vmgroup
	// where --name will be name of the vmgroup
	RemoveVMgroup = vmdkopsAdmin + "vmgroup rm --name="

	// AddVMToVMgroup referring to add vm to a vmgroup
	// where --name will be name of the vmgroup
	AddVMToVMgroup = vmdkopsAdmin + "vmgroup vm add --name="

	// RemoveVMFromVMgroup referring to remove a vm from vmgroup
	// where --name will be name of the vmgroup
	RemoveVMFromVMgroup = vmdkopsAdmin + "vmgroup vm rm --name="

	// ReplaceVMFromVMgroup referring replace a vm from vmgroup
	// where --name will be name of the vmgroup
	ReplaceVMFromVMgroup = vmdkopsAdmin + "vmgroup vm replace --name="

	// DefaultVMgroup referring name of default vmgroup
	DefaultVMgroup = "_DEFAULT "

	// VMHomeDatastore referring datastore where the docker host vm is created
	VMHomeDatastore = "_VM_DS"

	// InitLocalConfigDb referring to Initialize (local) Single Node Config DB
	InitLocalConfigDb = vmdkopsAdmin + "config init --local"

	// RemoveLocalConfigDb referring to Remove (local) Single Node Config DB
	RemoveLocalConfigDb = vmdkopsAdmin + "config rm --local --confirm"

	// ListVMgroups referring to vmdkops_admin vmgroups ls
	ListVMgroups = vmdkopsAdmin + "vmgroup ls "

	// UpdateVMgroup referring to updating a vmgroup
	// where --name will be name of the vmgroup
	UpdateVMgroup = vmdkopsAdmin + "vmgroup update --name="

	//RemoveVolumes option refers to removing all volumes from a vmgroup
	RemoveVolumes = " --remove-volumes"

	// ReadOnlyAccess read only rights for the volume
	ReadOnlyAccess = "read-only"

	// ReadWriteAccess read-write rights for the volume
	ReadWriteAccess = "read-write"
)
