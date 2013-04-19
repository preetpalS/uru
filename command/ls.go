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
		fmt.Println("---> No rubies registered with uru.")
		return
	}

	tag, _, err := env.CurrentRubyInfo()
	if err != nil {
		fmt.Printf("---> Unable to list rubies; try again.")
		os.Exit(1)
	}

	var me, desc string
	for k, v := range ctx.Rubies {
		if k == tag {
			me = `=>`
		} else {
			me = "  "
		}
		desc = v.Description
		if len(desc) >= 67 {
			desc = fmt.Sprintf("%.67s...", desc)
		}
		fmt.Printf(" %s %7.7s: %s\n", me, k, desc)
	}
}
