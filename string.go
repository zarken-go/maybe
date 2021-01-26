package maybe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type MaybeString struct { //nolint:golint
	ValidFlag
	String string
}

func String(v string) MaybeString {
	return MaybeString{
		ValidFlag: true,
		String:    v,
	}
}

func (m MaybeString) MarshalJSON() ([]byte, error) {
	if !m.ValidFlag {
		return nullBytes, nil
	}
	return json.Marshal(m.String)
}

func (m *MaybeString) Set(v string) {
	m.ValidFlag = true
	m.String = v
}

func (m *MaybeString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		m.String = ``
		m.ValidFlag = false
		return nil
	}

	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	m.Set(v)
	return nil
}

func (m *MaybeString) SetE(v interface{}) error {
	if v == nil {
		m.ValidFlag = false
		m.String = ``
		return nil
	}

	switch val := v.(type) {
	case string:
		m.ValidFlag = true
		m.String = val
		return nil
	default:
		return fmt.Errorf(`unsupported cast from %s to string`, reflect.TypeOf(v))
	}
}
