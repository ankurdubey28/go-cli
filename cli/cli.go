package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

func GetName(r io.Reader, w io.Writer) (string, error) {
	msg := "Your Name please? Press the Enter key when done. \n"
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

type Config struct {
	NumTimes   int
	PrintUsage bool
}

var UsageString = fmt.Sprintf(`usage: %s <Integer> [-h] --help 
A greeter application which prints the name you entered <Integer> times`, os.Args[0])

func PrintUsage(w io.Writer) {
	fmt.Fprint(w, UsageString)
}
func ParseArgs(args []string) (Config, error) {
	var numTimes int
	var err error
	c := Config{}
	if len(args) != 1 {
		return c, errors.New("invalid number of arguments")
	}
	if args[0] == "-h" || args[0] == "--help" {
		c.PrintUsage = true
		return c, nil
	}
	numTimes, err = strconv.Atoi(args[0])
	if err != nil {
		return c, err
	}
	c.NumTimes = numTimes
	return c, nil
}

func ValidateArgs(c Config) error {
	if c.PrintUsage {
		return nil
	}
	if !(c.NumTimes > 0) {
		return errors.New("must specify number greater than 0")
	}
	return nil
}

func RunCmd(r io.Reader, w io.Writer, c Config) error {
	if c.PrintUsage {
		PrintUsage(w)
		return nil
	}
	name, err := GetName(r, w)
	if err != nil {
		return err
	}
	greetUser(c, name, w)
	return nil
}

func greetUser(c Config, name string, w io.Writer) {
	msg := fmt.Sprintf("Nice to meet you %s\n", name)
	for i := 0; i < c.NumTimes; i++ {
		fmt.Fprint(w, msg)
	}
}
