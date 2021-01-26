package maybe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type MaybeFloat64 struct { //nolint:golint
	ValidFlag
	Float64 float64
}

func Float64(v float64) MaybeFloat64 {
	return MaybeFloat64{
		ValidFlag: true,
		Float64:   v,
	}
}

func (m MaybeFloat64) MarshalJSON() ([]byte, error) {
	if !m.ValidFlag {
		return nullBytes, nil
	}
	return json.Marshal(m.Float64)
}

func (m *MaybeFloat64) Set(v float64) {
	m.ValidFlag = true
	m.Float64 = v
}

func (m *MaybeFloat64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		m.Float64 = 0
		m.ValidFlag = false
		return nil
	}

	var v float64
	if err := json.Unmarshal(b, &v); err != nil {
		var s string
		if json.Unmarshal(b, &s) == nil {
			return m.SetE(s)
		}

		return err
	}
	m.Set(v)
	return nil
}

func (m *MaybeFloat64) SetE(v interface{}) error {
	if v == nil {
		m.ValidFlag = false
		m.Float64 = 0
		return nil
	}

	switch val := v.(type) {
	case float64:
		m.ValidFlag = true
		m.Float64 = val
		return nil
	case int64:
		m.ValidFlag = true
		m.Float64 = float64(val)
		return nil
	case string:
		parsed, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		m.ValidFlag = true
		m.Float64 = parsed
		return nil
	default:
		return fmt.Errorf(`unsupported cast from %s to float64`, reflect.TypeOf(v))
	}
}
