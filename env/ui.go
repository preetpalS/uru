// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"fmt"
	"regexp"
)

const (
	AppName    = "uru"
	AppVersion = "0.2.0"
)

var (
	yResp *regexp.Regexp
)

func init() {
	var err error
	yResp, err = regexp.Compile(`\A(?i)Y`)
	if err != nil {
		panic("unable to compile UI `yes` response regexp")
	}
}

// UIYesConfirm asks the user a question and returns Y or N with
// Y being the default.
func UIYesConfirm(prompt string) (resp string, err error) {
	resp = "Y"
	fmt.Printf("%s [Yn] ", prompt)
	_, err = fmt.Scanln(&resp)
	if yResp.MatchString(resp) {
		resp, err = "Y", nil
	} else {
		resp = "N"
	}

	return
}
