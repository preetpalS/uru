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
		Usage:   "ls",
		HelpMsg: "list all ruby installations",
		Eg:      `ls`}
}

// List all rubies registered with uru, identifying the currently active ruby
func List(ctx *env.Context) {
	if len(ctx.Rubies) == 0 {
		fmt.Println("---> No rubies registered with uru")
		return
	}

	tag, _, err := env.CurrentRubyInfo(ctx)
	if err != nil {
		fmt.Printf("---> Unable to list rubies; try again\n")
		os.Exit(1)
	}

	var me, desc string
	for t, v := range ctx.Rubies {
		if t == tag {
			me = `=>`
		} else {
			me = "  "
		}
		desc = v.Description
		if len(desc) > 65 {
			desc = fmt.Sprintf("%.65s...", desc)
		}
		fmt.Printf(" %s %10.10s: %s\n", me, v.TagLabel, desc)
	}
}
