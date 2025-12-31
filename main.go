// Copyright (c) Christian Surlykke
//
// This file is part of the WindowArranger project.
// It is distributed under the GPL v2 license.
// Please refer to the GPL2 file for a copy of the license.
package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	layout   Layout
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
		var program = CompileConfig()
		if dump {
			fmt.Println()
			for _, c := range program {
				fmt.Println(c)
			}
		} else {
			var sway = connectToSway()
			doWait(criteria, wait, sway)
			for _, cmd := range program {
				executeCmd(cmd, sway)
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
	var db = flag.Bool("debug", false, "On error, give a bit more info on where it occurred.")
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

func connectToSway() net.Conn {
	if swaysock := os.Getenv("SWAYSOCK"); swaysock == "" {
		panic("SWAYSOCK not set")
	} else if sock, err := net.Dial("unix", swaysock); err != nil {
		panic(err)
	} else {
		return sock
	}
}

func doWait(criteria []string, secs uint, sway net.Conn) {
	var deadLine = time.Now().Add(time.Duration(secs) * time.Second)
	for {
		var allFound = true
		for _, criterium := range criteria {
			if !sendCommand(fmt.Sprintf("[%s] focus", criterium), sway)[0]["success"].(bool) {
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

func executeCmd(cmd string, sway net.Conn) {
	if resp := sendCommand(cmd, sway); len(resp) == 0 {
		panic("Empty response")
	} else {
		for _, subresp := range resp {
			if success, ok := subresp["success"].(bool); !(ok && success) {
				panic("Command failed:" + cmd + ":" + subresp["error"].(string))
			}
		}
	}
}

func sendCommand(cmd string, sway net.Conn) []map[string]any {
	var (
		msg  = []byte("i3-ipc")
		resp = make([]byte, 2000)
		m    = make([]map[string]any, 10)
	)
	msg = binary.NativeEndian.AppendUint32(msg, uint32(len([]byte(cmd))))
	msg = binary.NativeEndian.AppendUint32(msg, 0)
	msg = append(msg, []byte(cmd)...)
	if _, err := sway.Write(msg); err != nil {
		panic(err)
	} else if n, err := sway.Read(resp); err != nil {
		panic(err)
	} else if err := json.Unmarshal(resp[14:n], &m); err != nil {
		panic(err)
	}

	return m
}
