// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"bitbucket.org/jonforums/uru/env"
)

type rbVersionFunc func(ctx *env.Context, dir string) (tags map[string]env.Ruby, err error)

func useRubyVersionFile(ctx *env.Context, verFunc rbVersionFunc) (tags map[string]env.Ruby, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		return nil, err
	}

	userHome := ``
	if runtime.GOOS == `windows` {
		userHome = os.Getenv(`USERPROFILE`)
	} else {
		userHome = os.Getenv(`HOME`)
	}
	if userHome == `` {
		return nil, err
	}
	userHome, err = filepath.Abs(userHome)
	if err != nil {
		return nil, err
	}

	atRoot := false
	for !atRoot {
		// TODO stdlib have anything more robust than string compare?
		if absCwd == userHome {
			absCwd = filepath.Dir(absCwd)
			continue
		}

		tags, err = verFunc(ctx, absCwd)
		if err == nil {
			return
		}

		absCwd = filepath.Dir(absCwd)
		err = os.Chdir(absCwd)
		if err != nil {
			return nil, err
		}

		var path string
		if runtime.GOOS == `windows` {
			path = strings.Split(absCwd, `:`)[1]
		} else {
			path = absCwd
		}
		// have walked back up to root so perform last check before fallback
		// check for .ruby-version in $HOME/%UserProfile%
		// TODO hoist further up to prevent double stat if starting at root
		if strings.HasPrefix(path, string(os.PathSeparator)) &&
			strings.HasSuffix(path, string(os.PathSeparator)) {
			atRoot = true

			tags, err = verFunc(ctx, absCwd)
			if err == nil {
				return
			}

			err = os.Chdir(userHome)
			if err != nil {
				return nil, err
			}

			tags, err = verFunc(ctx, userHome)
			if err == nil {
				return
			}
		}
	}

	return nil, errors.New("unable to find a .ruby-version file")
}

func versionator(ctx *env.Context, dir string) (tags map[string]env.Ruby, err error) {
	var path string
	if strings.HasSuffix(dir, string(os.PathSeparator)) {
		path = fmt.Sprintf("%s.ruby-version", dir)
	} else {
		path = fmt.Sprintf("%s%s.ruby-version", dir, string(os.PathSeparator))
	}
	log.Printf("[DEBUG] checking for `%s`\n", path)

	if _, err = os.Stat(path); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	b = bytes.Trim(b, " \r\n")
	b = bytes.ToLower(b)
	rbVer := bytes.NewBuffer(b).String()
	log.Printf("[DEBUG] .ruby-version data: %s\n", rbVer)

	return env.TagLabelToTag(ctx, rbVer)
}
