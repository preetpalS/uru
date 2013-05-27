// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Initialize uru's home directory, creating if necessary.
func initHome() {
	uruHome := os.Getenv(`URU_HOME`)
	if runtime.GOOS == `windows` {
		if uruHome == `` {
			ctx.SetHome(filepath.Join(os.Getenv(`USERPROFILE`), `.uru`))
		} else {
			ctx.SetHome(uruHome)
		}
	} else {
		if uruHome == `` {
			ctx.SetHome(filepath.Join(os.Getenv(`HOME`), `.uru`))
		} else {
			ctx.SetHome(uruHome)
		}
	}
	log.Printf("[DEBUG] uru HOME is %s\n", ctx.Home())

	_, err := os.Stat(ctx.Home())
	if os.IsNotExist(err) {
		log.Printf("[DEBUG] creating %s\n", ctx.Home())
		os.Mkdir(ctx.Home(), os.ModeDir|0750)
	}

	// purge existing runners to prevent bogus environment changes
	walk := func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(filepath.Base(path), `uru_lackee`) {
			log.Printf("[DEBUG] deleting runner script %s\n", path)
			_ = os.Remove(path) // TODO throw away the error?
		}
		return nil
	}
	filepath.Walk(ctx.Home(), walk)
}

// Make uru's Context ready for general use.
func initContext() {
	ctx.Init()
}

// Initialize uru's CLI parser.
func initCommandParser() {
	ctx.SetCmdRegex(`admin`, `\Aadmin\z`)
	ctx.SetCmdRegex(`add`, `\Aadd\z`)
	ctx.SetCmdRegex(`gem`, `\Agem\z`)
	ctx.SetCmdRegex(`help`, `\Ahelp\z`)
	ctx.SetCmdRegex(`install`, `\Ainstall|in\z`)
	ctx.SetCmdRegex(`ls`, `\Als|list\z`)
	ctx.SetCmdRegex(`refresh`, `\Arefresh\z`)
	ctx.SetCmdRegex(`rm`, `\Arm|del\z`)
	ctx.SetCmdRegex(`ruby`, `\Aruby|rb\z`)
	ctx.SetCmdRegex(`version`, `\Aver(?:sion)?\z`)
}

// Import all installed rubies that have been registered with uru.
func initRubies() {
	rubies := filepath.Join(ctx.Home(), `rubies.json`)
	_, err := os.Stat(rubies)
	if os.IsNotExist(err) {
		log.Printf("[DEBUG] %s does not exist\n", rubies)
		return
	}

	b, err := ioutil.ReadFile(rubies)
	if err != nil {
		log.Printf("[DEBUG] unable to read %s\n", rubies)
		panic("unable to read the JSON ruby registry")
	}

	err = json.Unmarshal(b, &ctx.Registry)
	if err != nil {
		log.Printf("[DEBUG] unable to unmarshal %s\n", rubies)
		panic("unable to unmarshal the JSON ruby registry")
	}
	log.Printf("[DEBUG] === ctx.Registry.Rubies ===\n%+v", ctx.Registry.Rubies)
}
