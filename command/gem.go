// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bitbucket.org/jonforums/uru/env"
)

func init() {
	CommandRegistry["gem"] = Command{
		Name:    "gem",
		Aliases: nil,
		Usage:   "gem ARGS ...",
		HelpMsg: "run a gem command with all registered rubies",
		Eg:      `gem install narray`}
}

func Gem(ctx *env.Context) {
	err := rubyExec(ctx)
	if err != nil {
		// TODO implement me
	}
}
