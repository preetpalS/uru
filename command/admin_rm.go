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

	tag, err := env.TagLabelToTag(ctx, ctx.CmdArgs()[0])
	if err != nil {
		fmt.Printf("---> Skipping. Unable to find ruby registered as `%s`\n", tag)
		return
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
