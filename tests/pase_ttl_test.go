package tests

import (
	"drop-test-db-job/pkgs/utils"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testData struct {
	input    string
	expected time.Duration
	error    error
}

var testSource = []testData{
	{input: "1h 15m", expected: 1*time.Hour + 15*time.Minute, error: nil},
	{input: " 1h  15m ", expected: 1*time.Hour + 15*time.Minute, error: nil},
	{input: "1h15m", expected: 1*time.Hour + 15*time.Minute, error: nil},
	{input: "1h", expected: 1 * time.Hour, error: nil},
	{input: "45m", expected: 45 * time.Minute, error: nil},
	{input: "155m", expected: 155 * time.Minute, error: nil},
	{input: "1h 15m 30s", expected: 0, error: fmt.Errorf(`ttl string: invalid format "1h 15m 30s", ttl examples: "1h", "2h 45m", "30m"`)},
	{input: "1h 15m 50m", expected: 0, error: fmt.Errorf(`ttl string: invalid format "1h 15m 50m", ttl examples: "1h", "2h 45m", "30m"`)},
	{input: "1h 1h", expected: 0, error: fmt.Errorf(`ttl string: invalid format "1h 1h", ttl examples: "1h", "2h 45m", "30m"`)},
	{input: "30m 2h", expected: 0, error: fmt.Errorf(`ttl string: invalid format "30m 2h", ttl examples: "1h", "2h 45m", "30m"`)},
	{input: "100", expected: 0, error: fmt.Errorf(`ttl string: invalid format "100", ttl examples: "1h", "2h 45m", "30m"`)},
	{input: "absdef", expected: 0, error: fmt.Errorf(`ttl string: invalid format "absdef", ttl examples: "1h", "2h 45m", "30m"`)},
	{input: "", expected: 0, error: fmt.Errorf(`ttl string: invalid format "", ttl examples: "1h", "2h 45m", "30m"`)},
}

func TestParallelTryParseTTL(t *testing.T) {
	for _, testCase := range testSource {
		tc := testCase
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			testCaseTryParseTTL(t, &tc)
		})
	}
}

func testCaseTryParseTTL(t *testing.T, testItem *testData) {
	duration, e := utils.TryParseTTL(testItem.input)
	assert.Equal(t, testItem.error, e)
	assert.Equal(t, testItem.expected, duration)
}
