package main

import (
	"fmt"
	"os"

	"ankurdubey28/github.com/go-cli/cli"
)

func main() {

	c, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		cli.PrintUsage(os.Stdout)
		os.Exit(1)
	}
	err = cli.ValidateArgs(c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		cli.PrintUsage(os.Stdout)
		os.Exit(1)
	}
	err = cli.RunCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
