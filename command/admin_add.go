// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["add"] = Command{
		Name:    "add",
		Aliases: nil,
		Usage:   "add RUBY_DIR | system",
		HelpMsg: "register an existing ruby installation",
		Eg:      `add C:\ruby200\bin`}
}

func adminAdd(ctx *env.Context) {
	if len(ctx.CmdArgs()) == 0 {
		fmt.Println("[ERROR] must specify a ruby bindir, or `system`, to register.")
		os.Exit(1)
	}

	loc := ctx.CmdArgs()[0]

	var rbPath, ext string
	if runtime.GOOS == `windows` {
		ext = ".exe"
	}
	switch loc {
	case `system`:
		var err error
		for _, v := range env.KnownRubies {
			rbPath, err = exec.LookPath(v)
			if err == nil {
				break
			}
		}
	default:
		for _, v := range env.KnownRubies {
			rbPath = filepath.Join(loc, fmt.Sprintf("%s%s", v, ext))
			_, err := os.Stat(rbPath)
			if os.IsNotExist(err) {
				rbPath = ""
				continue
			} else {
				break
			}
		}
		if rbPath == `` {
			fmt.Printf("---> Unable to find a known ruby at %s\n", loc)
			return
		}
	}

	tag, rbInfo, err := env.RubyInfo(rbPath)
	if err != nil {
		fmt.Printf("---> Unable to register %s due to missing ruby info\n", rbPath)
		return
	}

	// assume the vast majority of windows users install gems into the ruby
	// installation; clear GEM_HOME value source to prevent persisting a
	// GEM_HOME value for the ruby being registered.
	// XXX potential usage bug
	if runtime.GOOS == `windows` {
		rbInfo.GemHome = ``
	}

	// patch up if adding a system ruby
	if loc == `system` {
		tag = `system`
		rbInfo.GemHome = os.Getenv(`GEM_HOME`) // user configured value or empty
	}

	// TODO allow overwriting or force rm/add cycle?
	if _, ok := ctx.Rubies[tag]; ok {
		fmt.Printf("---> Skipping. `%s` is already registered.\n", rbPath)
		return
	}

	ctx.Rubies[tag] = rbInfo

	// persist the new and existing registered rubies to the filesystem
	err = env.MarshalRubies(ctx)
	if err != nil {
		fmt.Printf("---> Failed to register `%s`, try again\n", rbPath)
	}

	fmt.Printf("---> Registered %s at `%s` as `%s`\n", rbInfo.Exe, rbInfo.Home, tag)
}
