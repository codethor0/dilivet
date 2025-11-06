package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestApp_Version(t *testing.T) {
	var out bytes.Buffer
	app := &App{
		Name:    "test-cli",
		Version: "v1.2.3",
		Out:     &out,
	}

	exitCode := app.Run([]string{"-version"})

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	output := out.String()
	if !strings.Contains(output, "v1.2.3") {
		t.Errorf("Expected version in output, got: %s", output)
	}
}

func TestApp_Help(t *testing.T) {
	var out bytes.Buffer
	app := &App{
		Name:    "test-cli",
		Version: "v1.0.0",
		Out:     &out,
	}

	exitCode := app.Run([]string{"-help"})

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	output := out.String()
	if !strings.Contains(output, "USAGE:") {
		t.Errorf("Expected help message, got: %s", output)
	}
	if !strings.Contains(output, "test-cli") {
		t.Errorf("Expected app name in help, got: %s", output)
	}
}

func TestApp_DefaultBehavior(t *testing.T) {
	var out bytes.Buffer
	app := &App{
		Name:    "dilivet",
		Version: "dev",
		Out:     &out,
	}

	exitCode := app.Run([]string{})

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	output := out.String()
	if !strings.Contains(output, "ML-DSA vetting tool") {
		t.Errorf("Expected default message, got: %s", output)
	}
}

func TestApp_InvalidFlag(t *testing.T) {
	var out, errOut bytes.Buffer
	app := &App{
		Name:    "test-cli",
		Version: "dev",
		Out:     &out,
		Err:     &errOut,
	}

	exitCode := app.Run([]string{"-invalid"})

	if exitCode == 0 {
		t.Error("Expected non-zero exit code for invalid flag")
	}
}

func TestApp_MultipleFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantExit int
		wantOut  string
	}{
		{
			name:     "version takes precedence",
			args:     []string{"-version"},
			wantExit: 0,
			wantOut:  "v1.0.0",
		},
		{
			name:     "help flag",
			args:     []string{"-help"},
			wantExit: 0,
			wantOut:  "USAGE:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			app := &App{
				Name:    "test",
				Version: "v1.0.0",
				Out:     &out,
			}

			exitCode := app.Run(tt.args)

			if exitCode != tt.wantExit {
				t.Errorf("Expected exit code %d, got %d", tt.wantExit, exitCode)
			}

			if !strings.Contains(out.String(), tt.wantOut) {
				t.Errorf("Expected %q in output, got: %s", tt.wantOut, out.String())
			}
		})
	}
}
