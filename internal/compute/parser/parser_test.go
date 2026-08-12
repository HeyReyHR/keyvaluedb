package parser

import (
	"errors"
	"io"
	"log/slog"
	"testing"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
func TestParser_Parse(t *testing.T) {
	p := NewParser(newTestLogger())
	tests := []struct {
		name      string
		input     string
		wantCmdId int
		wantArgs  []string
		wantErr   error
	}{
		{
			name:      "valid GET",
			input:     "GET key1",
			wantCmdId: GetCommandId,
			wantArgs:  []string{"key1"},
			wantErr:   nil,
		},
		{
			name:      "valid SET",
			input:     "SET key1 value1",
			wantCmdId: SetCommandId,
			wantArgs:  []string{"key1", "value1"},
			wantErr:   nil,
		},
		{
			name:      "valid DEL",
			input:     "DEL key1",
			wantCmdId: DelCommandId,
			wantArgs:  []string{"key1"},
			wantErr:   nil,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: ErrInvalidQuery,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: ErrInvalidQuery,
		},
		{
			name:    "unknown command",
			input:   "LOL key1",
			wantErr: ErrInvalidCommand,
		},
		{
			name:    "SET missing value",
			input:   "SET key1",
			wantErr: ErrInvalidArgsNumber,
		},
		{
			name:    "SET too many args",
			input:   "SET key1 value1 extra",
			wantErr: ErrInvalidArgsNumber,
		},
		{
			name:    "GET no args",
			input:   "GET",
			wantErr: ErrInvalidArgsNumber,
		},
		{
			name:    "GET too many args",
			input:   "GET key1 key2",
			wantErr: ErrInvalidArgsNumber,
		},
		{
			name:      "extra whitespace between tokens",
			input:     "GET    key1",
			wantCmdId: GetCommandId,
			wantArgs:  []string{"key1"},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := p.Parse(tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Parse(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}

			if query.CommandId() != tt.wantCmdId {
				t.Errorf("Parse(%q) CommandId = %v, want %v", tt.input, query.CommandId(), tt.wantCmdId)
			}

			gotArgs := query.Arguments()
			if len(gotArgs) != len(tt.wantArgs) {
				t.Fatalf("Parse(%q) Arguments = %v, want %v", tt.input, gotArgs, tt.wantArgs)
			}
			for i := range gotArgs {
				if gotArgs[i] != tt.wantArgs[i] {
					t.Errorf("Parse(%q) Arguments[%d] = %q, want %q", tt.input, i, gotArgs[i], tt.wantArgs[i])
				}
			}
		})
	}
}
