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
	Config
	err  error
	args []string
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Config
		wantErr string
	}{
		{
			name: "help flag",
			args: []string{"-h"},
			want: Config{PrintUsage: true},
		},
		{
			name: "valid number",
			args: []string{"3"},
			want: Config{NumTimes: 3},
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
			got, err := ParseArgs(tt.args)

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
		cfg     Config
		wantErr string
	}{
		{
			name: "valid",
			cfg:  Config{NumTimes: 2},
		},
		{
			name:    "zero",
			cfg:     Config{NumTimes: 0},
			wantErr: "must specify number greater than 0",
		},
		{
			name:    "negative",
			cfg:     Config{NumTimes: -1},
			wantErr: "must specify number greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateArgs(tt.cfg)

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

			got, err := GetName(in, out)

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
		cfg      Config
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "print usage",
			cfg:      Config{PrintUsage: true},
			expected: UsageString,
		},
		{
			name:     "normal flow",
			cfg:      Config{NumTimes: 2},
			input:    "Ankur\n",
			expected: "Your Name please? Press the Enter key when done. \nNice to meet you Ankur\nNice to meet you Ankur\n",
		},
		{
			name:    "empty name",
			cfg:     Config{NumTimes: 1},
			input:   "\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.input)
			out := &bytes.Buffer{}

			err := RunCmd(in, out, tt.cfg)

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

func TestMain(m *testing.M) {
	cmd := exec.Command("go", "build", "-o", binaryPath)
	if err := cmd.Run(); err != nil {
		fmt.Print("failed to build binary:", err)
		os.Exit(1)
	}

	code := m.Run()

	os.Remove(binaryPath)
	os.Exit(code)
}

func runCLI(args []string, input string) (string, int, error) {
	cmd := exec.Command(binaryPath, args...)

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
