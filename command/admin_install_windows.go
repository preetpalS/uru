// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["install"] = Command{
		Name:    "install",
		Aliases: []string{"install", "in"},
		Usage:   "install",
		HelpMsg: "install uru",
		Eg:      `install`}
}

func adminInstall(ctx *env.Context) {
	_, err := exec.LookPath("uru_rt.exe")
	if err != nil {
		fmt.Printf("[ERROR] uru_rt.exe must be present in a directory on PATH\n")
		os.Exit(1)
	}
	_, err = os.Stat("uru_rt.exe")
	if os.IsNotExist(err) {
		fmt.Printf("[ERROR] must install from same directory as uru_rt.exe\n")
		os.Exit(1)
	}

	for _, v := range []string{"uru.bat", "uru.ps1"} {
		_, err := os.Stat(v)
		if err == nil {
			log.Printf("[DEBUG] creating backup of `%s`\n", v)
			_, e := env.CopyFile(fmt.Sprintf("%s.bak", v), v)
			if e != nil {
				log.Printf("[DEBUG] failed to backup `%s`; continuing", v)
			}
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ``
	} else {
		cwd = fmt.Sprintf("into %s", cwd)
	}
	fmt.Printf("---> Installing uru %s\n", cwd)

	for k, v := range map[string]string{"uru.bat": env.BatWrapper, "uru.ps1": env.PSWrapper} {
		script, err := os.Create(k)
		if err != nil {
			panic(fmt.Sprintf("unable to create `%s` script wrapper", k))
		}
		defer script.Close()

		_, err = script.WriteString(v)
		if err != nil {
			panic(fmt.Sprintf("failed to write `%s` script wrapper", k))
		}
	}
}
