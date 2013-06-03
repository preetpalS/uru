// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"
	"sort"

	"bitbucket.org/jonforums/uru/env"
)

type tagInfo struct {
	Tag      string
	TagLabel string
}

// tagInfoSorter sorts slices of tagInfo structs by implementing sort.Interface by
// providing Len, Swap, and Less
type tagInfoSorter struct {
	Tags []tagInfo
}

func (s *tagInfoSorter) Len() int {
	return len(s.Tags)
}

func (s *tagInfoSorter) Swap(i, j int) {
	s.Tags[i], s.Tags[j] = s.Tags[j], s.Tags[i]
}

func (s *tagInfoSorter) Less(i, j int) bool {
	return s.Tags[i].TagLabel < s.Tags[j].TagLabel
}

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
		fmt.Printf("---> Unable to list rubies; try again\n")
		os.Exit(1)
	}

	// sort tags by tag labels
	tis := new(tagInfoSorter)
	tis.Tags = []tagInfo{}
	for t, ri := range ctx.Registry.Rubies {
		tis.Tags = append(tis.Tags, tagInfo{Tag: t, TagLabel: ri.TagLabel})
	}
	sort.Sort(tis)

	var me, desc string
	indent := fmt.Sprintf("%17.17s", ``)
	for _, ti := range tis.Tags {
		t := ti.Tag
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
