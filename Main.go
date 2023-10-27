package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
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

	if scriptFilePath, err := build(configFilePath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		defer os.Remove(scriptFilePath)

		if *dumpFile == "" {
			if err := run(scriptFilePath); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

		} else if err := dump(scriptFilePath, *dumpFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}
	}

}

func run(scriptFilePath string) error {
	if err := os.Chmod(scriptFilePath, 0744); err != nil {
		return err
	} else if err := exec.Command(scriptFilePath).Run(); err != nil {
		return err
	} else {
		return nil
	}
}

func dump(scriptFilePath string, dumpFilePath string) error {
	if prog, err := os.ReadFile(scriptFilePath); err != nil {
		return err
	} else if dumpFilePath == "-" {
		fmt.Println(string(prog))
	} else if err := os.WriteFile(dumpFilePath, prog, 0744); err != nil {
		return err
	}

	return nil
}

func build(configFilePath string) (string, error) {
	var bytes []byte
	var workspaces []Workspace
	var scriptFileName string
	var scriptFile *os.File
	var err error

	if configFilePath == "" {
		if bytes, err = io.ReadAll(os.Stdin); err != nil {
			return "", err
		}
	} else if bytes, err = os.ReadFile(configFilePath); err != nil {
		return "", err
	}

	if scriptFile, err = os.CreateTemp("", "disp_config*.sh"); err != nil {
		return "", err
	} else {
		defer scriptFile.Close()
	}

	scriptFileName = scriptFile.Name()

	if workspaces, err = Parse(bytes); err != nil {
		return "", err
	} else {
		for _, line := range Generate(workspaces) {
			fmt.Fprintln(scriptFile, line)
		}
	}
	return scriptFileName, nil
}
