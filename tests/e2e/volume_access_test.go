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

// The goal of this test suite is to verify read/write consistency on volumes
// in accordance with the access updates on the volume
package e2e

import (
	"log"
	"os"
	"strings"

	. "gopkg.in/check.v1"

	adminconst "github.com/vmware/docker-volume-vsphere/tests/constants/admincli"
	"github.com/vmware/docker-volume-vsphere/tests/utils/admincli"
	"github.com/vmware/docker-volume-vsphere/tests/utils/dockercli"
	"github.com/vmware/docker-volume-vsphere/tests/utils/govc"
	"github.com/vmware/docker-volume-vsphere/tests/utils/inputparams"
)

const errorWriteVolume = "Read-only file system"

type VolumeAccessTestSuite struct {
	volumeName    string
	dockerHostIP  []string
	containerList [2][]string
	esxIP         string
	dsList        []string
}

func (s *VolumeAccessTestSuite) SetUpSuite(c *C) {
	s.dockerHostIP = []string{os.Getenv("VM1"), os.Getenv("VM2")}
	s.esxIP = os.Getenv("ESX")

	s.dsList = govc.GetDatastoreList()
	c.Assert(len(s.dsList), Not(Equals), 0, Commentf("Datastore list empty"))
}

func (s *VolumeAccessTestSuite) SetUpTest(c *C) {
	dsName := s.dsList[0]
	s.volumeName = inputparams.GetVolumeNameWithTimeStamp("vol_access") + "@" + dsName
}

func (s *VolumeAccessTestSuite) TearDownTest(c *C) {
	for i := 0; i < 2; i++ {
		cnameList := strings.Join(s.containerList[i], " ")
		out, err := dockercli.RemoveContainer(s.dockerHostIP[i], cnameList)
		c.Assert(err, IsNil, Commentf(out))
		s.containerList[i] = s.containerList[i][:0]
	}

	out, err := dockercli.DeleteVolume(s.dockerHostIP[0], s.volumeName)
	c.Assert(err, IsNil, Commentf(out))
}

var _ = Suite(&VolumeAccessTestSuite{})

func (s *VolumeAccessTestSuite) newCName(i int) string {
	cname := inputparams.GetContainerNameWithTimeStamp("vol_access")
	s.containerList[i] = append(s.containerList[i], cname)
	return cname
}

// Verify read, write is possible after volume access update
// 1. Write a message from host1 to a file on the volume
// 2. Read the content from host2 from same file on the volume
//    Verify the content is same.
// 3. Write another message from host2 to the same file on that volume
// 4. Read from host1 and verify the content is same
// 5. Update the volume access to read-only
// 6. Write from host1 to the file on volume should fail
// 7. Write from host2 should also fail
// 8. Update the volume access to read-write
// 9. Write from host1 should succeed
// 10. Read from host2 to verify the content is same
// 11. Write from host2 should succeed
// 12. Read from host1 to verify the content is same
func (s *VolumeAccessTestSuite) TestAccessUpdate(c *C) {
	log.Printf("START: volume_access_test.TestAccessUpdate")

	data1 := "message_by_host1"
	data2 := "message_by_host2"
	testFile := "test.txt"

	// Create a volume
	out, err := dockercli.CreateVolume(s.dockerHostIP[0], s.volumeName)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.WriteToVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile, data1)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.ReadFromVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile)
	c.Assert(err, IsNil, Commentf(out))
	c.Assert(out, Equals, data1)

	out, err = dockercli.WriteToVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile, data2)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.ReadFromVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile)
	c.Assert(err, IsNil, Commentf(out))
	c.Assert(out, Equals, data2)

	out, err = admincli.UpdateVolumeAccess(s.esxIP, s.volumeName, adminconst.DefaultVMgroup, adminconst.ReadOnlyAccess)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.WriteToVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile, data1)
	c.Assert(strings.Contains(out, errorWriteVolume), Equals, true, Commentf(out))

	out, err = dockercli.WriteToVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile, data2)
	c.Assert(strings.Contains(out, errorWriteVolume), Equals, true, Commentf(out))

	out, err = admincli.UpdateVolumeAccess(s.esxIP, s.volumeName, adminconst.DefaultVMgroup, adminconst.ReadWriteAccess)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.WriteToVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile, data1)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.ReadFromVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile)
	c.Assert(err, IsNil, Commentf(out))
	c.Assert(out, Equals, data1)

	out, err = dockercli.WriteToVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile, data2)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.ReadFromVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile)
	c.Assert(err, IsNil, Commentf(out))
	c.Assert(out, Equals, data2)

	log.Printf("END: volume_access_test.TestAccessUpdate")
}

// Verify read, write is possible after volume access update
// 1. Create a volume with read-only access
// 2. Write from host1 to the file on volume should fail
// 3. Write from host2 should also fail
// 4. Update the volume access to read-write
// 5. Write from host1 should succeed
// 6. Read from host2 to verify the content is same
// 7. Write from host2 should succeed
// 8. Read from host1 to verify the content is same
func (s *VolumeAccessTestSuite) TestAccessUpdate_R_RW(c *C) {
	log.Printf("START: volume_access_test.TestAccessUpdate_R_RW")

	data1 := "message_by_host1"
	data2 := "message_by_host2"
	testFile := "test.txt"

	// Create a volume
	out, err := dockercli.CreateVolumeWithOptions(s.dockerHostIP[0], s.volumeName, " -o access="+adminconst.ReadOnlyAccess)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.WriteToVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile, data1)
	c.Assert(strings.Contains(out, errorWriteVolume), Equals, true, Commentf(out))

	out, err = dockercli.WriteToVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile, data2)
	c.Assert(strings.Contains(out, errorWriteVolume), Equals, true, Commentf(out))

	out, err = admincli.UpdateVolumeAccess(s.esxIP, s.volumeName, adminconst.DefaultVMgroup, adminconst.ReadWriteAccess)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.WriteToVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile, data1)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.ReadFromVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile)
	c.Assert(err, IsNil, Commentf(out))
	c.Assert(out, Equals, data1)

	out, err = dockercli.WriteToVolume(s.dockerHostIP[1], s.volumeName, s.newCName(1), testFile, data2)
	c.Assert(err, IsNil, Commentf(out))

	out, err = dockercli.ReadFromVolume(s.dockerHostIP[0], s.volumeName, s.newCName(0), testFile)
	c.Assert(err, IsNil, Commentf(out))
	c.Assert(out, Equals, data2)

	log.Printf("END: volume_access_test.TestAccessUpdate_R_RW")
}
