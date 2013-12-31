// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"errors"
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["gemset"] = Command{
		Name:    "gemset",
		Aliases: []string{"gemset", "gs"},
		Usage:   "admin gemset init NAME... | ls [NAME] | rm",
		HelpMsg: "administer gemset installations",
		Eg:      `admin gemset init`}
}

func adminGemset(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	argsLen := len(cmdArgs)
	if argsLen == 0 {
		fmt.Println("[ERROR] must specify a gemset operation.")
		os.Exit(1)
	}

	var rubyName, gemsetName string
	var err error

	switch subCmd := cmdArgs[0]; subCmd {
	case `init`:
		if argsLen < 2 || argsLen > 20 { // artificial upper limit
			fmt.Println("[ERROR] invalid `admin gemset init NAME...` invocation.")
			os.Exit(1)
		}

		for _, v := range cmdArgs[1:] {
			if rubyName, gemsetName, err = parseGemsetName(v); err != nil {
				fmt.Println("---> invalid `admin gemset init NAME...` invocation.")
				continue
			}
			if err = gemsetInit(ctx, rubyName, gemsetName); err != nil {
				fmt.Println(err)
				continue
			}
		}
	case `ls`:
		if !(argsLen == 1 || argsLen == 2) {
			fmt.Println("[ERROR] invalid `admin gemset ls [NAME]` invocation.")
			os.Exit(1)
		}

		switch argsLen {
		case 1:
			rubyName = ``
		default:
			if rubyName, gemsetName, err = parseGemsetName(cmdArgs[1]); err != nil {
				fmt.Println("[ERROR] invalid `admin gemset ls [NAME]` invocation.")
				os.Exit(1)
			}
		}
		if err = gemsetList(ctx, rubyName, gemsetName); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case `rm`:
		if err := gemsetRemove(ctx); err != nil {
			fmt.Println("[ERROR] unable to remove the gemset.")
			os.Exit(1)
		}
	default:
		fmt.Printf("[ERROR] I don't understand the `%s` gemset sub-command\n\n", subCmd)
		ctx.SetCmdAndArgs(``, nil)
		Help(ctx)
	}
}

// create a skeleton gemset directory structure in the current directory with
// the following layout
//
//    .gem/$ENGINE/$RUBY_LIBRARY_VERSION
func gemsetInit(ctx *env.Context, ruby, gemset string) (err error) {
	if gemset != `gemset` {
		return errors.New("---> unable to initialize gemset. Only project gemsets supported.")
	}

	fmt.Printf("---> initializing project gemset for ruby matching `%s` label\n", ruby)

	return
}
