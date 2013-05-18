// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/env"
	"bitbucket.org/jonforums/uru/exec"
)

func Use(ctx *env.Context, msg string) {
	cmd := ctx.Cmd()

	tag, err := env.TagLabelToTag(ctx, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "---> unable to find ruby matching `%s`\n", cmd)
		os.Exit(1)
	}

	newRb := ctx.Rubies[tag]

	newPath, err := env.PathListForTag(ctx, tag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "---> unable to use ruby tagged as `%s`\n", tag)
		os.Exit(1)
	}

	// create and execute the environment changing runner script
	exec.CreateScript(ctx, &newPath, newRb.GemHome)
	exec.ExecScript(ctx)

	switch msg {
	case ``:
		tagAlias := ``
		if newRb.TagLabel != `` {
			tagAlias = fmt.Sprintf("tagged as `%s`", newRb.TagLabel)
		}
		fmt.Printf("---> Now using %s %s %s\n", newRb.Exe, newRb.ID, tagAlias)
	case ` `:
		break
	default:
		fmt.Printf(msg)
	}
}
