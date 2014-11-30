// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/jonforums/uru/env"
	"bitbucket.org/jonforums/uru/exec"
)

func useNil(ctx *env.Context) (err error) {
	envPath := os.Getenv(`PATH`)
	if envPath == `` {
		return errors.New("unable to get PATH env var value")
	}

	// uru-free environment as uru's Canary is not part of PATH
	if strings.Index(envPath, env.Canary) == -1 {
		return
	}

	// remove uru's effect on current PATH
	fmt.Println("---> removing non-system ruby from current environment")
	curPath := strings.Split(envPath, env.CanaryToken)
	newPath := strings.Split(curPath[1], string(os.PathListSeparator))

	// TODO handle pre-existing "system" GEM_HOME via URU_ORIGINAL_GEM_HOME envar
	exec.CreateSwitcherScript(ctx, &newPath, "")

	return
}
