package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {

	var scriptFilePath = build()
	//fmt.Println(scriptFilePath)
	defer os.Remove(scriptFilePath)

	os.Chmod(scriptFilePath, 0744)
	var cmd = exec.Command(scriptFilePath)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func build() string {
	var scriptFileName string
	if bytes, err := os.ReadFile(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else if scriptFile, err := os.CreateTemp("", "disp_config*.sh"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		scriptFileName = scriptFile.Name()
		defer scriptFile.Close()

		var workspaces = Parse(bytes)
		var program = Generate(workspaces)
		for _, line := range program {
			//fmt.Println(line)
			fmt.Fprintln(scriptFile, line)
		}
	}
	return scriptFileName
}
