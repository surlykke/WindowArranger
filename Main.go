// Copyright (c) Christian Surlykke
//
// This file is part of the WindowArranger project.
// It is distributed under the GPL v2 license.
// Please refer to the GPL2 file for a copy of the license.
//
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	defer func() {
		if cause := recover(); cause != nil {
			fmt.Fprintln(os.Stderr, cause)
			os.Exit(1)
		}
	}()

	var dumpFile, waitSeconds, configFilePath = getCliArgs()	
	var bytes = getConfig(configFilePath)
	var scriptFile = getScriptFile()

	build(scriptFile, bytes, waitSeconds)
	if err := scriptFile.Close(); err != nil{
		panic(err)
	}
	defer os.Remove(scriptFile.Name())
	if dumpFile == "" {
		run(scriptFile.Name())
	} else {
		dump(scriptFile.Name(), dumpFile)
	}
}

func getCliArgs() (dumpFile string, wait uint, configFilePath string) {
	flag.Usage = func() {
		var out = flag.CommandLine.Output()
		fmt.Fprintln(out, "usage:")
		fmt.Fprintln(out, "  WindowWrapper [option]... [configfile]")
		fmt.Fprintln(out, "        If configfile is not given, configuration will be read from standard input")
		fmt.Fprintln(out, "options:")
		flag.PrintDefaults()
	}

	var df = flag.String("dump", "", "Dont execute generated script, but write it to a file. '-' means standard out.")
	var ws = flag.Uint("wait", 0, "Wait <uint seconds> for all criteria in config to match a window")
	flag.Parse()
	
	if len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	} 
	
	if len(flag.Args()) == 1 {
		return *df, *ws, flag.Args()[0]
	} else {
		return *df, *ws, ""
	}
}

func getConfig(configFilePath string) []byte {
	var (
		bytes []byte
		err error
	)

	if configFilePath == "" {
		bytes, err = io.ReadAll(os.Stdin)
	} else {
		bytes, err = os.ReadFile(configFilePath)
	}
	
	if err != nil {
		panic(err)
	}
	
	return bytes 
}

func getScriptFile() *os.File {
	if f, err := os.CreateTemp("", "disp_config*.sh"); err != nil {
		panic(err)
	} else {
		return f
	}
}

func build(scriptFile *os.File, bytes []byte, waitSeconds uint) {
	var workSpaces = Parse(bytes)
	for _, line := range Generate(workSpaces, waitSeconds) {
		if _, err := fmt.Fprintln(scriptFile, line); err != nil {
			panic(err)
		}
	}
}

func run(scriptFilePath string) {
	if err := os.Chmod(scriptFilePath, 0744); err != nil {
		panic(err)
	} else {
		var cmd = exec.Command(scriptFilePath)
		cmd.Stdout = os.Stdout 
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			os.Exit(1)
		}
	}
}

func dump(scriptFilePath string, dumpFilePath string) {
	if prog, err := os.ReadFile(scriptFilePath); err != nil {
		panic(err)
	} else if dumpFilePath == "-" {
		fmt.Println(string(prog))
	} else if err := os.WriteFile(dumpFilePath, prog, 0744); err != nil {
		panic(err)
	}
}
