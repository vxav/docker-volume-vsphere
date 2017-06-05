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

// This test is going to create volume on the fresh testbed very first time.
// After installing vmdk volume plugin/driver, volume creation should not be
// failed very first time.

// This test is going to cover various vsan related test cases

package e2e

import (
	"log"
	"os"
	"strings"

	"github.com/vmware/docker-volume-vsphere/tests/utils/admincli"
	"github.com/vmware/docker-volume-vsphere/tests/utils/dockercli"
	"github.com/vmware/docker-volume-vsphere/tests/utils/govc"
	"github.com/vmware/docker-volume-vsphere/tests/utils/inputparams"
	"github.com/vmware/docker-volume-vsphere/tests/utils/verification"

	. "gopkg.in/check.v1"
)

const vsanPolicyFlag = "vsan-policy-name"

type VsanTestSuite struct {
	hostIP     string
	esxIP      string
	vsanDSName string
	volumeName string
}

func (s *VsanTestSuite) SetUpSuite(c *C) {
	s.hostIP = os.Getenv("VM2")
	s.esxIP = os.Getenv("ESX")
	s.vsanDSName = govc.GetDatastoreByType("vsan")
}

func (s *VsanTestSuite) TearDownTest(c *C) {
	out, err := dockercli.DeleteVolume(s.hostIP, s.volumeName)
	c.Assert(err, IsNil, Commentf(out))
}

var _ = Suite(&VsanTestSuite{})

// Vsan related test
// 1. Create a valid vsan policy
// 2. Create an invalid vsan policy (wrong content)
// 3. Volume creation with valid policy should pass
// 4. Valid volume should be accessible
// 5. Volume creation with non existing policy should fail
// 6. Volume creation with invalid policy should fail
func (s *VsanTestSuite) TestVSANPolicy(c *C) {
	if s.vsanDSName == "" {
		c.Skip("Vsan datastore unavailable")
	}

	log.Printf("START: vsan_test.TestVSANPolicy")

	policyName := "validPolicy"
	out, err := admincli.CreatePolicy(s.esxIP, policyName, "'((\"proportionalCapacity\" i50)''(\"hostFailuresToTolerate\" i0))'")
	c.Assert(err, IsNil, Commentf(out))

	invalidContentPolicyName := "invalidPolicy"
	out, err = admincli.CreatePolicy(s.esxIP, invalidContentPolicyName, "'((\"wrongKey\" i50)'")
	c.Assert(err, IsNil, Commentf(out))

	s.volumeName = inputparams.GetVolumeNameWithTimeStamp("vsanVol") + "@" + s.vsanDSName
	vsanOpts := " -o " + vsanPolicyFlag + "=" + policyName

	out, err = dockercli.CreateVolumeWithOptions(s.hostIP, s.volumeName, vsanOpts)
	c.Assert(err, IsNil, Commentf(out))
	isAvailable := verification.CheckVolumeAvailability(s.hostIP, s.volumeName)
	c.Assert(isAvailable, Equals, true, Commentf("Volume %s is not available after creation", s.volumeName))

	invalidVsanOpts := [2]string{"-o " + vsanPolicyFlag + "=IDontExist", "-o " + vsanPolicyFlag + "=" + invalidContentPolicyName}
	for _, option := range invalidVsanOpts {
		invalidVolName := inputparams.GetVolumeNameWithTimeStamp("vsanVol") + "@" + s.vsanDSName
		out, _ = dockercli.CreateVolumeWithOptions(s.hostIP, invalidVolName, option)
		c.Assert(strings.HasPrefix(out, ErrorVolumeCreate), Equals, true)
	}

	log.Printf("END: vsan_test.TestVSANPolicy")
}
