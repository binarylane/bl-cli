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
	"git.mammoth.com.au/github/bl-cli/commands/displayers"
	"github.com/spf13/cobra"
)

// Size creates the size commands hierarchy.
func Size() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "size",
			Short: "List available Server sizes",
			Long:  "The subcommands of `doctl compute size` retrieve information about Server sizes.",
		},
	}

	sizeDesc := `List the slug identifier, RAM, VCPU count, disk size, and pricing details for each Server size.

Use the slugs displayed by this command to specify the type of Server in other commands.
`
	CmdBuilder(cmd, RunSizeList, "list", "List available Server sizes", sizeDesc,
		Writer, aliasOpt("ls"), displayerType(&displayers.Size{}))

	return cmd
}

// RunSizeList all sizes.
func RunSizeList(c *CmdConfig) error {
	sizes := c.Sizes()

	list, err := sizes.List()
	if err != nil {
		return err
	}

	item := &displayers.Size{Sizes: list}
	return c.Display(item)
}
