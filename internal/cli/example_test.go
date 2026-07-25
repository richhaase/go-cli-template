package cli

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

// executeCommand runs the root command with the given args and returns
// the combined output. This is the standard pattern for testing Cobra
// commands: point the command at a buffer and drive it via SetArgs.
func executeCommand(t *testing.T, ctx context.Context, args ...string) (string, error) {
	t.Helper()

	// Reset flag-bound package globals so table cases don't leak state.
	exampleName = "World"
	exampleCount = 1

	// Cobra only propagates the root context to a subcommand whose own
	// context is still unset, so when re-executing a shared command tree
	// across tests, set the context on the subcommand explicitly.
	exampleCmd.SetContext(ctx)

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)

	err := rootCmd.ExecuteContext(ctx)
	return buf.String(), err
}

func TestExampleCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name: "default greeting",
			args: []string{"example"},
			want: "Hello, World!\n",
		},
		{
			name: "positional message",
			args: []string{"example", "Howdy"},
			want: "Howdy, World!\n",
		},
		{
			name: "name and count flags",
			args: []string{"example", "--name", "Gopher", "--count", "2"},
			want: "Hello, Gopher!\nHello, Gopher!\n",
		},
		{
			name:    "too many args",
			args:    []string{"example", "one", "two"},
			wantErr: true,
		},
		{
			name:    "unknown flag",
			args:    []string{"example", "--bogus"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executeCommand(t, context.Background(), tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExampleCommandCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already canceled: the command should do no work

	out, err := executeCommand(t, ctx, "example")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
	if out != "" {
		t.Errorf("output = %q, want empty", out)
	}
}
