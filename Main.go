// Copyright (c) Christian Surlykke
//
// This file is part of the WindowArranger project.
// It is distributed under the GPL v2 license.
// Please refer to the GPL2 file for a copy of the license.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Workspace struct {
	Output   string
	Layout   string
	Children []*Node
}

type Node struct {
	Criteria string
	Layout   string
	Children []*Node
}

func main() {
	var dumpFile, waitSeconds, configFilePath = getCliArgs()

	inputReader, err := getInputReader(configFilePath)
	if err != nil {
		panic(err)
	}

	outputWriter, err := getOutputWriter(dumpFile)
	if err != nil {
		panic(err)
	}

	translate(inputReader, outputWriter, waitSeconds)

	inputReader.Close()
	outputWriter.Close()

	if dumpFile == "" {
		run(outputWriter.(*os.File).Name())
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

func getInputReader(path string) (io.ReadCloser, error) {
	if path == "" {
		return os.Stdin, nil
	} else if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		return f, nil
	}
}

func getOutputWriter(dumpFile string) (io.WriteCloser, error) {
	switch dumpFile {
	case "":
		return os.CreateTemp("", "disp_config*.sh")
	case "-":
		return os.Stdout, nil
	default:
		return os.Create(dumpFile)
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
