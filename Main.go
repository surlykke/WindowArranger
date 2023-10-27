package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	flag.Usage = func() {
		var out = flag.CommandLine.Output()
		fmt.Fprintln(out, "usage:")
		fmt.Fprintln(out, "  WindowWrapper [option]... file")
		fmt.Fprintln(out, "options:")
		flag.PrintDefaults()
	}

	var dumpFile = flag.String("dump", "", "Dont execute but write script to a file")
	flag.Parse()

	var configFilePath string 
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	} else {
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
	var scriptFileName string
	var scriptFile *os.File
	var bytes []byte
	var err error

	if bytes, err = os.ReadFile(configFilePath); err != nil {
		return "", err
	}

	if scriptFile, err = os.CreateTemp("", "disp_config*.sh"); err != nil {
		return "", err
	} else {
		defer scriptFile.Close()
	}

	scriptFileName = scriptFile.Name()

	var program = Generate(Parse(bytes))
	for _, line := range program {
		fmt.Fprintln(scriptFile, line)
	}
	return scriptFileName, nil
}
