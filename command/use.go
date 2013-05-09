// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/jonforums/uru/env"
	"bitbucket.org/jonforums/uru/exec"
)

func Use(ctx *env.Context, msg string) {
	// canonicalize user provided tag from cmd line
	cmd := ctx.Cmd()
	tag := ``
	for k, _ := range ctx.Rubies {
		if strings.Contains(k, cmd) {
			tag = k
			break
		}
	}
	if tag == `` {
		fmt.Fprintf(os.Stderr, "---> unable to find ruby matching `%s`\n", cmd)
		os.Exit(1)
	}

	newRb := ctx.Rubies[tag]

	pth, err := env.PathListForTag(ctx, tag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "---> unable to use ruby tagged as `%s`\n", tag)
		os.Exit(1)
	}

	// create and execute the environment changing runner script
	exec.CreateScript(ctx, &pth, newRb.GemHome)
	exec.ExecScript(ctx)

	switch msg {
	case ``:
		fmt.Printf("---> Now using %s %s\n", newRb.Exe, newRb.ID)
	case ` `:
		break
	default:
		fmt.Printf(msg)
	}
}
