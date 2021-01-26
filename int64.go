package maybe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type MaybeInt64 struct { //nolint:golint
	ValidFlag
	Int64 int64
}

func Int64(v int64) MaybeInt64 {
	return MaybeInt64{
		ValidFlag: true,
		Int64:     v,
	}
}

func (m MaybeInt64) MarshalJSON() ([]byte, error) {
	if !m.ValidFlag {
		return nullBytes, nil
	}
	return json.Marshal(m.Int64)
}

func (m *MaybeInt64) Set(v int64) {
	m.ValidFlag = true
	m.Int64 = v
}

func (m *MaybeInt64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		m.Int64 = 0
		m.ValidFlag = false
		return nil
	}

	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		var s string
		if json.Unmarshal(b, &s) == nil {
			return m.SetE(s)
		}

		var f float64
		if json.Unmarshal(b, &f) == nil {
			return m.SetE(f)
		}

		return err
	}
	m.Set(v)
	return nil
}

func (m *MaybeInt64) SetE(v interface{}) error {
	if v == nil {
		m.ValidFlag = false
		m.Int64 = 0
		return nil
	}

	switch val := v.(type) {
	case int64:
		m.ValidFlag = true
		m.Int64 = val
		return nil
	case float64:
		if val == math.Trunc(val) &&
			!math.IsInf(val, -1) &&
			!math.IsInf(val, 1) &&
			!math.IsNaN(val) {
			m.ValidFlag = true
			m.Int64 = int64(val)
			return nil
		}
		return fmt.Errorf(`unsafe float64 to int64 cast: %f`, val)
	case string:
		parsed, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		m.ValidFlag = true
		m.Int64 = parsed
		return nil
	default:
		return fmt.Errorf(`unsupported cast from %s to int64`, reflect.TypeOf(v))
	}
}
