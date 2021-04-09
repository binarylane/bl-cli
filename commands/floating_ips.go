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
	"errors"
	"fmt"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/commands/displayers"
	"github.com/binarylane/go-binarylane"
	"github.com/spf13/cobra"
)

// FloatingIP creates the command hierarchy for floating ips.
func FloatingIP() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "floating-ip",
			Short: "Display commands to manage floating IP addresses",
			Long: `The sub-commands of ` + "`" + `bl compute floating-ip` + "`" + ` manage floating IP addresses.
Floating IPs are publicly-accessible static IP addresses that can be mapped to one of your Servers. They can be used to create highly available setups or other configurations requiring movable addresses.
Floating IPs are bound to a specific region.`,
			Aliases: []string{"fip"},
		},
	}

	cmdFloatingIPCreate := CmdBuilder(cmd, RunFloatingIPCreate, "create", "Create a new floating IP address", `Use this command to create a new floating IP address.

A floating IP address must be either assigned to a Server or reserved to a region.`, Writer,
		aliasOpt("c"), displayerType(&displayers.FloatingIP{}))
	AddStringFlag(cmdFloatingIPCreate, blcli.ArgRegionSlug, "", "",
		fmt.Sprintf("Region where to create the floating IP address. (mutually exclusive with `--%s`)",
			blcli.ArgServerID))
	AddIntFlag(cmdFloatingIPCreate, blcli.ArgServerID, "", 0,
		fmt.Sprintf("The ID of the Server to assign the floating IP to (mutually exclusive with `--%s`).",
			blcli.ArgRegionSlug))

	CmdBuilder(cmd, RunFloatingIPGet, "get <floating-ip>", "Retrieve information about a floating IP address", "Use this command to retrieve detailed information about a floating IP address.", Writer,
		aliasOpt("g"), displayerType(&displayers.FloatingIP{}))

	cmdRunFloatingIPDelete := CmdBuilder(cmd, RunFloatingIPDelete, "delete <floating-ip>", "Permanently delete a floating IP address", "Use this command to permanently delete a floating IP address. This is irreversible.", Writer, aliasOpt("d", "rm"))
	AddBoolFlag(cmdRunFloatingIPDelete, blcli.ArgForce, blcli.ArgShortForce, false, "Force floating IP delete")

	cmdFloatingIPList := CmdBuilder(cmd, RunFloatingIPList, "list", "List all floating IP addresses on your account", "Use this command to list all the floating IP addresses on your account.", Writer,
		aliasOpt("ls"), displayerType(&displayers.FloatingIP{}))
	AddStringFlag(cmdFloatingIPList, blcli.ArgRegionSlug, "", "", "The region the floating IP address resides in")

	return cmd
}

// RunFloatingIPCreate runs floating IP create.
func RunFloatingIPCreate(c *CmdConfig) error {
	fis := c.FloatingIPs()

	// ignore errors since we don't know which one is valid
	region, _ := c.Doit.GetString(c.NS, blcli.ArgRegionSlug)
	serverID, _ := c.Doit.GetInt(c.NS, blcli.ArgServerID)

	if region == "" && serverID == 0 {
		return blcli.NewMissingArgsErr("Region and Server ID can't both be blank.")
	}

	if region != "" && serverID != 0 {
		return fmt.Errorf("Specify region or Server ID when creating a floating IP address.")
	}

	req := &binarylane.FloatingIPCreateRequest{
		Region:   region,
		ServerID: serverID,
	}

	ip, err := fis.Create(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	item := &displayers.FloatingIP{FloatingIPs: bl.FloatingIPs{*ip}}
	return c.Display(item)
}

// RunFloatingIPGet retrieves a floating IP's details.
func RunFloatingIPGet(c *CmdConfig) error {
	fis := c.FloatingIPs()

	err := ensureOneArg(c)
	if err != nil {
		return err
	}

	ip := c.Args[0]

	if len(ip) < 1 {
		return errors.New("Invalid IP address")
	}

	fip, err := fis.Get(ip)
	if err != nil {
		return err
	}

	item := &displayers.FloatingIP{FloatingIPs: bl.FloatingIPs{*fip}}
	return c.Display(item)
}

// RunFloatingIPDelete runs floating IP delete.
func RunFloatingIPDelete(c *CmdConfig) error {
	fis := c.FloatingIPs()

	err := ensureOneArg(c)
	if err != nil {
		return err
	}

	force, err := c.Doit.GetBool(c.NS, blcli.ArgForce)
	if err != nil {
		return err
	}

	if force || AskForConfirmDelete("floating IP", 1) == nil {
		ip := c.Args[0]
		return fis.Delete(ip)
	}

	return fmt.Errorf("Operation aborted.")
}

// RunFloatingIPList runs floating IP create.
func RunFloatingIPList(c *CmdConfig) error {
	fis := c.FloatingIPs()

	region, err := c.Doit.GetString(c.NS, blcli.ArgRegionSlug)
	if err != nil {
		return err
	}

	list, err := fis.List()
	if err != nil {
		return err
	}

	fips := &displayers.FloatingIP{FloatingIPs: bl.FloatingIPs{}}
	for _, fip := range list {
		var skip bool
		if region != "" && region != fip.Region.Slug {
			skip = true
		}

		if !skip {
			fips.FloatingIPs = append(fips.FloatingIPs, fip)
		}
	}

	item := fips
	return c.Display(item)
}
