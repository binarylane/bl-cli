package commands

import (
	"strconv"
	"testing"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/go-binarylane"

	"github.com/stretchr/testify/assert"
)

var (
	testFirewall = bl.Firewall{
		Firewall: &binarylane.Firewall{
			Name: "my firewall",
		},
	}

	testFirewallList = bl.Firewalls{
		testFirewall,
	}
)

func TestFirewallCommand(t *testing.T) {
	cmd := Firewall()
	assert.NotNil(t, cmd)
	assertCommandNames(t, cmd, "get", "create", "update", "list", "list-by-server", "delete", "add-servers", "remove-servers", "add-tags", "remove-tags", "add-rules", "remove-rules")
}

func TestFirewallGet(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tm.firewalls.EXPECT().Get(fID).Return(&testFirewall, nil)

		config.Args = append(config.Args, fID)

		err := RunFirewallGet(config)
		assert.NoError(t, err)
	})
}

func TestFirewallCreate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		firewallCreateRequest := &binarylane.FirewallRequest{
			Name: "firewall",
			InboundRules: []binarylane.InboundRule{
				{
					Protocol:  "icmp",
					PortRange: "",
					Sources:   &binarylane.Sources{},
				},
				{
					Protocol:  "tcp",
					PortRange: "8000-9000",
					Sources: &binarylane.Sources{
						Addresses: []string{"127.0.0.0", "0::/0", "::/1"},
					},
				},
			},
			Tags:      []string{"backend"},
			ServerIDs: []int{1, 2},
		}
		tm.firewalls.EXPECT().Create(firewallCreateRequest).Return(&testFirewall, nil)

		config.Doit.Set(config.NS, blcli.ArgFirewallName, "firewall")
		config.Doit.Set(config.NS, blcli.ArgTagNames, []string{"backend"})
		config.Doit.Set(config.NS, blcli.ArgServerIDs, []string{"1", "2"})
		config.Doit.Set(config.NS, blcli.ArgInboundRules, "protocol:icmp protocol:tcp,ports:8000-9000,address:127.0.0.0,address:0::/0,address:::/1")

		err := RunFirewallCreate(config)
		assert.NoError(t, err)
	})
}

func TestFirewallUpdate(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		firewallUpdateRequest := &binarylane.FirewallRequest{
			Name: "firewall",
			InboundRules: []binarylane.InboundRule{
				{
					Protocol:  "tcp",
					PortRange: "8000-9000",
					Sources: &binarylane.Sources{
						Addresses: []string{"127.0.0.0"},
					},
				},
			},
			OutboundRules: []binarylane.OutboundRule{
				{
					Protocol:  "tcp",
					PortRange: "8080",
					Destinations: &binarylane.Destinations{
						LoadBalancerUIDs: []string{"lb-uuid"},
						Tags:             []string{"new-servers"},
					},
				},
				{
					Protocol:  "tcp",
					PortRange: "80",
					Destinations: &binarylane.Destinations{
						Addresses: []string{"192.168.0.0"},
					},
				},
			},
			ServerIDs: []int{1},
		}
		tm.firewalls.EXPECT().Update(fID, firewallUpdateRequest).Return(&testFirewall, nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgFirewallName, "firewall")
		config.Doit.Set(config.NS, blcli.ArgServerIDs, []string{"1"})
		config.Doit.Set(config.NS, blcli.ArgInboundRules, "protocol:tcp,ports:8000-9000,address:127.0.0.0")
		config.Doit.Set(config.NS, blcli.ArgOutboundRules, "protocol:tcp,ports:8080,load_balancer_uid:lb-uuid,tag:new-servers protocol:tcp,ports:80,address:192.168.0.0")

		err := RunFirewallUpdate(config)
		assert.NoError(t, err)
	})
}

func TestFirewallList(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		tm.firewalls.EXPECT().List().Return(testFirewallList, nil)

		err := RunFirewallList(config)
		assert.NoError(t, err)
	})
}

func TestFirewallListByServer(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		sID := 124
		tm.firewalls.EXPECT().ListByServer(sID).Return(testFirewallList, nil)
		config.Args = append(config.Args, strconv.Itoa(sID))

		err := RunFirewallListByServer(config)
		assert.NoError(t, err)
	})
}

func TestFirewallDelete(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tm.firewalls.EXPECT().Delete(fID).Return(nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgForce, true)

		err := RunFirewallDelete(config)
		assert.NoError(t, err)
	})
}

func TestFirewallAddServers(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		serverIDs := []int{1, 2}
		tm.firewalls.EXPECT().AddServers(fID, serverIDs[0], serverIDs[1]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgServerIDs, []string{"1", "2"})

		err := RunFirewallAddServers(config)
		assert.NoError(t, err)
	})
}

func TestFirewallRemoveServers(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "cde2c0d6-41e3-479e-ba60-ad971227232c"
		serverIDs := []int{1}
		tm.firewalls.EXPECT().RemoveServers(fID, serverIDs[0]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgServerIDs, []string{"1"})

		err := RunFirewallRemoveServers(config)
		assert.NoError(t, err)
	})
}

func TestFirewallAddTags(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tags := []string{"frontend", "backend"}
		tm.firewalls.EXPECT().AddTags(fID, tags[0], tags[1]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgTagNames, []string{"frontend", "backend"})

		err := RunFirewallAddTags(config)
		assert.NoError(t, err)
	})
}

func TestFirewallRemoveTags(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		tags := []string{"backend"}
		tm.firewalls.EXPECT().RemoveTags(fID, tags[0]).Return(nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgTagNames, []string{"backend"})

		err := RunFirewallRemoveTags(config)
		assert.NoError(t, err)
	})
}

func TestFirewallAddRules(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		inboundRules := []binarylane.InboundRule{
			{
				Protocol:  "tcp",
				PortRange: "80",
				Sources: &binarylane.Sources{
					Addresses: []string{"127.0.0.0", "0.0.0.0/0", "2604:A880:0002:00D0:0000:0000:32F1:E001"},
				},
			},
			{
				Protocol:  "tcp",
				PortRange: "8080",
				Sources: &binarylane.Sources{
					Tags:      []string{"backend"},
					ServerIDs: []int{1, 2, 3},
				},
			},
		}
		outboundRules := []binarylane.OutboundRule{
			{
				Protocol:  "tcp",
				PortRange: "22",
				Destinations: &binarylane.Destinations{
					LoadBalancerUIDs: []string{"lb-uuid"},
				},
			},
		}
		firewallRulesRequest := &binarylane.FirewallRulesRequest{
			InboundRules:  inboundRules,
			OutboundRules: outboundRules,
		}

		tm.firewalls.EXPECT().AddRules(fID, firewallRulesRequest).Return(nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgInboundRules, "protocol:tcp,ports:80,address:127.0.0.0,address:0.0.0.0/0,address:2604:A880:0002:00D0:0000:0000:32F1:E001 protocol:tcp,ports:8080,tag:backend,server_id:1,server_id:2,server_id:3")
		config.Doit.Set(config.NS, blcli.ArgOutboundRules, "protocol:tcp,ports:22,load_balancer_uid:lb-uuid")

		err := RunFirewallAddRules(config)
		assert.NoError(t, err)
	})
}

func TestFirewallRemoveRules(t *testing.T) {
	withTestClient(t, func(config *CmdConfig, tm *tcMocks) {
		fID := "ab06e011-6dd1-4034-9293-201f71aba299"
		inboundRules := []binarylane.InboundRule{
			{
				Protocol:  "tcp",
				PortRange: "80",
				Sources: &binarylane.Sources{
					Addresses: []string{"0.0.0.0/0"},
				},
			},
		}
		outboundRules := []binarylane.OutboundRule{
			{
				Protocol:  "tcp",
				PortRange: "22",
				Destinations: &binarylane.Destinations{
					Tags:      []string{"back:end"},
					Addresses: []string{"::/0"},
				},
			},
		}
		firewallRulesRequest := &binarylane.FirewallRulesRequest{
			InboundRules:  inboundRules,
			OutboundRules: outboundRules,
		}

		tm.firewalls.EXPECT().RemoveRules(fID, firewallRulesRequest).Return(nil)

		config.Args = append(config.Args, fID)
		config.Doit.Set(config.NS, blcli.ArgInboundRules, "protocol:tcp,ports:80,address:0.0.0.0/0")
		config.Doit.Set(config.NS, blcli.ArgOutboundRules, "protocol:tcp,ports:22,tag:back:end,address:::/0")

		err := RunFirewallRemoveRules(config)
		assert.NoError(t, err)
	})
}
