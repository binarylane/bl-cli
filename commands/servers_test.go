/*
Copyright 2018 The Doctl Authors All rights reserved.
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

package commands

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/go-binarylane"
	"github.com/stretchr/testify/assert"
)

var (
	testImage = bl.Image{Image: &binarylane.Image{
		ID:      1,
		Slug:    "slug",
		Regions: []string{"test0"},
	}}
	testImageSecondary = bl.Image{Image: &binarylane.Image{
		ID:      2,
		Slug:    "slug-secondary",
		Regions: []string{"test0"},
	}}
	testImageList = bl.Images{testImage, testImageSecondary}
)

func TestServerCommand(t *testing.T) {
	cmd := Server()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "actions", "backups", "create", "delete", "get", "kernels", "list", "neighbors", "snapshots", "tag", "untag")
}

func TestServerActionList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Actions(1).Return(testActionList, nil)

		config.Args = append(config.Args, "1")

		err := RunServerActions(config)
		assert.NoError(t, err)
	})
}

func TestServerBackupList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Backups(1).Return(testImageList, nil)

		config.Args = append(config.Args, "1")

		err := RunServerBackups(config)
		assert.NoError(t, err)
	})
}

func TestServerCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		volumeUUID := "00000000-0000-4000-8000-000000000000"
		vpcID := 2
		dcr := &binarylane.ServerCreateRequest{
			Name:   "server",
			Region: "dev0",
			Size:   "1gb",
			Image:  binarylane.ServerCreateImage{ID: 0, Slug: "image"},
			Volumes: []binarylane.ServerCreateVolume{
				{Name: "test-volume"},
				{ID: volumeUUID},
			},
			Backups:           false,
			IPv6:              false,
			PrivateNetworking: false,
			Monitoring:        false,
			VPCID:             vpcID,
			UserData:          "#cloud-config",
			Tags:              []string{"one", "two"},
		}
		tm.servers.EXPECT().Create(dcr, false).Return(&testServer, nil)

		config.Args = append(config.Args, "server")

		config.Doit.Set(config.NS, blcli.ArgRegionSlug, "dev0")
		config.Doit.Set(config.NS, blcli.ArgSizeSlug, "1gb")
		config.Doit.Set(config.NS, blcli.ArgImage, "image")
		config.Doit.Set(config.NS, blcli.ArgUserData, "#cloud-config")
		config.Doit.Set(config.NS, blcli.ArgVPCID, vpcID)
		config.Doit.Set(config.NS, blcli.ArgVolumeList, []string{"test-volume", volumeUUID})
		config.Doit.Set(config.NS, blcli.ArgTagNames, []string{"one", "two"})

		err := RunServerCreate(config)
		assert.NoError(t, err)
	})
}

func TestServerCreateWithNoneSSHKey(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		volumeUUID := "00000000-0000-4000-8000-000000000000"
		vpcID := 2
		dcr := &binarylane.ServerCreateRequest{
			Name:    "server",
			Region:  "dev0",
			Size:    "1gb",
			Image:   binarylane.ServerCreateImage{ID: 0, Slug: "image"},
			SSHKeys: []binarylane.ServerCreateSSHKey{},
			Volumes: []binarylane.ServerCreateVolume{
				{Name: "test-volume"},
				{ID: volumeUUID},
			},
			Backups:           false,
			IPv6:              false,
			PrivateNetworking: false,
			Monitoring:        false,
			VPCID:             vpcID,
			UserData:          "#cloud-config",
			Tags:              []string{"one", "two"},
		}
		tm.servers.EXPECT().Create(dcr, false).Return(&testServer, nil)

		config.Args = append(config.Args, "server")

		config.Doit.Set(config.NS, blcli.ArgRegionSlug, "dev0")
		config.Doit.Set(config.NS, blcli.ArgSizeSlug, "1gb")
		config.Doit.Set(config.NS, blcli.ArgImage, "image")
		config.Doit.Set(config.NS, blcli.ArgUserData, "#cloud-config")
		config.Doit.Set(config.NS, blcli.ArgVPCID, vpcID)
		config.Doit.Set(config.NS, blcli.ArgSSHKeys, []string{"none"})
		config.Doit.Set(config.NS, blcli.ArgVolumeList, []string{"test-volume", volumeUUID})
		config.Doit.Set(config.NS, blcli.ArgTagNames, []string{"one", "two"})

		err := RunServerCreate(config)
		assert.NoError(t, err)
	})
}

func TestServerCreateWithTag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		dcr := &binarylane.ServerCreateRequest{
			Name:              "server",
			Region:            "dev0",
			Size:              "1gb",
			Image:             binarylane.ServerCreateImage{ID: 0, Slug: "image"},
			Backups:           false,
			IPv6:              false,
			PrivateNetworking: false,
			UserData:          "#cloud-config",
			Tags:              []string{"my-tag"}}
		tm.servers.EXPECT().Create(dcr, false).Return(&testServer, nil)

		config.Args = append(config.Args, "server")

		config.Doit.Set(config.NS, blcli.ArgRegionSlug, "dev0")
		config.Doit.Set(config.NS, blcli.ArgSizeSlug, "1gb")
		config.Doit.Set(config.NS, blcli.ArgImage, "image")
		config.Doit.Set(config.NS, blcli.ArgUserData, "#cloud-config")
		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")

		err := RunServerCreate(config)
		assert.NoError(t, err)
	})
}

func TestServerCreateUserDataFile(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		userData := `
coreos:
  etcd2:
    discovery: https://discovery.etcd.io/<token>
    advertise-client-urls: http://$private_ipv4:2379,http://$private_ipv4:4001
    initial-advertise-peer-urls: http://$private_ipv4:2380
    listen-client-urls: http://0.0.0.0:2379,http://0.0.0.0:4001
    listen-peer-urls: http://$private_ipv4:2380
  units:
    - name: etcd2.service
      command: start
    - name: fleet.service
      command: start
`

		tmpFile, err := ioutil.TempFile(os.TempDir(), "blServersTest-*.yml")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(userData)
		assert.NoError(t, err)

		err = tmpFile.Close()
		assert.NoError(t, err)

		dcr := &binarylane.ServerCreateRequest{
			Name:   "server",
			Region: "dev0",
			Size:   "1gb",
			Image: binarylane.ServerCreateImage{
				ID:   0,
				Slug: "image",
			},
			Backups:           false,
			IPv6:              false,
			PrivateNetworking: false,
			UserData:          userData,
		}
		tm.servers.EXPECT().Create(dcr, false).Return(&testServer, nil)

		config.Args = append(config.Args, "server")

		config.Doit.Set(config.NS, blcli.ArgRegionSlug, "dev0")
		config.Doit.Set(config.NS, blcli.ArgSizeSlug, "1gb")
		config.Doit.Set(config.NS, blcli.ArgImage, "image")
		config.Doit.Set(config.NS, blcli.ArgUserDataFile, tmpFile.Name())

		err = RunServerCreate(config)
		assert.NoError(t, err)
	})
}

func TestServerDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Delete(1).Return(nil)

		config.Args = append(config.Args, strconv.Itoa(testServer.ID))
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunServerDelete(config)
		assert.NoError(t, err)

	})
}

func TestServerDeleteByTag_ServersExist(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().ListByTag("my-tag").Return(testServerList, nil)
		tm.servers.EXPECT().DeleteByTag("my-tag").Return(nil)

		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunServerDelete(config)
		assert.NoError(t, err)
	})
}

func TestServerDeleteByTag_ServersMissing(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().ListByTag("my-tag").Return(bl.Servers{}, nil)

		var buf bytes.Buffer
		config.Out = &buf
		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunServerDelete(config)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Nothing to delete")
	})
}

func TestServerDeleteRepeatedID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Delete(1).Return(nil).Times(1)

		id := strconv.Itoa(testServer.ID)
		config.Args = append(config.Args, id, id)
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunServerDelete(config)
		assert.NoError(t, err)
	})
}

func TestServerDeleteByName(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().List().Return(testServerList, nil)
		tm.servers.EXPECT().Delete(1).Return(nil)

		config.Args = append(config.Args, testServer.Name)
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunServerDelete(config)
		assert.NoError(t, err)
	})
}

func TestServerDeleteByName_Ambiguous(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		list := bl.Servers{testServer, testServer}
		tm.servers.EXPECT().List().Return(list, nil)

		config.Args = append(config.Args, testServer.Name)
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunServerDelete(config)
		t.Log(err)
		assert.Error(t, err)
	})
}

func TestServerDelete_MixedNameAndType(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().List().Return(testServerList, nil)
		tm.servers.EXPECT().Delete(1).Return(nil).Times(1)

		id := strconv.Itoa(testServer.ID)
		config.Args = append(config.Args, id, testServer.Name)
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunServerDelete(config)
		assert.NoError(t, err)
	})

}

func TestServerGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Get(testServer.ID).Return(&testServer, nil)

		config.Args = append(config.Args, strconv.Itoa(testServer.ID))

		err := RunServerGet(config)
		assert.NoError(t, err)
	})
}

func TestServerGet_Template(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Get(testServer.ID).Return(&testServer, nil)

		config.Args = append(config.Args, strconv.Itoa(testServer.ID))
		config.Doit.Set(config.NS, blcli.ArgTemplate, "{{.Name}}")

		err := RunServerGet(config)
		assert.NoError(t, err)
	})
}

func TestServerKernelList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Kernels(testServer.ID).Return(testKernelList, nil)

		config.Args = append(config.Args, "1")

		err := RunServerKernels(config)
		assert.NoError(t, err)
	})
}

func TestServerNeighbors(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Neighbors(testServer.ID).Return(testServerList, nil)

		config.Args = append(config.Args, "1")

		err := RunServerNeighbors(config)
		assert.NoError(t, err)
	})
}

func TestServerSnapshotList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().Snapshots(testServer.ID).Return(testImageList, nil)

		config.Args = append(config.Args, "1")

		err := RunServerSnapshots(config)
		assert.NoError(t, err)
	})
}

func TestServersList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().List().Return(testServerList, nil)

		err := RunServerList(config)
		assert.NoError(t, err)
	})
}

func TestServersListByTag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.servers.EXPECT().ListByTag("my-tag").Return(testServerList, nil)

		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")

		err := RunServerList(config)
		assert.NoError(t, err)
	})
}

func TestServersTag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &binarylane.TagResourcesRequest{
			Resources: []binarylane.Resource{
				{ID: "1", Type: binarylane.ServerResourceType},
			},
		}
		tm.tags.EXPECT().TagResources("my-tag", trr).Return(nil)

		config.Args = append(config.Args, "1")
		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")

		err := RunServerTag(config)
		assert.NoError(t, err)
	})
}

func TestServersTagMultiple(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &binarylane.TagResourcesRequest{
			Resources: []binarylane.Resource{
				{ID: "1", Type: binarylane.ServerResourceType},
				{ID: "2", Type: binarylane.ServerResourceType},
			},
		}
		tm.tags.EXPECT().TagResources("my-tag", trr).Return(nil)

		config.Args = append(config.Args, "1")
		config.Args = append(config.Args, "2")
		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")

		err := RunServerTag(config)
		assert.NoError(t, err)
	})
}

func TestServersTagByName(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &binarylane.TagResourcesRequest{
			Resources: []binarylane.Resource{
				{ID: "1", Type: binarylane.ServerResourceType},
			},
		}
		tm.tags.EXPECT().TagResources("my-tag", trr).Return(nil)
		tm.servers.EXPECT().List().Return(testServerList, nil)

		config.Args = append(config.Args, testServer.Name)
		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")

		err := RunServerTag(config)
		assert.NoError(t, err)
	})
}

func TestServersTagMultipleNameAndID(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		trr := &binarylane.TagResourcesRequest{
			Resources: []binarylane.Resource{
				{ID: "1", Type: binarylane.ServerResourceType},
				{ID: "3", Type: binarylane.ServerResourceType},
			},
		}
		tm.tags.EXPECT().TagResources("my-tag", trr).Return(nil)
		tm.servers.EXPECT().List().Return(testServerList, nil)

		config.Args = append(config.Args, testServer.Name)
		config.Args = append(config.Args, strconv.Itoa(anotherTestServer.ID))
		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")

		err := RunServerTag(config)
		assert.NoError(t, err)
	})
}

func TestServersUntag(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		urr := &binarylane.UntagResourcesRequest{
			Resources: []binarylane.Resource{
				{ID: "1", Type: binarylane.ServerResourceType},
			},
		}

		tm.tags.EXPECT().UntagResources("my-tag", urr).Return(nil)
		tm.servers.EXPECT().List().Return(testServerList, nil)

		config.Args = []string{testServer.Name}
		config.Doit.Set(config.NS, blcli.ArgTagName, "my-tag")

		err := RunServerUntag(config)
		assert.NoError(t, err)
	})
}

func Test_extractSSHKey(t *testing.T) {
	cases := []struct {
		in       []string
		expected []binarylane.ServerCreateSSHKey
	}{
		{
			in:       []string{"1"},
			expected: []binarylane.ServerCreateSSHKey{{ID: 1}},
		},
		{
			in:       []string{"fingerprint"},
			expected: []binarylane.ServerCreateSSHKey{{Fingerprint: "fingerprint"}},
		},
		{
			in:       []string{"1", "2"},
			expected: []binarylane.ServerCreateSSHKey{{ID: 1}, {ID: 2}},
		},
		{
			in:       []string{"1", "fingerprint"},
			expected: []binarylane.ServerCreateSSHKey{{ID: 1}, {Fingerprint: "fingerprint"}},
		},
		{
			in:       []string{"none"},
			expected: []binarylane.ServerCreateSSHKey{},
		},
		{
			in:       []string{},
			expected: nil,
		},
	}

	for _, c := range cases {
		got := extractSSHKeys(c.in)
		assert.Equal(t, c.expected, got)
	}
}
