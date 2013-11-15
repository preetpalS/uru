// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	AppName    = `uru`
	AppVersion = `0.6.4`
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

// SelectRubyFromList presents a list of registered rubies and asks the user
// to select one. It returns the identifying tag for the selected ruby, or an
// error if unable to get the users selection.
func SelectRubyFromList(tags RubyMap, label, verb string) (tag string, err error) {
	var i, choice uint8
	choices := make(map[uint8]string)
	indent := fmt.Sprintf("%19.19s", ``)

	fmt.Printf("---> these rubies match your `%s` tag:\n\n", label)
	for t, ri := range tags {
		i++
		choices[i] = t
		fmt.Printf(" [%d] %-12.12s: %s\n%sHome: %s\n",
			i,
			ri.TagLabel,
			ri.Description,
			indent,
			ri.Home)
	}
	fmt.Printf("\nselect [1]-[%d] to %s that specific ruby (0 to exit) [0]: ", i, verb)
	_, err = fmt.Scanln(&choice)
	if err != nil || choice == 0 || choice > i {
		return ``, errors.New("error: unable to get users ruby selection response")
	} else {
		tag = choices[choice]
	}

	return
}
