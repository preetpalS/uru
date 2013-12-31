// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"errors"
	"fmt"

	"bitbucket.org/jonforums/uru/env"
)

// list gems for all gemset members or just the gems relevant to the given
// `ruby@gemset` name
func gemsetList(ctx *env.Context, ruby, gemset string) (err error) {
	if gemset != `gemset` && gemset != `` {
		return errors.New("[ERROR] unable to list gemset gems. Only project gemsets supported.")
	}

	var str string
	if ruby == `` {
		str = `for all gemset members`
	} else {
		str = fmt.Sprintf("for ruby matching `%s` label", ruby)
	}
	fmt.Printf("---> listing gemset gems %s\n", str)

	return
}
