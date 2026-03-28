package flag_cli

import (
	"errors"
	"flag"
	"io"
)

type config struct {
	numTimes int
}

func ParseArgs(w io.Writer, r io.Reader, args []string) (config, error) {
	c := config{}
	fs := flag.NewFlagSet("Greeter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.IntVar(&c.numTimes, "n", 0, "Number of times to greet")

	err := fs.Parse(args)
	if err != nil {
		return c, err
	}
	if fs.NArg() != 0 {
		return c, errors.New("positional Argument specified")
	}
	return c, nil
}
