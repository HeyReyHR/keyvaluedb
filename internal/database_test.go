package internal

import (
	"io"
	"log/slog"
	"testing"

	"github.com/HeyReyHR/keyvaluedb/internal/compute/parser"
	"github.com/HeyReyHR/keyvaluedb/internal/storage"
)

type mockParser struct {
	parseFunc func(string) (parser.Query, error)
}

func (m *mockParser) Parse(input string) (parser.Query, error) {
	return m.parseFunc(input)
}

type mockStorage struct {
	setFunc func(key, value string) error
	getFunc func(key string) (string, error)
	delFunc func(key string) error
}

func (m *mockStorage) Set(key, value string) error    { return m.setFunc(key, value) }
func (m *mockStorage) Get(key string) (string, error) { return m.getFunc(key) }
func (m *mockStorage) Del(key string) error           { return m.delFunc(key) }

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestNewDatabase(t *testing.T) {
	validParser := &mockParser{}
	validStorage := &mockStorage{}

	tests := []struct {
		name         string
		computeLayer computeLayer
		storageLayer storageLayer
		wantErr      bool
	}{
		{"valid", validParser, validStorage, false},
		{"nil compute layer", nil, validStorage, true},
		{"nil storage layer", validParser, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDatabase(tt.computeLayer, tt.storageLayer, newTestLogger())

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if db != nil {
					t.Error("expected nil database on error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if db == nil {
				t.Fatal("expected non-nil database")
			}
		})
	}
}

func TestDatabase_HandleQuery(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		parseFunc  func(string) (parser.Query, error)
		setFunc    func(key, value string) error
		getFunc    func(key string) (string, error)
		delFunc    func(key string) error
		wantOutput string
	}{
		{
			name:  "parse error",
			input: "KEK",
			parseFunc: func(string) (parser.Query, error) {
				return parser.Query{}, parser.ErrInvalidCommand
			},
			wantOutput: "[error] invalid command",
		},
		{
			name:  "GET found",
			input: "GET key1",
			parseFunc: func(string) (parser.Query, error) {
				return parser.NewQuery(parser.GetCommandId, []string{"key1"}), nil
			},
			getFunc: func(key string) (string, error) {
				return "value1", nil
			},
			wantOutput: "[ok] value1",
		},
		{
			name:  "GET not found",
			input: "GET missing",
			parseFunc: func(string) (parser.Query, error) {
				return parser.NewQuery(parser.GetCommandId, []string{"missing"}), nil
			},
			getFunc: func(key string) (string, error) {
				return "", storage.ErrNotFound
			},
			wantOutput: "[not found]",
		},
		{
			name:  "SET ok",
			input: "SET key1 value1",
			parseFunc: func(string) (parser.Query, error) {
				return parser.NewQuery(parser.SetCommandId, []string{"key1", "value1"}), nil
			},
			setFunc: func(key, value string) error {
				return nil
			},
			wantOutput: "[ok]",
		},
		{
			name:  "DEL ok",
			input: "DEL key1",
			parseFunc: func(string) (parser.Query, error) {
				return parser.NewQuery(parser.DelCommandId, []string{"key1"}), nil
			},
			delFunc: func(key string) error {
				return nil
			},
			wantOutput: "[ok]",
		},
		{
			name:  "unknown command id",
			input: "???",
			parseFunc: func(string) (parser.Query, error) {
				return parser.NewQuery(999, nil), nil
			},
			wantOutput: "[error] internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &mockParser{parseFunc: tt.parseFunc}
			ms := &mockStorage{
				setFunc: tt.setFunc,
				getFunc: tt.getFunc,
				delFunc: tt.delFunc,
			}

			db, err := NewDatabase(mp, ms, newTestLogger())
			if err != nil {
				t.Fatalf("unexpected error creating database: %v", err)
			}

			got := db.HandleQuery(tt.input)
			if got != tt.wantOutput {
				t.Errorf("HandleQuery(%q) = %q, want %q", tt.input, got, tt.wantOutput)
			}
		})
	}
}
