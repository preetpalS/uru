// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["retag"] = Command{
		Name:    "retag",
		Aliases: []string{"retag", "tag"},
		Usage:   "admin retag CURRENT NEW",
		HelpMsg: "retag CURRENT tag value to NEW",
		Eg:      `admin retag 200p197 200p197-x64`}
}

func adminRetag(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	if len(cmdArgs) != 2 {
		fmt.Println("[ERROR] must specify both CURRENT and NEW tag labels")
		os.Exit(1)
	}

	oldLabel, newLabel := cmdArgs[0], cmdArgs[1]

	for _, ri := range ctx.Registry.Rubies {
		if newLabel == ri.TagLabel {
			fmt.Printf("---> `%s` collides with an existing registered ruby\n", newLabel)
			os.Exit(1)
		}
	}

	tags, err := env.TagLabelToTag(ctx, oldLabel)
	if err != nil {
		fmt.Printf("---> unable to find registered ruby matching `%s`\n", oldLabel)
		os.Exit(1)
	}

	tag := ``
	if len(tags) == 1 {
		// XXX less convoluted way to get the key of a 1 element map?
		for t, _ := range tags {
			tag = t
			break
		}
	} else {
		// multiple rubies match the given tag label, ask the user for the
		// correct one.
		tag, err = env.SelectRubyFromList(tags, oldLabel, `retag`)
		if err != nil {
			os.Exit(1)
		}
	}

	rb := ctx.Registry.Rubies[tag]
	origLabel := rb.TagLabel

	rb.TagLabel = newLabel
	ctx.Registry.Rubies[tag] = rb

	err = env.MarshalRubies(ctx)
	if err != nil {
		fmt.Printf("---> Failed to retag `%s` to `%s`. Try again", origLabel, newLabel)
	}

	fmt.Printf("---> retagged `%s` to `%s`\n", origLabel, newLabel)
}
