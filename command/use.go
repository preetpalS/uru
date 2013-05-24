// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"errors"
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/env"
	"bitbucket.org/jonforums/uru/exec"
)

func Use(ctx *env.Context, msg string) {
	cmd := ctx.Cmd()

	// use .ruby-version file contents to select which ruby to use
	// credit: thanks to Luis Lavena for the idea
	var tags map[string]env.Ruby
	var err error
	switch cmd {
	case `.`:
		tags, err = useRubyVersionFile(ctx)
		if err != nil {
			fmt.Println("---> someday soon I'll understand the `.ruby-version` file")
			os.Exit(1)
		}
	default:
		tags, err = env.TagLabelToTag(ctx, cmd)
		if err != nil {
			fmt.Printf("---> unable to find registered ruby matching `%s`\n", cmd)
			os.Exit(1)
		}
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
		tag, err = env.SelectRubyFromList(tags, cmd, `use`)
		if err != nil {
			os.Exit(1)
		}
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

// TODO implement
func useRubyVersionFile(ctx *env.Context) (tags map[string]env.Ruby, err error) {
	return nil, errors.New("not implemented")
}
