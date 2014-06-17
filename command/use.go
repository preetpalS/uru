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

	// use .ruby-version file contents to select which ruby to activate
	var tags env.RubyMap
	var err error
	switch cmd {
	case `auto`:
		tags, err = useRubyVersionFile(ctx, versionator)
		if err != nil {
			fmt.Println("---> unable to find or process a `.ruby-version` file")
			os.Exit(1)
		}
	case `nil`:
		useNil(ctx)
		os.Exit(0)
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

	newRb := ctx.Registry.Rubies[tag]

	newPath, err := env.PathListForTag(ctx, tag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "---> unable to use ruby tagged as `%s`\n", tag)
		os.Exit(1)
	}

	// create the environment switcher script
	exec.CreateSwitcherScript(ctx, &newPath, newRb.GemHome)

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
