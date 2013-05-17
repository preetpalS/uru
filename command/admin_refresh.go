// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["refresh"] = Command{
		Name:    "refresh",
		Aliases: nil,
		Usage:   "refresh",
		HelpMsg: "refresh all registered rubies",
		Eg:      `refresh`}
}

func adminRefresh(ctx *env.Context) {

	freshRubies := make(env.RubyRegistry, 4)

	for tag, info := range ctx.Rubies {
		rb := filepath.Join(info.Home, info.Exe)

		newTag, freshInfo, err := env.RubyInfo(rb)
		if err != nil {
			fmt.Println("---> Unable to determine ruby metadata while refreshing")
			os.Exit(1)
		}

		// XXX assume windows users always install gems into the ruby installation
		// so GEM_HOME is always empty except in the case of a system ruby in which
		// the GEM_HOME env var was active at system ruby registration.
		if runtime.GOOS == `windows` {
			freshInfo.GemHome = ``
		}
		// patch up freshened ruby GEM_HOME with registered system ruby GEM_HOME as
		// `RubyInfo` only generates a default value.
		if tag == `system` {
			newTag = `system`
			freshInfo.GemHome = ctx.Rubies[tag].GemHome
		}

		fmt.Printf("---> Refreshing %s tagged as %s\n", info.Exe, tag)
		freshRubies[newTag] = freshInfo
	}

	log.Printf("[DEBUG] === fresh ruby metadata ===\n%+v\n", freshRubies)
	ctx.Rubies = freshRubies

	err := env.MarshalRubies(ctx)
	if err != nil {
		fmt.Println("---> unable to persist refreshed ruby metadata")
		os.Exit(1)
	}
}
