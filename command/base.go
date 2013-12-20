// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

	for tag, info := range ctx.Registry.Rubies {
		fmt.Printf("%s\n\n", info.Description)

		pth, err := env.PathListForTag(ctx, tag)
		if err != nil {
			fmt.Printf("[ERROR] getting path list, unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			break
		}

		// set env vars in current process to be inherited by the child process
		err = os.Setenv(`PATH`, strings.Join(pth, string(os.PathListSeparator)))
		if err != nil {
			fmt.Printf("[ERROR] setting PATH, unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))
			break
		}
		if info.GemHome != `` {
			// XXX oddly, GEM_HOME must be set in current process so that users .gemrc
			// is consulted. Setting os/exec's `Command.Env` causes users .gemrc to
			// be ignored.
			err = os.Setenv(`GEM_HOME`, info.GemHome)
			if err != nil {
				fmt.Printf("[ERROR] setting GEM_HOME, unable to run `%s %s`\n\n", ctx.Cmd(),
					strings.Join(ctx.CmdArgs(), " "))
				break
			}
		}

		// run the command in a child process and capture stdout/stderr
		cmd := ctx.Cmd()
		if runtime.GOOS == `windows` || cmd == `ruby` {
			// on windows, bypass .bat wrappers; always select correct ruby exe
			cmd = info.Exe
		}
		cmdArgs := ctx.CmdArgs()
		if runtime.GOOS == `windows` && ctx.Cmd() == `gem` {
			// on windows, bypass gem.bat wrapper; always run gem via ruby exe
			cmdArgs = append([]string{filepath.Join(info.Home, `gem`)}, cmdArgs...)
		}
		log.Printf("[DEBUG] === exec.Command args ===\n  cmd: %s\n  cmdArgs: %v\n",
			cmd, cmdArgs)

		runner := exec.Command(cmd, cmdArgs...)
		// TODO while allowing communication with the child process, intermediate
		// output from the child process is not currently displayed making this
		// capability almost useless.
		runner.Stdin = os.Stdin
		out, err := runner.CombinedOutput()
		if err != nil {
			fmt.Printf("---> unable to run `%s %s`\n\n", ctx.Cmd(),
				strings.Join(ctx.CmdArgs(), " "))

			log.Printf("--- returned error message ---\n%s\n\n", err.Error())
			log.Printf("--- combined child output ---\n%s\n", out)
		} else {
			fmt.Printf("%s\n", out)
		}
	}

	// revert to the original ruby
	os.Setenv(`PATH`, curPath)
	if curGemHome != `` {
		os.Setenv(`GEM_HOME`, curGemHome)
	}

	return
}
