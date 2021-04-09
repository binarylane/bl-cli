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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/commands/displayers"
	"github.com/binarylane/go-binarylane"

	"github.com/spf13/cobra"
)

// Firewall creates the firewall command.
func Firewall() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "firewall",
			Short: "Display commands to manage cloud firewalls",
			Long: `The sub-commands of ` + "`" + `bl compute firewall` + "`" + ` manage BinaryLane cloud firewalls.

Cloud firewalls provide the ability to restrict network access to and from a Server, allowing you to define which ports accept inbound or outbound connections. With these commands, you can list, create, or delete Cloud firewalls, as well as modify access rules.

A firewall's ` + "`" + `inbound_rules` + "`" + ` and ` + "`" + `outbound_rules` + "`" + ` attributes contain arrays of objects as their values. These objects contain the standard attributes of their associated types, which can be found below.

Inbound access rules specify the protocol (TCP, UDP, or ICMP), ports, and sources for inbound traffic that will be allowed through the Firewall to the target Servers. The ` + "`" + `ports` + "`" + ` attribute may contain a single port, a range of ports (e.g. ` + "`" + `8000-9000` + "`" + `), or ` + "`" + `all` + "`" + ` to allow traffic on all ports for the specified protocol. The ` + "`" + `sources` + "`" + ` attribute will contain an object specifying a whitelist of sources from which traffic will be accepted.`,
		},
	}
	fwDetail := `

	- The firewall's ID
	- The firewall's name
	- The status of the firewall. This can be ` + "`" + `waiting` + "`" + `, ` + "`" + `succeeded` + "`" + `, or ` + "`" + `failed` + "`" + `.
	- The firewall's creation date, in ISO8601 combined date and time format.
	- Any pending changes to the firewall. These can be ` + "`" + `server_id` + "`" + `, ` + "`" + `removing` + "`" + `, and ` + "`" + `status` + "`" + `.
	  When empty, all changes have been successfully applied.
	- The inbound rules for the firewall.
	- The outbound rules for the firewall.
	- The IDs of Servers assigned to the firewall.
	- The tags assigned to the firewall.
`
	inboundRulesTxt := "A comma-separated key-value list that defines an inbound rule, e.g.: `protocol:tcp,ports:22,server_id:123`. Use a quoted string of space-separated values for multiple rules."
	outboundRulesTxt := "A comma-separate key-value list the defines an outbound rule, e.g.: `protocol:tcp,ports:22,address:0.0.0.0/0`. Use a quoted string of space-separated values for multiple rules."
	serverIDRulesTxt := "A comma-separated list of Server IDs to place behind the cloud firewall, e.g.: `123,456`"
	tagNameRulesTxt := "A comma-separated list of tag names to apply to the cloud firewall, e.g.: `frontend,backend`"

	CmdBuilder(cmd, RunFirewallGet, "get <id>", "Retrieve information about a cloud firewall", `Use this command to get information about an existing cloud firewall, including:`+fwDetail, Writer, aliasOpt("g"), displayerType(&displayers.Firewall{}))

	cmdFirewallCreate := CmdBuilder(cmd, RunFirewallCreate, "create", "Create a new cloud firewall", `Use this command to create a cloud firewall. This command must contain at least one inbound or outbound access rule.`, Writer, aliasOpt("c"), displayerType(&displayers.Firewall{}))
	AddStringFlag(cmdFirewallCreate, blcli.ArgFirewallName, "", "", "Firewall name", requiredOpt())
	AddStringFlag(cmdFirewallCreate, blcli.ArgInboundRules, "", "", inboundRulesTxt)
	AddStringFlag(cmdFirewallCreate, blcli.ArgOutboundRules, "", "", outboundRulesTxt)
	AddStringSliceFlag(cmdFirewallCreate, blcli.ArgServerIDs, "", []string{}, serverIDRulesTxt)
	AddStringSliceFlag(cmdFirewallCreate, blcli.ArgTagNames, "", []string{}, tagNameRulesTxt)

	cmdFirewallUpdate := CmdBuilder(cmd, RunFirewallUpdate, "update <id>", "Update a cloud firewall's configuration", `Use this command to update the configuration of an existing cloud firewall. The request should contain a full representation of the Firewall, including existing attributes. Note: Any attributes that are not provided will be reset to their default values.`, Writer, aliasOpt("u"), displayerType(&displayers.Firewall{}))
	AddStringFlag(cmdFirewallUpdate, blcli.ArgFirewallName, "", "", "Firewall name", requiredOpt())
	AddStringFlag(cmdFirewallUpdate, blcli.ArgInboundRules, "", "", inboundRulesTxt)
	AddStringFlag(cmdFirewallUpdate, blcli.ArgOutboundRules, "", "", outboundRulesTxt)
	AddStringSliceFlag(cmdFirewallUpdate, blcli.ArgServerIDs, "", []string{}, serverIDRulesTxt)
	AddStringSliceFlag(cmdFirewallUpdate, blcli.ArgTagNames, "", []string{}, tagNameRulesTxt)

	CmdBuilder(cmd, RunFirewallList, "list", "List the cloud firewalls on your account", `Use this command to retrieve a list of cloud firewalls.`, Writer, aliasOpt("ls"), displayerType(&displayers.Firewall{}))

	CmdBuilder(cmd, RunFirewallListByServer, "list-by-server <server_id>", "List firewalls by Server", `Use this command to list cloud firewalls by the ID of a Server assigned to the firewall.`, Writer, displayerType(&displayers.Firewall{}))

	cmdRunRecordDelete := CmdBuilder(cmd, RunFirewallDelete, "delete <id>...", "Permanently delete a cloud firewall", `Use this command to permanently delete a cloud firewall. This is irreversable, but does not delete any Servers assigned to the cloud firewall.`, Writer, aliasOpt("d", "rm"))
	AddBoolFlag(cmdRunRecordDelete, blcli.ArgForce, blcli.ArgShortForce, false, "Delete firewall without confirmation prompt")

	cmdAddServers := CmdBuilder(cmd, RunFirewallAddServers, "add-servers <id>", "Add Servers to a cloud firewall", `Use this command to add Servers to the cloud firewall.`, Writer)
	AddStringSliceFlag(cmdAddServers, blcli.ArgServerIDs, "", []string{}, serverIDRulesTxt)

	cmdRemoveServers := CmdBuilder(cmd, RunFirewallRemoveServers, "remove-servers <id>", "Remove Servers from a cloud firewall", `Use this command to remove Servers from a cloud firewall.`, Writer)
	AddStringSliceFlag(cmdRemoveServers, blcli.ArgServerIDs, "", []string{}, serverIDRulesTxt)

	cmdAddTags := CmdBuilder(cmd, RunFirewallAddTags, "add-tags <id>", "Add tags to a cloud firewall", `Use this command to add tags to a cloud firewall. This adds all assets using that tag to the firewall`, Writer)
	AddStringSliceFlag(cmdAddTags, blcli.ArgTagNames, "", []string{}, tagNameRulesTxt)

	cmdRemoveTags := CmdBuilder(cmd, RunFirewallRemoveTags, "remove-tags <id>", "Remove tags from a cloud firewall", `Use this command to remove tags from a cloud firewall. This removes all assets using that tag from the firewall.`, Writer)
	AddStringSliceFlag(cmdRemoveTags, blcli.ArgTagNames, "", []string{}, tagNameRulesTxt)

	cmdAddRules := CmdBuilder(cmd, RunFirewallAddRules, "add-rules <id>", "Add inbound or outbound rules to a cloud firewall", `Use this command to add inbound or outbound rules to a cloud firewall.`, Writer)
	AddStringFlag(cmdAddRules, blcli.ArgInboundRules, "", "", inboundRulesTxt)
	AddStringFlag(cmdAddRules, blcli.ArgOutboundRules, "", "", outboundRulesTxt)

	cmdRemoveRules := CmdBuilder(cmd, RunFirewallRemoveRules, "remove-rules <id>", "Remove inbound or outbound rules from a cloud firewall", `Use this command to remove inbound or outbound rules from a cloud firewall.`, Writer)
	AddStringFlag(cmdRemoveRules, blcli.ArgInboundRules, "", "", inboundRulesTxt)
	AddStringFlag(cmdRemoveRules, blcli.ArgOutboundRules, "", "", outboundRulesTxt)

	return cmd
}

// RunFirewallGet retrieves an existing Firewall by its identifier.
func RunFirewallGet(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	id := c.Args[0]

	fs := c.Firewalls()
	f, err := fs.Get(id)
	if err != nil {
		return err
	}

	item := &displayers.Firewall{Firewalls: bl.Firewalls{*f}}
	return c.Display(item)
}

// RunFirewallCreate creates a new Firewall with a given configuration.
func RunFirewallCreate(c *CmdConfig) error {
	r := new(binarylane.FirewallRequest)
	if err := buildFirewallRequestFromArgs(c, r); err != nil {
		return err
	}

	fs := c.Firewalls()
	f, err := fs.Create(r)
	if err != nil {
		return err
	}

	item := &displayers.Firewall{Firewalls: bl.Firewalls{*f}}
	return c.Display(item)
}

// RunFirewallUpdate updates an existing Firewall with new configuration.
func RunFirewallUpdate(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return blcli.NewMissingArgsErr(c.NS)
	}
	fID := c.Args[0]

	r := new(binarylane.FirewallRequest)
	if err := buildFirewallRequestFromArgs(c, r); err != nil {
		return err
	}

	fs := c.Firewalls()
	f, err := fs.Update(fID, r)
	if err != nil {
		return err
	}

	item := &displayers.Firewall{Firewalls: bl.Firewalls{*f}}
	return c.Display(item)
}

// RunFirewallList lists Firewalls.
func RunFirewallList(c *CmdConfig) error {
	fs := c.Firewalls()
	list, err := fs.List()
	if err != nil {
		return err
	}

	items := &displayers.Firewall{Firewalls: list}
	return c.Display(items)
}

// RunFirewallListByServer lists Firewalls for a given Server.
func RunFirewallListByServer(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	sID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return fmt.Errorf("invalid server id: [%v]", c.Args[0])
	}

	fs := c.Firewalls()
	list, err := fs.ListByServer(sID)
	if err != nil {
		return err
	}

	items := &displayers.Firewall{Firewalls: list}
	return c.Display(items)
}

// RunFirewallDelete deletes a Firewall by its identifier.
func RunFirewallDelete(c *CmdConfig) error {
	if len(c.Args) < 1 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	force, err := c.Doit.GetBool(c.NS, blcli.ArgForce)
	if err != nil {
		return err
	}

	fs := c.Firewalls()
	if force || AskForConfirmDelete("firewall", len(c.Args)) == nil {
		for _, id := range c.Args {
			if err := fs.Delete(id); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("Operation aborted.")
	}

	return nil
}

// RunFirewallAddServers adds servers to a Firewall.
func RunFirewallAddServers(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	fID := c.Args[0]

	serverIDsList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgServerIDs)
	if err != nil {
		return err
	}

	serverIDs, err := extractServerIDs(serverIDsList)
	if err != nil {
		return err
	}

	return c.Firewalls().AddServers(fID, serverIDs...)
}

// RunFirewallRemoveServers removes servers from a Firewall.
func RunFirewallRemoveServers(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	fID := c.Args[0]

	serverIDsList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgServerIDs)
	if err != nil {
		return err
	}

	serverIDs, err := extractServerIDs(serverIDsList)
	if err != nil {
		return err
	}

	return c.Firewalls().RemoveServers(fID, serverIDs...)
}

// RunFirewallAddTags adds tags to a Firewall.
func RunFirewallAddTags(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	fID := c.Args[0]

	tagList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgTagNames)
	if err != nil {
		return err
	}

	return c.Firewalls().AddTags(fID, tagList...)
}

// RunFirewallRemoveTags removes tags from a Firewall.
func RunFirewallRemoveTags(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	fID := c.Args[0]

	tagList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgTagNames)
	if err != nil {
		return err
	}

	return c.Firewalls().RemoveTags(fID, tagList...)
}

// RunFirewallAddRules adds rules to a Firewall.
func RunFirewallAddRules(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	fID := c.Args[0]

	rr := new(binarylane.FirewallRulesRequest)
	if err := buildFirewallRulesRequestFromArgs(c, rr); err != nil {
		return err
	}

	return c.Firewalls().AddRules(fID, rr)
}

// RunFirewallRemoveRules removes rules from a Firewall.
func RunFirewallRemoveRules(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	fID := c.Args[0]

	rr := new(binarylane.FirewallRulesRequest)
	if err := buildFirewallRulesRequestFromArgs(c, rr); err != nil {
		return err
	}

	return c.Firewalls().RemoveRules(fID, rr)
}

func buildFirewallRequestFromArgs(c *CmdConfig, r *binarylane.FirewallRequest) error {
	name, err := c.Doit.GetString(c.NS, blcli.ArgFirewallName)
	if err != nil {
		return err
	}
	r.Name = name

	ira, err := c.Doit.GetString(c.NS, blcli.ArgInboundRules)
	if err != nil {
		return err
	}

	inboundRules, err := extractInboundRules(ira)
	if err != nil {
		return err
	}
	r.InboundRules = inboundRules

	ora, err := c.Doit.GetString(c.NS, blcli.ArgOutboundRules)
	if err != nil {
		return err
	}

	outboundRules, err := extractOutboundRules(ora)
	if err != nil {
		return err
	}
	r.OutboundRules = outboundRules

	serverIDsList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgServerIDs)
	if err != nil {
		return err
	}

	serverIDs, err := extractServerIDs(serverIDsList)
	if err != nil {
		return err
	}
	r.ServerIDs = serverIDs

	tagsList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgTagNames)
	if err != nil {
		return err
	}
	r.Tags = tagsList

	return nil
}

func buildFirewallRulesRequestFromArgs(c *CmdConfig, rr *binarylane.FirewallRulesRequest) error {
	ira, err := c.Doit.GetString(c.NS, blcli.ArgInboundRules)
	if err != nil {
		return err
	}

	inboundRules, err := extractInboundRules(ira)
	if err != nil {
		return err
	}
	rr.InboundRules = inboundRules

	ora, err := c.Doit.GetString(c.NS, blcli.ArgOutboundRules)
	if err != nil {
		return err
	}

	outboundRules, err := extractOutboundRules(ora)
	if err != nil {
		return err
	}
	rr.OutboundRules = outboundRules

	return nil
}

func extractInboundRules(s string) (rules []binarylane.InboundRule, err error) {
	if len(s) == 0 {
		return nil, nil
	}

	list := strings.Split(s, " ")
	for _, v := range list {
		rule, err := extractRule(v, "sources")
		if err != nil {
			return nil, err
		}
		mr, _ := json.Marshal(rule)
		ir := &binarylane.InboundRule{}
		json.Unmarshal(mr, ir)
		rules = append(rules, *ir)
	}

	return rules, nil
}

func extractOutboundRules(s string) (rules []binarylane.OutboundRule, err error) {
	if len(s) == 0 {
		return nil, nil
	}

	list := strings.Split(s, " ")
	for _, v := range list {
		rule, err := extractRule(v, "destinations")
		if err != nil {
			return nil, err
		}
		mr, _ := json.Marshal(rule)
		or := &binarylane.OutboundRule{}
		json.Unmarshal(mr, or)
		rules = append(rules, *or)
	}

	return rules, nil
}

func extractRule(ruleStr string, sd string) (map[string]interface{}, error) {
	rule := map[string]interface{}{}
	var serverIDs []int
	var addresses, lbUIDs, tags []string

	kvs := strings.Split(ruleStr, ",")
	for _, v := range kvs {
		pair := strings.SplitN(v, ":", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("Unexpected input value [%v], must be a key:value pair", pair)
		}

		switch pair[0] {
		case "address":
			addresses = append(addresses, pair[1])
		case "server_id":
			i, err := strconv.Atoi(pair[1])
			if err != nil {
				return nil, fmt.Errorf("Provided value [%v] for server id is not of type int", pair[0])
			}
			serverIDs = append(serverIDs, i)
		case "load_balancer_uid":
			lbUIDs = append(lbUIDs, pair[1])
		case "tag":
			tags = append(tags, pair[1])
		default:
			rule[pair[0]] = pair[1]
		}
	}

	rule[sd] = map[string]interface{}{
		"addresses":          addresses,
		"server_ids":         serverIDs,
		"load_balancer_uids": lbUIDs,
		"tags":               tags,
	}

	return rule, nil
}
