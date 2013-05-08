// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"

	"bitbucket.org/jonforums/uru/env"
)

// generate no help information as it adds no value and clutters the CLI

func Version(ctx *env.Context) {
	fmt.Printf("%s v%s\n", env.AppName, env.AppVersion)
}
