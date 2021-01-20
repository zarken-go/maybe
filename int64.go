package maybe

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type MaybeInt64 struct {
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
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	return m.SetE(v)
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
		// TODO: Test +INF, -INF, NaN
		// verify lossless conversion
		if val == math.Trunc(val) {
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
