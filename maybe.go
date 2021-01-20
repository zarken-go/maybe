package maybe

import "encoding/json"

var nullBytes = []byte(`null`)

type Maybe interface {
	Valid() bool
	SetE(v interface{}) error
	json.Unmarshaler
	json.Marshaler
}

type ValidFlag bool

func (v ValidFlag) Valid() bool {
	return bool(v)
}
