// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"

	"bitbucket.org/jonforums/uru/env"
)

func gemsetRemove(ctx *env.Context) (err error) {
	rv, err := env.UIYesConfirm(`Delete all gems for all rubies of the gemset?`)
	if err != nil {
		// TODO implement
	}

	if rv == `Y` {
		fmt.Println("---> removing all gemset gems")
	}

	return
}
