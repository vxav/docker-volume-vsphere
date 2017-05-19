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

// This test suite includes test cases to verify basic functionality
// in most common configurations

package e2e

import (
	"log"
	"os"

	. "gopkg.in/check.v1"

	"github.com/vmware/docker-volume-vsphere/tests/utils/dockercli"
	"github.com/vmware/docker-volume-vsphere/tests/utils/inputparams"
	"github.com/vmware/docker-volume-vsphere/tests/utils/misc"
	"github.com/vmware/docker-volume-vsphere/tests/utils/verification"
)

type BasicTestSuite struct {
	esxName       string
	volumeName    string
	containerName string
}

func (s *BasicTestSuite) SetUpTest(c *C) {
	s.esxName = os.Getenv("ESX")
	s.volumeName = inputparams.GetVolumeNameWithTimeStamp()
	s.containerName = inputparams.GetContainerNameWithTimeStamp()
}

var _ = Suite(&BasicTestSuite{})

// Test volume lifecycle management on different datastores:
// VM1 - local VMFS datastore
// VM2 - shared VMFS datastore
// VM3 - shared VSAN datastore
func (s *BasicTestSuite) TestVolumeCreation(c *C) {
	log.Printf("START: basic-test.TestVolumeCreation")

	dockerHosts := []string{os.Getenv("VM1"), os.Getenv("VM2"), os.Getenv("VM3")}
	for _, host := range dockerHosts {
		//TODO: Remove this check once VM3 is available from the CI testbed
		if host == "" {
			continue
		}

		// Create a volume
		out, err := dockercli.CreateVolume(host, s.volumeName)
		c.Assert(err, IsNil, Commentf(misc.FormatOutput(out)))

		// Attach the volume
		out, err = dockercli.AttachVolume(host, s.volumeName, s.containerName)
		c.Assert(err, IsNil, Commentf(misc.FormatOutput(out)))

		// Verify volume status
		status := verification.VerifyAttachedStatus(s.volumeName, host, s.esxName)
		c.Assert(status, Equals, true, Commentf("Volume %s is not attached", s.volumeName))

		// Clean up the container
		out, err = dockercli.RemoveContainer(host, s.containerName)
		c.Assert(err, IsNil, Commentf(misc.FormatOutput(out)))

		// Verify volume status
		status = verification.VerifyDetachedStatus(s.volumeName, host, s.esxName)
		c.Assert(status, Equals, true, Commentf("Volume %s is still attached", s.volumeName))

		// Clean up the volume
		out, err = dockercli.DeleteVolume(host, s.volumeName)
		c.Assert(err, IsNil, Commentf(misc.FormatOutput(out)))
	}

	log.Printf("END: basic-test.TestVolumeCreation")
}

// Test volume visibility: volume is created on the local datastore attached to VM1:
// VM1 - accessible to the volume
// VM2 - inaccessible to the volume
func (s *BasicTestSuite) TestVolumeAccessibility(c *C) {
	log.Printf("START: basic-test.TestVolumeAccessibility")

	vm1, vm2 := os.Getenv("VM1"), os.Getenv("VM2")

	// Create a volume from VM1
	out, err := dockercli.CreateVolume(vm1, s.volumeName)
	c.Assert(err, IsNil, Commentf(misc.FormatOutput(out)))

	// Verify the volume is accessible from VM1
	accessible := verification.CheckVolumeAccessibility(vm1, s.volumeName)
	c.Assert(accessible, Equals, true, Commentf("Volume %s is not accessible", s.volumeName))

	// Verify the volume is inaccessible from VM4
	accessible = verification.CheckVolumeAccessibility(vm2, s.volumeName)
	//TODO: VM2 inaccessible to this volume is currently not available
	//c.Assert(accessible, Equals, false, Commentf("Volume %s is accessible", s.volumeName))

	// Clean up the volume
	out, err = dockercli.DeleteVolume(vm1, s.volumeName)
	c.Assert(err, IsNil, Commentf(misc.FormatOutput(out)))

	log.Printf("END: basic-test.TestVolumeAccessibility")
}
