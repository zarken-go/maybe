package maybe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type MaybeBool struct { //nolint:golint
	ValidFlag
	Bool bool
}

func Bool(v bool) MaybeBool {
	return MaybeBool{
		ValidFlag: true,
		Bool:      v,
	}
}

func (m MaybeBool) MarshalJSON() ([]byte, error) {
	if !m.ValidFlag {
		return nullBytes, nil
	}
	return json.Marshal(m.Bool)
}

func (m *MaybeBool) Set(v bool) {
	m.ValidFlag = true
	m.Bool = v
}

func (m *MaybeBool) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		m.Bool = false
		m.ValidFlag = false
		return nil
	}

	var v bool
	if err := json.Unmarshal(b, &v); err != nil {
		var i int64
		if json.Unmarshal(b, &i) == nil {
			return m.SetE(i)
		}

		var s string
		if json.Unmarshal(b, &s) == nil {
			return m.SetE(s)
		}

		return err
	}
	m.Set(v)
	return nil
}

func (m *MaybeBool) SetE(v interface{}) error {
	if v == nil {
		m.ValidFlag = false
		m.Bool = false
		return nil
	}

	switch val := v.(type) {
	case bool:
		m.ValidFlag = true
		m.Bool = val
		return nil
	case int64:
		if val == 0 || val == 1 {
			m.ValidFlag = true
			m.Bool = val > 0
			return nil
		}
		return fmt.Errorf(`unsafe int64 to bool cast: %d`, val)
	case string:
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		m.ValidFlag = true
		m.Bool = parsed
		return nil
	default:
		return fmt.Errorf(`unsupported cast from %s to bool`, reflect.TypeOf(v))
	}
}
