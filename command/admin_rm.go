// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["rm"] = Command{
		Name:    "rm",
		Aliases: []string{"rm", "del"},
		Usage:   "rm TAG",
		HelpMsg: "deregister a ruby installation",
		Eg:      `rm 193p193`}
}

func adminRemove(ctx *env.Context) {
	if len(ctx.CmdArgs()) == 0 {
		fmt.Println("[ERROR] must specify the tag of the ruby to deregister")
		os.Exit(1)
	}

	cmd := ctx.CmdArgs()[0]
	tags, err := env.TagLabelToTag(ctx, cmd)
	if err != nil {
		fmt.Printf("---> unable to find registered ruby matching `%s`\n", cmd)
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
		tag, err = env.SelectRubyFromList(tags, cmd, `deregister`)
		if err != nil {
			os.Exit(1)
		}
	}

	rb := ctx.Rubies[tag]

	resp, err := env.UIYesConfirm(fmt.Sprintf("\nOK to deregister `%s`?", rb.Description))
	if err != nil {
		fmt.Println("---> Unable to understand your response. Try again")
		return
	}
	if resp == `N` {
		return
	}

	delete(ctx.Rubies, tag)

	err = env.MarshalRubies(ctx)
	if err != nil {
		fmt.Printf("---> Failed to remove `%s`. Try again", tag)
	}
}
