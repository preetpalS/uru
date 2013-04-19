// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	rbRegex, rbVerRegex *regexp.Regexp
	KnownRubies         []string

	SysRbRegex *regexp.Regexp

	Canary = fmt.Sprintf("%s%s%s", string(os.PathListSeparator),
		string(os.PathListSeparator), string(os.PathListSeparator))
)

type Ruby struct {
	ID          string // ruby version including patch number
	Exe         string // ruby executable name
	Home        string // full path to ruby executable directory
	GemHome     string // full path to a ruby's gem home directory
	Description string // full ruby description
}

func init() {
	var err error
	rbRegex, err = regexp.Compile(`\A(j?ruby)\s+(\d\.\d\.\d(?:\w+)?)`)
	if err != nil {
		panic("unable to compile ruby parsing regexp")
	}

	rbVerRegex, err = regexp.Compile(`\A(\d\.\d\.\d)`)
	if err != nil {
		panic("unable to compile ruby version parsing regexp")
	}

	SysRbRegex, err = regexp.Compile(`\Asys`)
	if err != nil {
		panic("unable to compile ruby parsing regexp")
	}

	KnownRubies = []string{`ruby`, `jruby`}
}

// CurrentRubyInfo returns the tag for the ruby currently in use.
func CurrentRubyInfo() (tag string, info Ruby, err error) {
	envPath := os.Getenv(`PATH`)
	if envPath == `` {
		err = errors.New("Unable to read PATH environment variable")
		return
	}

	if strings.Index(envPath, Canary) != -1 {
		// modified PATH looks like:
		//   GEM_HOME;RUBY_DIR;;;... -or- GEM_HOME:RUBY_DIR:::...
		curRbPath := strings.Split(envPath, string(os.PathListSeparator))[1]
		for _, v := range KnownRubies {
			tstRb := []string{curRbPath, v}
			tag, info, err = RubyInfo(strings.Join(tstRb, string(os.PathSeparator)))
			if err == nil {
				break
			}
		}
	} else {
		tag = `system`
	}

	return
}

// RubyInfo returns information about a specific ruby. It accepts a string with
// either the simple name of the ruby executable, or the absolute path the the
// ruby executable.
func RubyInfo(ruby string) (tag string, info Ruby, err error) {
	rb, err := exec.LookPath(ruby)
	if err != nil {
		return
	}

	info.Home = filepath.Dir(rb)

	c := exec.Command(rb, `--version`)
	b, err := c.Output()
	if err != nil {
		err = errors.New("unable to capture ruby version info")
		return
	}

	info.Description = strings.TrimSpace(string(b))
	res := rbRegex.FindStringSubmatch(info.Description)
	if res != nil {
		info.ID = res[2]
		info.Exe = res[1]
		tag = strings.Replace(info.ID, `.`, ``, -1)
		usrHome := ``
		if runtime.GOOS == `windows` {
			usrHome = os.Getenv(`USERPROFILE`)
		} else {
			usrHome = os.Getenv(`HOME`)
		}
		info.GemHome = filepath.Join(usrHome, `.gem`, info.Exe,
			rbVerRegex.FindStringSubmatch(info.ID)[0])
	} else {
		err = errors.New("unable to parse ruby name and version info")
		return
	}
	log.Printf("[DEBUG] tag: %s, %+v\n", tag, info)

	return
}

// MarshallRubies persists the registered rubies to a JSON formatted file.
func MarshalRubies(ctx *Context) (err error) {
	src := filepath.Join(ctx.Home(), `rubies.json`)
	dst := filepath.Join(ctx.Home(), `rubies.json.bak`)

	_, err = os.Stat(src)
	if os.IsNotExist(err) {
		log.Printf("[DEBUG] %s does not exist; creating\n", src)
		f, e := os.Create(src)
		if e != nil {
			log.Printf("[DEBUG] unable to create new %s\n", src)
			return e
		}
		defer f.Close()
	} else {
		_, e := CopyFile(dst, src)
		if e != nil {
			log.Println("[DEBUG] unable to backup JSON ruby registry")
			return e
		}
	}

	b, err := json.Marshal(ctx.Rubies)
	if err != nil {
		log.Println("[DEBUG] unable to marshall the ruby registry to JSON")
		return
	}

	buf := new(bytes.Buffer)
	err = json.Indent(buf, b, ``, `  `)
	if err != nil {
		log.Println("[DEBUG] unable to format the JSON marshalled ruby registry")
		return
	}

	err = ioutil.WriteFile(src, buf.Bytes(), 0)
	if err != nil {
		os.Remove(src)
		os.Rename(dst, src)
		log.Println("[DEBUG] unable to persist the updated JSON ruby registry")
		return
	}

	os.Remove(dst)
	return
}
