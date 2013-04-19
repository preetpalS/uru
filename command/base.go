// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"bitbucket.org/jonforums/uru/env"
)

var (
	CommandRegistry  = make(map[string]Command)
	AdminCmdRegistry = make(map[string]Command)
)

type Command struct {
	Name    string
	Aliases []string
	Usage   string
	HelpMsg string
	Eg      string
}

func rubyExec(ctx *env.Context) (err error) {
	// TODO error check for empty PATH string
	curPath := os.Getenv(`PATH`)
	curGemHome := os.Getenv(`GEM_HOME`)

	for rb, info := range ctx.Rubies {
		fmt.Printf("%s\n\n", info.Description)

		pth, err := env.PathListForTag(ctx, rb)
		if err != nil {
			fmt.Println("Unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			break
		}

		// set env vars in this process so they'll be injected into the child process
		err = os.Setenv(`PATH`, strings.Join(pth, string(os.PathListSeparator)))
		if err != nil {
			fmt.Println("Unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			break
		}
		err = os.Setenv(`GEM_HOME`, info.GemHome)
		if err != nil {
			fmt.Println("Unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			break
		}

		// run the command in a child process and reflect the child's stdout
		cmd := ctx.Cmd()
		if cmd == `ruby` {
			// invoke correct ruby executable
			cmd = info.Exe
		}
		out, err := exec.Command(cmd, ctx.CmdArgs()...).Output()
		if err != nil {
			fmt.Printf("---> Unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
		} else {
			fmt.Printf("%s\n", out)
		}
	}

	// switch back to the original ruby
	os.Setenv(`PATH`, curPath)
	if curGemHome != `` {
		os.Setenv(`GEM_HOME`, curGemHome)
	}

	return
}
