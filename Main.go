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

	flag.Usage = func() {
		var out = flag.CommandLine.Output()
		fmt.Fprintln(out, "usage:")
		fmt.Fprintln(out, "  WindowWrapper [option]... [configfile]")
		fmt.Fprintln(out, "        if configfile is not given, configuration will be read from standard input")
		fmt.Fprintln(out, "options:")
		flag.PrintDefaults()
	}

	var dumpFile = flag.String("dump", "", "Dont execute but write script to a file")
	flag.Parse()

	var configFilePath string = ""
	if len(flag.Args()) > 1 {
		flag.Usage()
		os.Exit(1)
	} else if len(flag.Args()) == 1 {
		configFilePath = flag.Args()[0]
	}

	scriptFilePath := build(configFilePath)
	defer os.Remove(scriptFilePath)

	if *dumpFile == "" {
		run(scriptFilePath)
	} else {
		dump(scriptFilePath, *dumpFile)
	}

}

func run(scriptFilePath string) {
	if err := os.Chmod(scriptFilePath, 0744); err != nil {
		panic(err)
	} else if err := exec.Command(scriptFilePath).Run(); err != nil {
		panic(err)
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

func build(configFilePath string) string {
	var bytes []byte
	var scriptFileName string
	var scriptFile *os.File
	var err error

	if configFilePath == "" {
		bytes, err = io.ReadAll(os.Stdin)
	} else {
		bytes, err = os.ReadFile(configFilePath)
	}

	if err != nil {
		panic(err)
	}

	if scriptFile, err = os.CreateTemp("", "disp_config*.sh"); err != nil {
		panic(err)
	} else {
		defer scriptFile.Close()
	}

	scriptFileName = scriptFile.Name()

	for _, line := range Generate(Parse(bytes)) {
		fmt.Fprintln(scriptFile, line)
	}

	return scriptFileName
}
