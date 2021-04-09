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
	"io/ioutil"
	"sort"
	"testing"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	blmocks "github.com/binarylane/bl-cli/bl/mocks"
	"github.com/binarylane/go-binarylane"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testServer = bl.Server{
		Server: &binarylane.Server{
			ID: 1,
			Image: &binarylane.Image{
				ID:           1,
				Name:         "an-image",
				Distribution: "DOOS",
			},
			Name: "a-server",
			Networks: &binarylane.Networks{
				V4: []binarylane.NetworkV4{
					{IPAddress: "8.8.8.8", Type: "public"},
					{IPAddress: "172.16.1.2", Type: "private"},
				},
			},
			Region: &binarylane.Region{
				Slug: "test0",
				Name: "test 0",
			},
		},
	}

	anotherTestServer = bl.Server{
		Server: &binarylane.Server{
			ID: 3,
			Image: &binarylane.Image{
				ID:           1,
				Name:         "an-image",
				Distribution: "DOOS",
			},
			Name: "another-server",
			Networks: &binarylane.Networks{
				V4: []binarylane.NetworkV4{
					{IPAddress: "8.8.8.9", Type: "public"},
					{IPAddress: "172.16.1.4", Type: "private"},
				},
			},
			Region: &binarylane.Region{
				Slug: "test0",
				Name: "test 0",
			},
		},
	}

	testPrivateServer = bl.Server{
		Server: &binarylane.Server{
			ID: 1,
			Image: &binarylane.Image{
				ID:           1,
				Name:         "an-image",
				Distribution: "DOOS",
			},
			Name: "a-server",
			Networks: &binarylane.Networks{
				V4: []binarylane.NetworkV4{
					{IPAddress: "172.16.1.2", Type: "private"},
				},
			},
			Region: &binarylane.Region{
				Slug: "test0",
				Name: "test 0",
			},
		},
	}

	testServerList        = bl.Servers{testServer, anotherTestServer}
	testPrivateServerList = bl.Servers{testPrivateServer}
	testKernel            = bl.Kernel{Kernel: &binarylane.Kernel{ID: 1}}
	testKernelList        = bl.Kernels{testKernel}
	testFloatingIP        = bl.FloatingIP{
		FloatingIP: &binarylane.FloatingIP{
			Server: testServer.Server,
			Region: testServer.Region,
			IP:     "127.0.0.1",
		},
	}
	testFloatingIPList = bl.FloatingIPs{testFloatingIP}

	testSnapshot = bl.Snapshot{
		Snapshot: &binarylane.Snapshot{
			ID:      "1",
			Name:    "test-snapshot",
			Regions: []string{"dev0"},
		},
	}
	testSnapshotSecondary = bl.Snapshot{
		Snapshot: &binarylane.Snapshot{
			ID:      "2",
			Name:    "test-snapshot-2",
			Regions: []string{"dev1", "dev2"},
		},
	}

	testSnapshotList = bl.Snapshots{testSnapshot, testSnapshotSecondary}
)

func assertCommandNames(t *testing.T, cmd *Command, expected ...string) {
	var names []string

	for _, c := range cmd.Commands() {
		names = append(names, c.Name())
		if c.Name() == "list" {
			assert.Contains(t, c.Aliases, "ls", "Missing 'ls' alias for 'list' command.")
			assert.NotNil(t, c.Flags().Lookup("format"), "Missing 'format' flag for 'list' command.")
		}
		if c.Name() == "get" {
			assert.NotNil(t, c.Flags().Lookup("format"), "Missing 'format' flag for 'get' command.")
		}
	}

	sort.Strings(expected)
	sort.Strings(names)
	assert.Equal(t, expected, names)
}

type testFn func(c *CmdConfig, tm *tcMocks)

type tcMocks struct {
	account           *blmocks.MockAccountService
	actions           *blmocks.MockActionsService
	balance           *blmocks.MockBalanceService
	billingHistory    *blmocks.MockBillingHistoryService
	serverActions     *blmocks.MockServerActionsService
	servers           *blmocks.MockServersService
	keys              *blmocks.MockKeysService
	sizes             *blmocks.MockSizesService
	regions           *blmocks.MockRegionsService
	images            *blmocks.MockImagesService
	imageActions      *blmocks.MockImageActionsService
	invoices          *blmocks.MockInvoicesService
	floatingIPs       *blmocks.MockFloatingIPsService
	floatingIPActions *blmocks.MockFloatingIPActionsService
	domains           *blmocks.MockDomainsService
	tags              *blmocks.MockTagsService
	snapshots         *blmocks.MockSnapshotsService
	loadBalancers     *blmocks.MockLoadBalancersService
	firewalls         *blmocks.MockFirewallsService
	projects          *blmocks.MockProjectsService
	sshRunner         *blmocks.MockRunner
	vpcs              *blmocks.MockVPCsService
}

func withTestClient(t *testing.T, tFn testFn) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tm := &tcMocks{
		account:           blmocks.NewMockAccountService(ctrl),
		actions:           blmocks.NewMockActionsService(ctrl),
		balance:           blmocks.NewMockBalanceService(ctrl),
		billingHistory:    blmocks.NewMockBillingHistoryService(ctrl),
		keys:              blmocks.NewMockKeysService(ctrl),
		sizes:             blmocks.NewMockSizesService(ctrl),
		regions:           blmocks.NewMockRegionsService(ctrl),
		images:            blmocks.NewMockImagesService(ctrl),
		imageActions:      blmocks.NewMockImageActionsService(ctrl),
		invoices:          blmocks.NewMockInvoicesService(ctrl),
		floatingIPs:       blmocks.NewMockFloatingIPsService(ctrl),
		floatingIPActions: blmocks.NewMockFloatingIPActionsService(ctrl),
		servers:           blmocks.NewMockServersService(ctrl),
		serverActions:     blmocks.NewMockServerActionsService(ctrl),
		domains:           blmocks.NewMockDomainsService(ctrl),
		tags:              blmocks.NewMockTagsService(ctrl),
		snapshots:         blmocks.NewMockSnapshotsService(ctrl),
		loadBalancers:     blmocks.NewMockLoadBalancersService(ctrl),
		firewalls:         blmocks.NewMockFirewallsService(ctrl),
		projects:          blmocks.NewMockProjectsService(ctrl),
		sshRunner:         blmocks.NewMockRunner(ctrl),
		vpcs:              blmocks.NewMockVPCsService(ctrl),
	}

	config := &CmdConfig{
		NS:   "test",
		Doit: blcli.NewTestConfig(),
		Out:  ioutil.Discard,

		// can stub this out, since the return is dictated by the mocks.
		initServices: func(c *CmdConfig) error { return nil },

		getContextAccessToken: func() string {
			return viper.GetString(blcli.ArgAccessToken)
		},

		setContextAccessToken: func(token string) {},

		Keys:              func() bl.KeysService { return tm.keys },
		Sizes:             func() bl.SizesService { return tm.sizes },
		Regions:           func() bl.RegionsService { return tm.regions },
		Images:            func() bl.ImagesService { return tm.images },
		ImageActions:      func() bl.ImageActionsService { return tm.imageActions },
		FloatingIPs:       func() bl.FloatingIPsService { return tm.floatingIPs },
		FloatingIPActions: func() bl.FloatingIPActionsService { return tm.floatingIPActions },
		Servers:           func() bl.ServersService { return tm.servers },
		ServerActions:     func() bl.ServerActionsService { return tm.serverActions },
		Domains:           func() bl.DomainsService { return tm.domains },
		Actions:           func() bl.ActionsService { return tm.actions },
		Account:           func() bl.AccountService { return tm.account },
		Balance:           func() bl.BalanceService { return tm.balance },
		BillingHistory:    func() bl.BillingHistoryService { return tm.billingHistory },
		Invoices:          func() bl.InvoicesService { return tm.invoices },
		Tags:              func() bl.TagsService { return tm.tags },
		Snapshots:         func() bl.SnapshotsService { return tm.snapshots },
		LoadBalancers:     func() bl.LoadBalancersService { return tm.loadBalancers },
		Firewalls:         func() bl.FirewallsService { return tm.firewalls },
		Projects:          func() bl.ProjectsService { return tm.projects },
		VPCs:              func() bl.VPCsService { return tm.vpcs },
	}

	tFn(config, tm)
}
