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
	"time"

	"github.com/surlykke/WindowArranger/compile"
	"github.com/surlykke/WindowArranger/sway"
)

var (
	input    io.ReadCloser
	out      io.WriteCloser
	wait     uint
	criteria []string
)

func main() {
	var (
		err            error
		dump           bool
		debug          bool
		configFilePath string
	)

	defer func() {
		if !debug {
			if r := recover(); r != nil {
				fmt.Fprintln(os.Stderr, "\nError:", r)
				os.Exit(1)
			}
		}
	}()

	dump, wait, debug, configFilePath = getCliArgs()
	fmt.Fprintln(os.Stderr, "dump: ", dump)
	fmt.Fprintln(os.Stderr, "waitSeconds: ", wait)
	if configFilePath != "" {
		fmt.Fprintln(os.Stderr, "input file: ", configFilePath)
	}

	if input, err = getInputReader(configFilePath); err != nil {
		panic(err)
	} else {
		defer input.Close()
		var program, criteria = compile.CompileConfig(input)
		if dump {
			fmt.Println()
			for _, c := range program {
				fmt.Println(c)
			}
		} else {
			doWait(criteria, wait)
			for _, cmd := range program {
				if debug {
					fmt.Fprintln(os.Stderr, cmd)
				}
				var responses = sway.Execute(cmd)
				for i, response := range responses {
					if !response.Success {
						panic(fmt.Sprintf("Command '%s' failed at subcommand %d: %s", cmd, i+1, response.Error))
					}
				}
			}
		}
	}
}

func getCliArgs() (dump bool, wait uint, debug bool, configFilePath string) {
	flag.Usage = func() {
		var out = flag.CommandLine.Output()
		fmt.Fprintln(out, "usage:")
		fmt.Fprintln(out, "  WindowWrapper [option]... [configfile]")
		fmt.Fprintln(out, "        If configfile is not given, configuration will be read from standard input")
		fmt.Fprintln(out, "options:")
		flag.PrintDefaults()
	}

	var df = flag.Bool("dump", false, "Dont execute generated commands, but write them to standard out.")
	var ws = flag.Uint("wait", 0, "Wait <uint seconds> for all criteria in config to match a window")
	var db = flag.Bool("debug", false, "Write commands to stdout. On error, give a bit info on where it occurred.")
	flag.Parse()

	if len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	}

	if len(flag.Args()) == 1 {
		return *df, *ws, *db, flag.Args()[0]
	} else {
		return *df, *ws, *db, ""
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

func doWait(criteria []string, secs uint) {
	var deadLine = time.Now().Add(time.Duration(secs) * time.Second)
	for {
		var allFound = true
		for _, criterium := range criteria {
			if !sway.Execute(fmt.Sprintf("[%s] focus", criterium))[0].Success {
				allFound = false
				break
			}
		}
		if allFound {
			return
		}
		if time.Now().After(deadLine) {
			panic("Not all criteria could be matched")
		}
		time.Sleep(1 * time.Second)
	}
}
