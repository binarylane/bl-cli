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
	"fmt"
	"strconv"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/commands/displayers"
	"github.com/binarylane/go-binarylane"
	"github.com/spf13/cobra"
)

// VPCs creates the vpcs command.
func VPCs() *Command {

	cmd := &Command{
		Command: &cobra.Command{
			Use:   "vpcs",
			Short: "Display commands that manage VPCs",
			Long: `The commands under ` + "`" + `bl vpcs` + "`" + ` are for managing your VPCs.

With the vpcs command, you can list, create, or delete VPCs, and manage their configuration details.`,
		},
	}

	vpcDetail := `

- The VPC's ID
- The uniform resource name (URN) for the VPC
- The VPC's name
- The VPC's description
- The range of IP addresses in the VPC in CIDR notation
- The datacenter region slug the VPC is located in
- The VPC's default boolean value indicating whether or not the VPC is the default one for the region
- The VPC's creation date, in ISO8601 combined date and time format
`

	CmdBuilder(cmd, RunVPCGet, "get <id>", "Retrieve a VPC", "Use this command to retrieve information about a VPC, including:"+vpcDetail, Writer,
		aliasOpt("g"), displayerType(&displayers.VPC{}))

	cmdRecordCreate := CmdBuilder(cmd, RunVPCCreate, "create",
		"Create a new VPC", "Use this command to create a new VPC on your account.", Writer, aliasOpt("c"))
	AddStringFlag(cmdRecordCreate, blcli.ArgVPCName, "", "",
		"The VPC's name", requiredOpt())
	AddStringFlag(cmdRecordCreate, blcli.ArgVPCDescription, "", "", "The VPC's name")
	AddStringFlag(cmdRecordCreate, blcli.ArgVPCIPRange, "", "",
		"The range of IP addresses in the VPC in CIDR notation, e.g.: `10.240.0.0/16`")
	AddStringFlag(cmdRecordCreate, blcli.ArgRegionSlug, "", "", "The VPC's region slug, e.g.: `syd`", requiredOpt())

	cmdRecordUpdate := CmdBuilder(cmd, RunVPCUpdate, "update <id>",
		"Update a VPC's configuration", `Use this command to update the configuration of a specified VPC.`, Writer, aliasOpt("u"))
	AddStringFlag(cmdRecordUpdate, blcli.ArgVPCName, "", "",
		"The VPC's name", requiredOpt())
	AddStringFlag(cmdRecordUpdate, blcli.ArgVPCDescription, "", "",
		"The VPC's description")
	AddBoolFlag(cmdRecordUpdate, blcli.ArgVPCDefault, "", false,
		"The VPC's default state")

	CmdBuilder(cmd, RunVPCList, "list", "List VPCs", "Use this command to get a list of the VPCs on your account, including the following information for each:"+vpcDetail, Writer,
		aliasOpt("ls"), displayerType(&displayers.VPC{}))

	cmdRunRecordDelete := CmdBuilder(cmd, RunVPCDelete, "delete <id>",
		"Permanently delete a VPC", `Use this command to permanently delete the specified VPC. This is irreversible.`, Writer, aliasOpt("d", "rm"))
	AddBoolFlag(cmdRunRecordDelete, blcli.ArgForce, blcli.ArgShortForce, false,
		"Delete the VPC without a confirmation prompt")

	return cmd
}

// RunVPCGet retrieves an existing VPC by its identifier.
func RunVPCGet(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}

	vpcID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	vpcs := c.VPCs()
	vpc, err := vpcs.Get(vpcID)
	if err != nil {
		return err
	}

	item := &displayers.VPC{VPCs: bl.VPCs{*vpc}}
	return c.Display(item)
}

// RunVPCList lists VPCs.
func RunVPCList(c *CmdConfig) error {
	vpcs := c.VPCs()
	list, err := vpcs.List()
	if err != nil {
		return err
	}

	item := &displayers.VPC{VPCs: list}
	return c.Display(item)
}

// RunVPCCreate creates a new VPC with a given configuration.
func RunVPCCreate(c *CmdConfig) error {
	r := new(binarylane.VPCCreateRequest)
	name, err := c.Doit.GetString(c.NS, blcli.ArgVPCName)
	if err != nil {
		return err
	}
	r.Name = name

	desc, err := c.Doit.GetString(c.NS, blcli.ArgVPCDescription)
	if err != nil {
		return err
	}
	r.Description = desc

	ipRange, err := c.Doit.GetString(c.NS, blcli.ArgVPCIPRange)
	if err != nil {
		return err
	}
	r.IPRange = ipRange

	rSlug, err := c.Doit.GetString(c.NS, blcli.ArgRegionSlug)
	if err != nil {
		return err
	}
	r.RegionSlug = rSlug

	vpcs := c.VPCs()
	vpc, err := vpcs.Create(r)
	if err != nil {
		return err
	}

	item := &displayers.VPC{VPCs: bl.VPCs{*vpc}}
	return c.Display(item)
}

// RunVPCUpdate updates an existing VPC with new configuration.
func RunVPCUpdate(c *CmdConfig) error {
	if len(c.Args) == 0 {
		return blcli.NewMissingArgsErr(c.NS)
	}
	vpcID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	r := new(binarylane.VPCUpdateRequest)
	name, err := c.Doit.GetString(c.NS, blcli.ArgVPCName)
	if err != nil {
		return err
	}
	r.Name = name

	desc, err := c.Doit.GetString(c.NS, blcli.ArgVPCDescription)
	if err != nil {
		return err
	}
	r.Description = desc

	def, err := c.Doit.GetBoolPtr(c.NS, blcli.ArgVPCDefault)
	if err != nil {
		return err
	}

	if def != nil {
		r.Default = boolPtr(true)
	}

	vpcs := c.VPCs()
	vpc, err := vpcs.Update(vpcID, r)
	if err != nil {
		return err
	}

	item := &displayers.VPC{VPCs: bl.VPCs{*vpc}}
	return c.Display(item)
}

// RunVPCDelete deletes a VPC by its identifier.
func RunVPCDelete(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	vpcID, err := strconv.Atoi(c.Args[0])
	if err != nil {
		return err
	}

	force, err := c.Doit.GetBool(c.NS, blcli.ArgForce)
	if err != nil {
		return err
	}

	if force || AskForConfirmDelete("VPC", 1) == nil {
		vpcs := c.VPCs()
		if err := vpcs.Delete(vpcID); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("operation aborted")
	}

	return nil
}
