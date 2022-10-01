package internal

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTimestampTypeUnmarshalJSON(t *testing.T) {
	type testData struct {
		TsType TimestampType `json:"tsType"`
	}
	tests := []struct {
		name    string
		input   string
		want    testData
		wantErr bool
	}{
		{
			name:  "success",
			input: `{"tsType":"prefix"}`,
			want: testData{
				TsType: TimestampPrefix,
			},
			wantErr: false,
		},
		{
			name:    "error",
			input:   `{"tsType":"invalid"}`,
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
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("UnmarshalJSON() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
