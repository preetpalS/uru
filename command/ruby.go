// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bitbucket.org/jonforums/uru/env"
)

func init() {
	CommandRegistry["ruby"] = Command{
		Name:    "ruby",
		Aliases: []string{"ruby", "rb"},
		Usage:   "ruby ARGS ...",
		HelpMsg: "run a ruby command with all registered rubies",
		Eg:      `ruby -e "puts RUBY_VERSION"`}
}

func Ruby(ctx *env.Context) {
	ctx.SetCmd(`ruby`)
	err := rubyExec(ctx)
	if err != nil {
		// TODO implement me
	}
}
