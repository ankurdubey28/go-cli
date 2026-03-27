package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var binaryPath = "./application.exe"

type testConfig struct {
	config
	err  error
	args []string
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    config
		wantErr string
	}{
		{
			name: "help flag",
			args: []string{"-h"},
			want: config{printUsage: true},
		},
		{
			name: "valid number",
			args: []string{"3"},
			want: config{numTimes: 3},
		},
		{
			name:    "no args",
			args:    []string{},
			wantErr: "invalid number of arguments",
		},
		{
			name:    "invalid number",
			args:    []string{"abc"},
			wantErr: `strconv.Atoi: parsing "abc": invalid syntax`,
		},
		{
			name:    "too many args",
			args:    []string{"1", "foo"},
			wantErr: "invalid number of arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config
		wantErr string
	}{
		{
			name: "valid",
			cfg:  config{numTimes: 2},
		},
		{
			name:    "zero",
			cfg:     config{numTimes: 0},
			wantErr: "must specify number greater than 0",
		},
		{
			name:    "negative",
			cfg:     config{numTimes: -1},
			wantErr: "must specify number greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArgs(tt.cfg)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr string
	}{
		{
			name:  "valid input",
			input: "Ankur\n",
			want:  "Ankur",
		},
		{
			name:    "empty input",
			input:   "\n",
			wantErr: "you did not enter any name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}

			got, err := getName(in, out)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "print usage",
			cfg:      config{printUsage: true},
			expected: usageString,
		},
		{
			name:     "normal flow",
			cfg:      config{numTimes: 2},
			input:    "Ankur\n",
			expected: "Your Name please? Press the Enter key when done. \nNice to meet you Ankur\nNice to meet you Ankur\n",
		},
		{
			name:    "empty name",
			cfg:     config{numTimes: 1},
			input:   "\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}

			err := runCmd(in, out, tt.cfg)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if out.String() != tt.expected {
				t.Fatalf("expected:\n%q\ngot:\n%q", tt.expected, out.String())
			}
		})
	}
}

// TestMain
// 1. This is not a separate test case. It is a test runner controller.
// 2. This is a custom entry point for the test binary itself.

func TestMain(m *testing.M) {
	// build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath)
	if err := cmd.Run(); err != nil {
		fmt.Print("failed to build binary:", err)
		os.Exit(1)
	}

	// run tests
	code := m.Run()

	// cleanup
	os.Remove(binaryPath)
	os.Exit(code)
}

func runCLI(args []string, input string) (string, int, error) {
	cmd := exec.Command(binaryPath, args...)

	// simulate stdin
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			exitCode = e.ExitCode()
		} else {
			return "", -1, err
		}
	}
	return out.String(), exitCode, nil
}

func TestMain_Success(t *testing.T) {
	out, code, err := runCLI([]string{"2"}, "Ankur\n")
	if err != nil {
		t.Fatal(err)
	}

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	expected := "Nice to meet you Ankur\nNice to meet you Ankur\n"
	if !strings.Contains(out, expected) {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestMain_NoArgs(t *testing.T) {
	out, code, _ := runCLI([]string{}, "")

	if code == 0 {
		t.Fatalf("expected non-zero exit code")
	}

	if !strings.Contains(out, "invalid number of arguments") {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestMain_InvalidNumber(t *testing.T) {
	out, code, _ := runCLI([]string{"abc"}, "")

	if code == 0 {
		t.Fatalf("expected non-zero exit code")
	}

	if !strings.Contains(out, "invalid syntax") {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestMain_Zero(t *testing.T) {
	out, code, _ := runCLI([]string{"0"}, "")

	if code == 0 {
		t.Fatalf("expected non-zero exit code")
	}

	if !strings.Contains(out, "must specify number greater than 0") {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestMain_Help(t *testing.T) {
	out, code, _ := runCLI([]string{"-h"}, "")

	if code != 0 {
		t.Fatalf("expected exit code 0")
	}

	if !strings.Contains(out, "usage:") {
		t.Fatalf("expected usage output")
	}
}
