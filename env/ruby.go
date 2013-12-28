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

const RubyRegistryVersion = `1.0.0`

var (
	rbRegex, rbVerRegex, rbMajMinRegex, SysRbRegex *regexp.Regexp
	KnownRubies                                    []string

	Canary = fmt.Sprintf("%s%s%s", string(os.PathListSeparator),
		string(os.PathListSeparator), string(os.PathListSeparator))
)

type RubyMap map[string]Ruby

type RubyRegistry struct {
	Version string
	Rubies  RubyMap
}

type MarshalFunc func(ctx *Context) error

type RubyMarshaller struct {
	marshaller MarshalFunc
}

func NewRubyMarshaller(m MarshalFunc) *RubyMarshaller {
	return &RubyMarshaller{marshaller: m}
}

func (m *RubyMarshaller) MarshalRubyRegistry(ctx *Context) (err error) {
	return m.marshaller(ctx)
}

type Ruby struct {
	ID          string // ruby version including patch number
	TagLabel    string // user friendly ruby tag value
	Exe         string // ruby executable name
	Home        string // full path to ruby executable directory
	GemHome     string // full path to a ruby's gem home directory
	Description string // full ruby description
}

func init() {
	var err error
	rbRegex, err = regexp.Compile(`\A(j?ruby|rubinius)\s+(\d\.\d\.\d)(\w+)?(?:.+patchlevel )?(\d{1,3})?`)
	if err != nil {
		panic("unable to compile ruby parsing regexp")
	}

	rbVerRegex, err = regexp.Compile(`\A(\d\.\d\.\d)`)
	if err != nil {
		panic("unable to compile ruby version parsing regexp")
	}

	rbMajMinRegex, err = regexp.Compile(`\A(\d\.\d)`)
	if err != nil {
		panic("unable to compile ruby major/minor version parsing regexp")
	}

	SysRbRegex, err = regexp.Compile(`\Asys`)
	if err != nil {
		panic("unable to compile system ruby parsing regexp")
	}

	// list of known ruby executables
	KnownRubies = []string{`rbx`, `ruby`, `jruby`}
}

// CurrentRubyInfo returns the identifying tag and metadata information for the
// ruby currently in use.
func CurrentRubyInfo(ctx *Context) (tag string, info Ruby, err error) {
	envPath := os.Getenv(`PATH`)
	if envPath == `` {
		err = errors.New("Unable to read PATH environment variable")
		return
	}

	if strings.Index(envPath, Canary) != -1 {
		// prepended PATH looks like this, where GEM_HOME element is optional:
		//   GEM_HOME;RUBY_DIR;;;... -or- GEM_HOME:RUBY_DIR:::...
		head := strings.Split(envPath, string(os.PathListSeparator))[:3]
		var curRbPath string
		// test if the last element of the 2-element `head` slice is a blank
		// string which means GEM_HOME wasn't prepended to PATH
		if head[1] == `` {
			// scenario: RUBY_DIR;;;... -or- RUBY_DIR:::...
			curRbPath = head[0]
		} else {
			// scenario: GEM_HOME;RUBY_DIR;;;... -or- GEM_HOME:RUBY_DIR:::...
			curRbPath = head[1]
		}
		for _, v := range KnownRubies {
			tstRb := []string{curRbPath, v}
			tag, info, err = RubyInfo(ctx, strings.Join(tstRb, string(os.PathSeparator)))
			if err == nil {
				break
			}
		}
	} else {
		tags, err := TagLabelToTag(ctx, `system`)
		if err != nil {
			if len(ctx.Registry.Rubies) > 0 {
				// gracefully handle the scenario where a system ruby isn't included
				// in the registered rubies and PATH is the base PATH
				return ``, info, nil
			} else {
				return ``, info, errors.New("Unable to find tag for system ruby")
			}
		}
		for t, ri := range tags {
			if ri.TagLabel == `system` {
				tag = t
				break
			}
		}
		info = ctx.Registry.Rubies[tag]
	}

	return
}

// RubyInfo returns an identifying tag and metadata information about a specific
// ruby. It accepts a string of either the simple name of the ruby executable, or
// the ruby executables absolute path.
func RubyInfo(ctx *Context, ruby string) (tag string, info Ruby, err error) {
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
		if exe := res[1]; exe == `rubinius` {
			info.Exe = `rbx`
		} else {
			info.Exe = exe
		}
		if patch := res[3]; patch == `` {
			// patch up patchlevel for MRI 1.8.7's version string
			if patch187 := res[4]; patch187 != `` {
				info.ID = fmt.Sprintf("%s-p%s", res[2], patch187)
			} else {
				info.ID = res[2]
			}
		} else {
			info.ID = fmt.Sprintf("%s-%s", res[2], patch)
		}
		info.TagLabel = strings.Replace(strings.Replace(info.ID, `.`, ``, -1), `-`, ``, -1)
		tag, err = NewTag(ctx, info)
		if err != nil {
			// TODO implement
			panic("unable to create new tag for ruby")
		}
		info.GemHome = gemHome(info)
	} else {
		err = errors.New("unable to parse ruby name and version info")
		return
	}
	log.Printf("[DEBUG] tag: %s, %+v\n", tag, info)

	return
}

// marshalRubies persists the registered rubies to a JSON formatted file.
func marshalRubies(ctx *Context) (err error) {
	src := filepath.Join(ctx.Home(), `rubies.json`)
	dst := filepath.Join(ctx.Home(), `rubies.json.bak`)

	// TODO extract backup functionality to a utility function
	_, err = os.Stat(src)
	if err == nil {
		log.Printf("[DEBUG] backing up JSON ruby registry\n")
		_, e := CopyFile(dst, src)
		if e != nil {
			log.Println("[DEBUG] unable to backup JSON ruby registry")
			return e
		}
	}
	if os.IsNotExist(err) {
		log.Printf("[DEBUG] %s does not exist; creating\n", src)
		f, e := os.Create(src)
		if e != nil {
			log.Printf("[DEBUG] unable to create new %s\n", src)
			return e
		}
		defer f.Close()
	}

	b, err := json.Marshal(ctx.Registry)
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

// gemHome returns a string containing the filesystem location of a particular
// Ruby's gem home and is used to the the Ruby's GEM_HOME envar.
func gemHome(rb Ruby) string {
	usrHome := ``
	if runtime.GOOS == `windows` {
		usrHome = os.Getenv(`USERPROFILE`)
	} else {
		usrHome = os.Getenv(`HOME`)
	}

	rbLibVersion := rbVerRegex.FindStringSubmatch(rb.ID)[0]
	switch {
	case rbLibVersion >= `2.1.0`:
		rbLibVersion = fmt.Sprintf("%s.0", rbMajMinRegex.FindStringSubmatch(rbLibVersion)[0])
	}

	return filepath.Join(usrHome, `.gem`, rb.Exe, rbLibVersion)
}
