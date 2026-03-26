package main

import (
	"bytes"
	"strings"
	"testing"
)

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


