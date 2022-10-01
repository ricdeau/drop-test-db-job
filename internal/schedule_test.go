package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchedule_UnmarshalJSON(t *testing.T) {
	type testData struct {
		Schedule Schedule `json:"schedule"`
	}
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "standard",
			input:   `{"schedule":"* * * * ?"}`,
			wantErr: false,
		},
		{
			name:    "every",
			input:   `{"schedule":"@every 1h10m"}`,
			wantErr: false,
		},
		{
			name:    "special",
			input:   `{"schedule":"@hourly"}`,
			wantErr: false,
		},
		{
			name:    "error",
			input:   `{"schedule":"invalid"}`,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got testData
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				t.Log(err)
			} else {
				require.NotEmpty(t, got.Schedule)
			}
		})
	}
}
