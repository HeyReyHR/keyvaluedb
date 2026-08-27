package common

import "testing"

func TestParseSize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "plain bytes", input: "100", want: 100},
		{name: "plain bytes with whitespace", input: "  100  ", want: 100},
		{name: "kilobytes", input: "10KB", want: 10 * 1024},
		{name: "megabytes", input: "2MB", want: 2 * 1024 * 1024},
		{name: "gigabytes", input: "1GB", want: 1024 * 1024 * 1024},
		{name: "bytes suffix", input: "5B", want: 5},
		{name: "zero bytes", input: "0B", want: 0},
		{name: "lowercase suffix", input: "1gb", want: 1024 * 1024 * 1024},
		{name: "space between number and suffix", input: "1 GB", want: 1024 * 1024 * 1024},
		{name: "negative plain number is not validated", input: "-5", wantErr: true},
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
		{name: "garbage without suffix", input: "abc", wantErr: true},
		{name: "suffix without a number", input: "GB", wantErr: true},
		{name: "non-numeric value with suffix", input: "XGB", wantErr: true},
		{name: "float value is not supported", input: "1.5GB", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSize(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseSize(%q) error = nil, want error", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("ParseSize(%q) unexpected error = %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseSize(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}
