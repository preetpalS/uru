// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package exec

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"bitbucket.org/jonforums/uru/env"
)

var scriptName string

// runner script templates
var batScript = `@ECHO OFF
REM autogenerated by uru

SET PATH=%s
SET GEM_HOME=%s
`

var ps1Script = `# autogenerated by uru

$env:PATH = "%s"
$env:GEM_HOME = "%s"
`

var bashScript = `# autogenerated by uru

NEW_PATH=%s
`

// CreateScript creates a runner script tuned to the type of shell that called
// the uru runtime.
func CreateScript(ctx *env.Context, path *[]string, gemHome string) {
	scriptType := os.Getenv(`URU_INVOKER`)

	var script string
	switch scriptType {
	case `powershell`:
		script = ps1Script
		scriptName = "uru_lackee.ps1"
	case `batch`:
		script = batScript
		scriptName = "uru_lackee.bat"
	case `bash`:
		script = bashScript
		scriptName = "uru_lackee"
	default:
		panic("uru invoked from unknown shell (check URU_INVOKER env var)")
	}
	log.Printf("[DEBUG] runner script: %s\n", scriptName)

	runner := filepath.Join(ctx.Home(), scriptName)
	f, err := os.Create(runner)
	if err != nil {
		panic(fmt.Sprintf("unable to create `%s` runner script", runner))
	}
	defer f.Close()

	if runtime.GOOS != `windows` {
		f.Chmod(0755)
	}

	// TODO refactor to use `text\template` awesomeness?
	// hackishly modify the bash script template to suppress GEM_HOME creation
	// when nonexistent
	var content string
	if scriptType == `bash`  {
		if gemHome != `` {
			script = strings.Join([]string{script, "NEW_GEM_HOME=%s\n"}, ``)
			content = fmt.Sprintf(script, strings.Join(*path, string(os.PathListSeparator)), gemHome)
		} else {
			script = strings.Join([]string{script, "UNSET_GEM_HOME=yes\n"}, ``)
			content = fmt.Sprintf(script, strings.Join(*path, string(os.PathListSeparator)))
		}
	}

	_, err = f.WriteString(content)
	if err != nil {
		panic(fmt.Sprintf("failed to write `%s` runner script", runner))
	}
}

// ExecScript executes the runner script based upon the type of shell that
// called the uru runtime.
func ExecScript(ctx *env.Context) {

	runner := filepath.Join(ctx.Home(), scriptName)
	var cmd *exec.Cmd
	switch os.Getenv(`URU_INVOKER`) {
	case `powershell`:
		cmd = exec.Command("powershell", "-file", runner)
	case `batch`:
		cmd = exec.Command("cmd", runner)
	case `bash`:
		cmd = exec.Command("sh", runner)
	default:
		panic("uru invoked from unknown shell (check URU_INVOKER env var)")
	}

	cmd.Run()
}
