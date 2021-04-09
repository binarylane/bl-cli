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
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/binarylane/bl-cli"
	"github.com/binarylane/bl-cli/bl"
	"github.com/binarylane/bl-cli/commands/displayers"
	"github.com/binarylane/go-binarylane"
	"github.com/gobwas/glob"
	"github.com/spf13/cobra"
)

// Server creates the server command.
func Server() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "server",
			Aliases: []string{"d"},
			Short:   "Display commands to manage servers",
			Long:    "Use the subcommands of `bl compute server` to list, create, or delete servers.",
		},
	}
	serverDetails := `

	- The Server's ID
	- The Server's name
	- The Server's Public IPv4 Address
	- The Server's Private IPv4 Address
	- The Server's IPv6 Address
	- The memory size of the Server in MB
	- The number of vCPUs on the Server
	- The size of the Server's disk in GB
	- The Server's region
	- The image the Server was created from
	- The status of the Server; can be ` + "`" + `new` + "`" + `, ` + "`" + `active` + "`" + `, ` + "`" + `off` + "`" + `, or ` + "`" + `archive` + "`" + `
	- The tags assigned to the Server
	- A list of features enabled for the Server. Examples are ` + "`" + `backups` + "`" + `, ` + "`" + `ipv6` + "`" + `, ` + "`" + `monitoring` + "`" + `, ` + "`" + `private_networking` + "`" + `
	- The IDs of block storage volumes attached to the Server
	`
	CmdBuilder(cmd, RunServerActions, "actions <server-id>", "List Server actions", `Use this command to list the available actions that can be taken on a Server. These can be things like rebooting, resizing, and snapshotting the Server.`, Writer,
		aliasOpt("a"), displayerType(&displayers.Action{}))

	CmdBuilder(cmd, RunServerBackups, "backups <server-id>", "List Server backups", `Use this command to list Server backups.`, Writer,
		aliasOpt("b"), displayerType(&displayers.Image{}))

	serverCreateLongDesc := `Use this command to create a new Server. Required values are name, region, size, and image. For example, to create an Ubuntu 20.04 with 1 vCPU and 1 GB of RAM in the Sydney region, run:

	bl compute server create --image ubuntu-20-04-lts --size std-min --region syd example.com
`

	cmdServerCreate := CmdBuilder(cmd, RunServerCreate, "create <server-name>...", "Create a new Server", serverCreateLongDesc, Writer,
		aliasOpt("c"), displayerType(&displayers.Server{}))
	AddStringSliceFlag(cmdServerCreate, blcli.ArgSSHKeys, "", []string{}, "A list of SSH key fingerprints or IDs of the SSH keys to embed in the Server's root account upon creation")
	AddStringFlag(cmdServerCreate, blcli.ArgUserData, "", "", "User-data to configure the Server on first boot")
	AddStringFlag(cmdServerCreate, blcli.ArgUserDataFile, "", "", "The path to a file containing user-data to configure the Server on first boot")
	AddBoolFlag(cmdServerCreate, blcli.ArgCommandWait, "", false, "Wait for Server creation to complete before returning")
	AddStringFlag(cmdServerCreate, blcli.ArgRegionSlug, "", "", "A slug indicating the region where the Server will be created (e.g. `syd`). Run `bl compute region list` for a list of valid regions.",
		requiredOpt())
	AddStringFlag(cmdServerCreate, blcli.ArgSizeSlug, "", "", "A slug indicating the size of the Server (e.g. `std-min`). Run `bl compute size list` for a list of valid sizes.",
		requiredOpt())
	AddBoolFlag(cmdServerCreate, blcli.ArgBackups, "", false, "Enables backups for the Server")
	AddBoolFlag(cmdServerCreate, blcli.ArgIPv6, "", false, "Enables IPv6 support and assigns an IPv6 address")
	AddBoolFlag(cmdServerCreate, blcli.ArgPrivateNetworking, "", false, "Enables private networking for the Server by provisioning it inside of your account's default VPC for the region")
	AddBoolFlag(cmdServerCreate, blcli.ArgMonitoring, "", false, "Install the BinaryLane agent for additional monitoring")
	AddStringFlag(cmdServerCreate, blcli.ArgImage, "", "", "An ID or slug indicating the image the Server will be based-on (e.g. `ubuntu-20-04-lts`). Use the commands under `bl compute image` to find additional images.",
		requiredOpt())
	AddStringFlag(cmdServerCreate, blcli.ArgTagName, "", "", "A tag name to be applied to the Server")
	AddStringFlag(cmdServerCreate, blcli.ArgVPCID, "", "", "The ID of a non-default VPC to create the Server in")
	AddStringSliceFlag(cmdServerCreate, blcli.ArgTagNames, "", []string{}, "A list of tag names to be applied to the Server")

	AddStringSliceFlag(cmdServerCreate, blcli.ArgVolumeList, "", []string{}, "A list of block storage volume IDs to attach to the Server")

	cmdRunServerDelete := CmdBuilder(cmd, RunServerDelete, "delete <server-id|server-name>...", "Permanently delete a Server", `Use this command to permanently delete a Server. This is irreversible.`, Writer,
		aliasOpt("d", "del", "rm"))
	AddBoolFlag(cmdRunServerDelete, blcli.ArgForce, blcli.ArgShortForce, false, "Delete the Server without a confirmation prompt")
	AddStringFlag(cmdRunServerDelete, blcli.ArgTagName, "", "", "Tag name")

	cmdRunServerGet := CmdBuilder(cmd, RunServerGet, "get <server-id|server-name>", "Retrieve information about a Server", `Use this command to retrieve information about a Server, including:`+serverDetails, Writer,
		aliasOpt("g"), displayerType(&displayers.Server{}))
	AddStringFlag(cmdRunServerGet, blcli.ArgTemplate, "", "", "Go template format. Sample values: `{{.ID}}`, `{{.Name}}`, `{{.Memory}}`, `{{.Region.Name}}`, `{{.Image}}`, `{{.Tags}}`")

	CmdBuilder(cmd, RunServerKernels, "kernels <server-id>", "List available Server kernels", `Use this command to retrieve a list of all kernels available to a Server.`, Writer,
		aliasOpt("k"), displayerType(&displayers.Kernel{}))

	cmdRunServerList := CmdBuilder(cmd, RunServerList, "list [GLOB]", "List Servers on your account", `Use this command to retrieve a list of Servers, including the following information about each:`+serverDetails, Writer,
		aliasOpt("ls"), displayerType(&displayers.Server{}))
	AddStringFlag(cmdRunServerList, blcli.ArgRegionSlug, "", "", "Server region")
	AddStringFlag(cmdRunServerList, blcli.ArgTagName, "", "", "Tag name")

	CmdBuilder(cmd, RunServerNeighbors, "neighbors <server-id>", "List a Server's neighbors on your account", `Use this command to get a list of your Servers that are on the same physical hardware, including the following details:`+serverDetails, Writer,
		aliasOpt("n"), displayerType(&displayers.Server{}))

	CmdBuilder(cmd, RunServerSnapshots, "snapshots <server-id>", "List all snapshots for a Server", `Use this command to get a list of snapshots created from this Server.`, Writer,
		aliasOpt("s"), displayerType(&displayers.Image{}))

	cmdRunServerTag := CmdBuilder(cmd, RunServerTag, "tag <server-id|server-name>", "Add a tag to a Server", "Use this command to tag a Server. Specify the tag with the `--tag-name` flag.", Writer)
	AddStringFlag(cmdRunServerTag, blcli.ArgTagName, "", "", "Tag name to use; can be a new or existing tag",
		requiredOpt())

	cmdRunServerUntag := CmdBuilder(cmd, RunServerUntag, "untag <server-id|server-name>", "Remove a tag from a Server", "Use this command to remove a tag from a Server, specified with the `--tag-name` flag.", Writer)
	AddStringSliceFlag(cmdRunServerUntag, blcli.ArgTagName, "", []string{}, "Tag name to remove from Server")

	return cmd
}

// RunServerActions returns a list of actions for a server.
func RunServerActions(c *CmdConfig) error {

	ds := c.Servers()

	id, err := getServerIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Actions(id)
	if err != nil {
		return err
	}
	item := &displayers.Action{Actions: list}
	return c.Display(item)
}

// RunServerBackups returns a list of backup images for a server.
func RunServerBackups(c *CmdConfig) error {

	ds := c.Servers()

	id, err := getServerIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Backups(id)
	if err != nil {
		return err
	}

	item := &displayers.Image{Images: list}
	return c.Display(item)
}

// RunServerCreate creates a server.
func RunServerCreate(c *CmdConfig) error {

	if len(c.Args) < 1 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	region, err := c.Doit.GetString(c.NS, blcli.ArgRegionSlug)
	if err != nil {
		return err
	}

	size, err := c.Doit.GetString(c.NS, blcli.ArgSizeSlug)
	if err != nil {
		return err
	}

	backups, err := c.Doit.GetBool(c.NS, blcli.ArgBackups)
	if err != nil {
		return err
	}

	ipv6, err := c.Doit.GetBool(c.NS, blcli.ArgIPv6)
	if err != nil {
		return err
	}

	privateNetworking, err := c.Doit.GetBool(c.NS, blcli.ArgPrivateNetworking)
	if err != nil {
		return err
	}

	monitoring, err := c.Doit.GetBool(c.NS, blcli.ArgMonitoring)
	if err != nil {
		return err
	}

	keys, err := c.Doit.GetStringSlice(c.NS, blcli.ArgSSHKeys)
	if err != nil {
		return err
	}

	tagName, err := c.Doit.GetString(c.NS, blcli.ArgTagName)
	if err != nil {
		return err
	}

	vpcID, err := c.Doit.GetInt(c.NS, blcli.ArgVPCID)
	if err != nil {
		return err
	}

	tagNames, err := c.Doit.GetStringSlice(c.NS, blcli.ArgTagNames)
	if err != nil {
		return err
	}
	if len(tagName) > 0 {
		tagNames = append(tagNames, tagName)
	}

	sshKeys := extractSSHKeys(keys)

	userData, err := c.Doit.GetString(c.NS, blcli.ArgUserData)
	if err != nil {
		return err
	}

	volumeList, err := c.Doit.GetStringSlice(c.NS, blcli.ArgVolumeList)
	if err != nil {
		return err
	}
	volumes := extractVolumes(volumeList)

	filename, err := c.Doit.GetString(c.NS, blcli.ArgUserDataFile)
	if err != nil {
		return err
	}

	userData, err = extractUserData(userData, filename)
	if err != nil {
		return err
	}

	imageStr, err := c.Doit.GetString(c.NS, blcli.ArgImage)
	if err != nil {
		return err
	}

	createImage := binarylane.ServerCreateImage{Slug: imageStr}

	i, err := strconv.Atoi(imageStr)
	if err == nil {
		createImage = binarylane.ServerCreateImage{ID: i}
	}

	wait, err := c.Doit.GetBool(c.NS, blcli.ArgCommandWait)
	if err != nil {
		return err
	}

	ds := c.Servers()

	var wg sync.WaitGroup
	var createdList bl.Servers
	errs := make(chan error, len(c.Args))
	for _, name := range c.Args {
		dcr := &binarylane.ServerCreateRequest{
			Name:              name,
			Region:            region,
			Size:              size,
			Image:             createImage,
			Volumes:           volumes,
			Backups:           backups,
			IPv6:              ipv6,
			PrivateNetworking: privateNetworking,
			Monitoring:        monitoring,
			SSHKeys:           sshKeys,
			UserData:          userData,
			VPCID:             vpcID,
			Tags:              tagNames,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			d, err := ds.Create(dcr, wait)
			if err != nil {
				errs <- err
				return
			}

			createdList = append(createdList, *d)
		}()
	}

	wg.Wait()
	close(errs)

	item := &displayers.Server{Servers: createdList}

	for err := range errs {
		if err != nil {
			return err
		}
	}
	c.Display(item)

	return nil
}

// RunServerTag adds a tag to a server.
func RunServerTag(c *CmdConfig) error {
	ds := c.Servers()
	ts := c.Tags()

	if len(c.Args) < 1 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	tag, err := c.Doit.GetString(c.NS, blcli.ArgTagName)
	if err != nil {
		return err
	}

	fn := func(ids []int) error {
		trr := &binarylane.TagResourcesRequest{}
		for _, id := range ids {
			r := binarylane.Resource{
				ID:   strconv.Itoa(id),
				Type: binarylane.ServerResourceType,
			}
			trr.Resources = append(trr.Resources, r)
		}

		return ts.TagResources(tag, trr)
	}

	return matchServers(c.Args, ds, fn)
}

// RunServerUntag untags a server.
func RunServerUntag(c *CmdConfig) error {
	ds := c.Servers()
	ts := c.Tags()

	if len(c.Args) < 1 {
		return blcli.NewMissingArgsErr(c.NS)
	}

	serverIDStrs := c.Args

	tagNames, err := c.Doit.GetStringSlice(c.NS, blcli.ArgTagName)
	if err != nil {
		return err
	}

	fn := func(ids []int) error {
		urr := &binarylane.UntagResourcesRequest{}

		for _, id := range ids {
			for _, tagName := range tagNames {
				r := binarylane.Resource{
					ID:   strconv.Itoa(id),
					Type: binarylane.ServerResourceType,
				}

				urr.Resources = append(urr.Resources, r)

				err := ts.UntagResources(tagName, urr)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	return matchServers(serverIDStrs, ds, fn)
}

func extractSSHKeys(keys []string) []binarylane.ServerCreateSSHKey {
	sshKeys := []binarylane.ServerCreateSSHKey{}

	for _, k := range keys {
		if i, err := strconv.Atoi(k); err == nil {
			if i > 0 {
				sshKeys = append(sshKeys, binarylane.ServerCreateSSHKey{ID: i})
			}
			continue
		}

		if k != "" {
			sshKeys = append(sshKeys, binarylane.ServerCreateSSHKey{Fingerprint: k})
		}
	}

	return sshKeys
}

func extractUserData(userData, filename string) (string, error) {
	if userData == "" && filename != "" {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", err
		}
		userData = string(data)
	}

	return userData, nil
}

func extractVolumes(volumeList []string) []binarylane.ServerCreateVolume {
	var volumes []binarylane.ServerCreateVolume

	for _, v := range volumeList {
		var req binarylane.ServerCreateVolume
		if looksLikeUUID(v) {
			req.ID = v
		} else {
			req.Name = v
		}
		volumes = append(volumes, req)
	}

	return volumes
}

func allInt(in []string) ([]int, error) {
	out := []int{}
	seen := map[string]bool{}

	for _, i := range in {
		if seen[i] {
			continue
		}

		seen[i] = true

		id, err := strconv.Atoi(i)
		if err != nil {
			return nil, fmt.Errorf("%s is not an int", i)
		}
		out = append(out, id)
	}
	return out, nil
}

// RunServerDelete destroy a server by id.
func RunServerDelete(c *CmdConfig) error {
	ds := c.Servers()

	force, err := c.Doit.GetBool(c.NS, blcli.ArgForce)
	if err != nil {
		return err
	}

	tagName, err := c.Doit.GetString(c.NS, blcli.ArgTagName)
	if err != nil {
		return err
	}

	if len(c.Args) < 1 && tagName == "" {
		return blcli.NewMissingArgsErr(c.NS)
	} else if len(c.Args) > 0 && tagName != "" {
		return fmt.Errorf("Please specify Server identifier or a tag name.")
	} else if tagName != "" {
		// Collect affected Server IDs to show in confirmation message.
		var affectedIDs string
		list, err := ds.ListByTag(tagName)
		if err != nil {
			return err
		}
		if len(list) == 0 {
			fmt.Fprintf(c.Out, "Nothing to delete: no Servers are using the \"%s\" tag\n", tagName)
			return nil
		}
		ids := make([]string, 0, len(list))
		for _, server := range list {
			ids = append(ids, strconv.Itoa(server.ID))
		}
		affectedIDs = strings.Join(ids, " ")
		resourceType := "Server"
		if len(list) > 1 {
			resourceType = "Servers"
		}

		if force || AskForConfirm(fmt.Sprintf("delete %d %s tagged \"%s\"? [affected %s: %s]", len(list), resourceType, tagName, resourceType, affectedIDs)) == nil {
			return ds.DeleteByTag(tagName)
		}
		return fmt.Errorf("Operation aborted.")
	}

	if force || AskForConfirmDelete("Server", len(c.Args)) == nil {

		fn := func(ids []int) error {
			for _, id := range ids {
				if err := ds.Delete(id); err != nil {
					return fmt.Errorf("Unable to delete Server %d: %v", id, err)
				}
			}
			return nil
		}
		return matchServers(c.Args, ds, fn)
	}
	return fmt.Errorf("Operation aborted.")
}

type matchServersFn func(ids []int) error

func matchServers(ids []string, ds bl.ServersService, fn matchServersFn) error {
	if extractedIDs, err := allInt(ids); err == nil {
		return fn(extractedIDs)
	}

	sum, err := buildServerSummary(ds)
	if err != nil {
		return err
	}

	matchedMap := map[int]bool{}
	for _, idStr := range ids {
		count := sum.count[idStr]
		if count > 1 {
			return fmt.Errorf("There are %d Servers with the name %q; please provide a specific Server ID. [%s]",
				count, idStr, strings.Join(sum.byName[idStr], ", "))
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			id, ok := sum.byID[idStr]
			if !ok {
				return fmt.Errorf("Server with the name %q could not be found.", idStr)
			}

			matchedMap[id] = true
			continue
		}

		matchedMap[id] = true
	}

	var extractedIDs []int
	for id := range matchedMap {
		extractedIDs = append(extractedIDs, id)
	}

	sort.Ints(extractedIDs)
	return fn(extractedIDs)
}

// RunServerGet returns a server.
func RunServerGet(c *CmdConfig) error {
	err := ensureOneArg(c)
	if err != nil {
		return err
	}

	getTemplate, err := c.Doit.GetString(c.NS, blcli.ArgTemplate)
	if err != nil {
		return err
	}

	ds := c.Servers()
	fn := func(ids []int) error {
		for _, id := range ids {
			d, err := ds.Get(id)
			if err != nil {
				return err
			}

			item := &displayers.Server{Servers: bl.Servers{*d}}

			if getTemplate != "" {
				t := template.New("Get template")
				t, err = t.Parse(getTemplate)
				if err != nil {
					return err
				}
				return t.Execute(c.Out, d)
			}
			return c.Display(item)
		}
		return nil
	}
	return matchServers(c.Args, ds, fn)

}

// RunServerKernels returns a list of available kernels for a server.
func RunServerKernels(c *CmdConfig) error {

	ds := c.Servers()
	id, err := getServerIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Kernels(id)
	if err != nil {
		return err
	}

	item := &displayers.Kernel{Kernels: list}
	return c.Display(item)
}

// RunServerList returns a list of servers.
func RunServerList(c *CmdConfig) error {

	ds := c.Servers()

	region, err := c.Doit.GetString(c.NS, blcli.ArgRegionSlug)
	if err != nil {
		return err
	}

	tagName, err := c.Doit.GetString(c.NS, blcli.ArgTagName)
	if err != nil {
		return err
	}

	matches := []glob.Glob{}
	for _, globStr := range c.Args {
		g, err := glob.Compile(globStr)
		if err != nil {
			return fmt.Errorf("Unknown glob %q", globStr)
		}

		matches = append(matches, g)
	}

	var matchedList bl.Servers

	var list bl.Servers
	if tagName == "" {
		list, err = ds.List()
	} else {
		list, err = ds.ListByTag(tagName)
	}
	if err != nil {
		return err
	}

	for _, server := range list {
		var skip = true
		if len(matches) == 0 {
			skip = false
		} else {
			for _, m := range matches {
				if m.Match(server.Name) {
					skip = false
				}
			}
		}

		if !skip && region != "" {
			if region != server.Region.Slug {
				skip = true
			}
		}

		if !skip {
			matchedList = append(matchedList, server)
		}
	}

	item := &displayers.Server{Servers: matchedList}
	return c.Display(item)
}

// RunServerNeighbors returns a list of server neighbors.
func RunServerNeighbors(c *CmdConfig) error {

	ds := c.Servers()

	id, err := getServerIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Neighbors(id)
	if err != nil {
		return err
	}

	item := &displayers.Server{Servers: list}
	return c.Display(item)
}

// RunServerSnapshots returns a list of available snapshots for a server.
func RunServerSnapshots(c *CmdConfig) error {

	ds := c.Servers()
	id, err := getServerIDArg(c.NS, c.Args)
	if err != nil {
		return err
	}

	list, err := ds.Snapshots(id)
	if err != nil {
		return err
	}

	item := &displayers.Image{Images: list}
	return c.Display(item)
}

func getServerIDArg(ns string, args []string) (int, error) {
	if len(args) != 1 {
		return 0, blcli.NewMissingArgsErr(ns)
	}

	return strconv.Atoi(args[0])
}

type serverSummary struct {
	count  map[string]int
	byID   map[string]int
	byName map[string][]string
}

func buildServerSummary(ds bl.ServersService) (*serverSummary, error) {
	list, err := ds.List()
	if err != nil {
		return nil, err
	}

	var sum serverSummary

	sum.count = map[string]int{}
	sum.byID = map[string]int{}
	sum.byName = map[string][]string{}
	for _, d := range list {
		sum.count[d.Name]++
		sum.byID[d.Name] = d.ID
		sum.byName[d.Name] = append(sum.byName[d.Name], strconv.Itoa(d.ID))
	}

	return &sum, nil
}
