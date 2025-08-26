package repository

import (
	"database/sql"
	"testing"
)

func TestToNullString(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  sql.NullString
	}{
		{
			name:  "null input",
			input: []byte("null"),
			want:  sql.NullString{String: "", Valid: false},
		},
		{
			name:  "valid json",
			input: []byte(`{"key":"value"}`),
			want:  sql.NullString{String: `{"key":"value"}`, Valid: true},
		},
		{
			name:  "empty string",
			input: []byte(""),
			want:  sql.NullString{String: "", Valid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toNullString(tt.input); got.String != tt.want.String || got.Valid != tt.want.Valid {
				t.Errorf("toNullString() = %v, want %v", got, tt.want)
			}
		})
	}
}
