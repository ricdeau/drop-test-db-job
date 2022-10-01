package internal

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	_ json.Unmarshaler = (*TimestampType)(nil)
)

type TimestampType int8

const (
	_ TimestampType = iota
	TimestampPrefix
	TimestampPostfix
)

func (t *TimestampType) UnmarshalJSON(bytes []byte) error {
	s := strings.Trim(string(bytes), `"`)

	var v TimestampType
	switch strings.ToLower(s) {
	case "prefix":
		v = TimestampPrefix
	case "postfix":
		v = TimestampPostfix
	default:
		return fmt.Errorf("unknown TimestampType: %q", s)
	}

	*t = v
	return nil
}
