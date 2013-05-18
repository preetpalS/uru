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
		Usage:   "add RUBY_DIR [--tag TAG] | system",
		HelpMsg: "register an existing ruby installation",
		Eg:      `add C:\ruby200\bin`}
}

func adminAdd(ctx *env.Context) {
	argsLen := 0
	cmdArgs := ctx.CmdArgs()

	if argsLen = len(cmdArgs); argsLen == 0 {
		fmt.Println("[ERROR] must specify a ruby bindir or `system`.")
		os.Exit(1)
	}

	loc := cmdArgs[0]

	tagAlias := ``
	for i, v := range ctx.CmdArgs() {
		if v == `--tag` {
			if i < argsLen {
				tagAlias = cmdArgs[i+1]
				break
			} else {
				fmt.Println("[ERROR] invalid `admin add --tag TAG` invocation.")
				os.Exit(1)
			}
		}
	}

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
			fmt.Printf("---> Unable to find a known ruby at `%s`\n", loc)
			return
		}
	}

	tag, rbInfo, err := env.RubyInfo(ctx, rbPath)
	if err != nil {
		fmt.Printf("---> Unable to register `%s` due to missing ruby info\n", rbPath)
		return
	}

	// set tag alias if given
	if tagAlias != `` {
		rbInfo.TagLabel = tagAlias
	}

	// assume the vast majority of windows users install gems into the ruby
	// installation; clear GEM_HOME value source to prevent persisting a
	// GEM_HOME value for the ruby being registered.
	// XXX potential usage bug
	if runtime.GOOS == `windows` {
		rbInfo.GemHome = ``
	}

	// patch metadata if adding a ruby with the same default tag label as an
	// existing registered ruby.
	for t, i := range ctx.Rubies {
		// default tag labels are the same but tag (description/home hash) is different
		if rbInfo.TagLabel == i.TagLabel && tag != t {
			if tagAlias != `` && len(tagAlias) <= 10 {
				rbInfo.TagLabel = tagAlias
			} else {
				fmt.Printf(`
---> So sorry, but I'm not able to register the following ruby
--->
--->   %s
--->
---> because its tag label conflicts with a previously registered
---> ruby. Please re-register the ruby with a unique tag alias by
---> running the following command:
--->
--->   %s admin add RUBY_DIR --tag TAG
--->
---> where TAG is 10 characters or less.`, loc, env.AppName)
				os.Exit(1)
			}
		}
	}

	// patch metadata if adding a system ruby
	if loc == `system` {
		rbInfo.TagLabel = `system`
		rbInfo.GemHome = os.Getenv(`GEM_HOME`) // user configured value or empty
	}

	// TODO allow overwriting or force rm/add cycle?
	if _, ok := ctx.Rubies[tag]; ok {
		fmt.Printf("---> Skipping. `%s` is already registered\n", rbPath)
		return
	}

	ctx.Rubies[tag] = rbInfo

	// persist the new and existing registered rubies to the filesystem
	err = env.MarshalRubies(ctx)
	if err != nil {
		fmt.Printf("---> Failed to register `%s`, try again\n", rbPath)
	} else {
		fmt.Printf("---> Registered %s at `%s` as `%s`\n", rbInfo.Exe, rbInfo.Home, rbInfo.TagLabel)
	}
}
