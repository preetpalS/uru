// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	CommandRegistry["ls"] = Command{
		Name:    "ls",
		Aliases: []string{"ls", "list"},
		Usage:   "ls [--verbose]",
		HelpMsg: "list all registered ruby installations",
		Eg:      `ls`}
}

// List all rubies registered with uru, identifying the currently active ruby
func List(ctx *env.Context) {
	if len(ctx.Registry.Rubies) == 0 {
		fmt.Println("---> No rubies registered with uru")
		return
	}

	verbose := false
	for _, v := range ctx.CmdArgs() {
		if v == `--verbose` {
			verbose = true
			break
		}
	}

	tag, _, err := env.CurrentRubyInfo(ctx)
	if err != nil {
		fmt.Println("---> unable to list rubies; try again\n")
		os.Exit(1)
	}

	sortedTags, err := env.SortTagsByTagLabel(ctx)
	if err != nil {
		fmt.Println("---> unable to list sorted rubies; try again")
		os.Exit(1)
	}

	var me, desc string
	indent := fmt.Sprintf("%17.17s", ``)
	for _, t := range sortedTags {
		ri := ctx.Registry.Rubies[t]

		if t == tag {
			me = `=>`
		} else {
			me = "  "
		}

		desc = ri.Description
		if len(desc) > 64 {
			desc = fmt.Sprintf("%.64s...", desc)
		}

		fmt.Printf(" %s %12.12s: %s\n", me, ri.TagLabel, desc)
		if verbose {
			fmt.Printf("%s ID: %s\n%s Home: %s\n%s GemHome: %s\n\n",
				indent, ri.ID, indent, ri.Home, indent, ri.GemHome)
		}
	}
}
