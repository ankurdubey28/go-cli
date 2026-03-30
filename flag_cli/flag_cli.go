package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Config struct {
	numTimes int
	name     string
}

var errPosArgSpecified = errors.New("more than one positional argument specified")

func main() {
	c, err := ParseArgs(os.Stdout, os.Args[1:])
	if err != nil && errors.Is(err, errPosArgSpecified) {
		fmt.Fprint(os.Stdout, err.Error())
		os.Exit(1)
	}
	err = ValidateArgs(c)
	if err != nil {
		fmt.Fprint(os.Stdout, err.Error())
		os.Exit(1)
	}
	err = RunCmd(os.Stdin, os.Stdout, c)
	if err != nil {
		fmt.Fprint(os.Stdout, err.Error())
		os.Exit(1)
	}
}

func ParseArgs(w io.Writer, args []string) (Config, error) {
	c := Config{}
	fs := flag.NewFlagSet("Greeter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Usage = func() {
		var usageString = `
A greeter application which prints the name you entered a specified 
number of times.
Usage of %s: <options> [name]`
		fmt.Fprintf(w, usageString, fs.Name())
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Options: ")
		fs.PrintDefaults()
	}
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	if fs.NArg() > 1 {
		return c, errPosArgSpecified
	}
	if fs.NArg() == 1 {
		c.name = fs.Arg(0)
	}
	return c, nil
}

func GetName(w io.Writer, r io.Reader) (string, error) {
	msg := "Tell your name and then press Enter\n"
	fmt.Fprint(w, msg)
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}
	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("you did not enter any name")
	}
	return name, nil

}

func greetUser(c Config, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", c.name)
	for i := 0; i < c.numTimes; i++ {
		fmt.Fprint(w, msg)
	}
}

func RunCmd(r io.Reader, w io.Writer, c Config) error {
	var err error
	if len(c.name) == 0 {
		c.name, err = GetName(w, r)
		if err != nil {
			return err
		}
	}
	greetUser(c, w)
	return nil
}

func ValidateArgs(c Config) error {
	if !(c.numTimes > 0) {
		return errors.New("must specify number greater than 0")
	}
	return nil
}
