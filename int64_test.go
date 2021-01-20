package maybe

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MaybeInt64Suite struct {
	MaybeSuite
}

func (Suite *MaybeInt64Suite) TestSet() {
	var v MaybeInt64
	Suite.False(v.Valid())
	v.Set(12345)
	Suite.True(v.Valid())
	Suite.Equal(int64(12345), v.Int64)
}

func (Suite *MaybeInt64Suite) TestMarshalJSON() {
	Suite.AssertMarshalJson(MaybeInt64{}, nullBytes, nil)
	Suite.AssertMarshalJson(MaybeInt64{
		ValidFlag: false,
		Int64:     12345,
	}, nullBytes, nil)
	Suite.AssertMarshalJson(MaybeInt64{
		ValidFlag: true,
		Int64:     12345,
	}, []byte(`12345`), nil)
}

func (Suite *MaybeInt64Suite) TestUnmarshalJSON() {
	Suite.AssertUnmarshalJSON([]byte(`12345`), 12345, true, nil)
}

func (Suite *MaybeInt64Suite) TestUnmarshalJSONNull() {
	Suite.AssertUnmarshalJSON(nullBytes, 0, false, nil)
}

func (Suite *MaybeInt64Suite) TestUnmarshalJSONUnsafe() {
	Suite.AssertUnmarshalJSON([]byte(`12345.11`), 0, false,
		errors.New(`unsafe float64 to int64 cast: 12345.110000`))
}

func (Suite *MaybeInt64Suite) TestUnmarshalJSONArrayErr() {
	Suite.AssertUnmarshalJSON([]byte(`[]`), 0, false,
		errors.New(`unsupported cast from []interface {} to int64`))
}

func (Suite *MaybeInt64Suite) TestUnmarshalString() {
	Suite.AssertUnmarshalJSON([]byte(`"12345"`), 12345, true, nil)
}

func (Suite *MaybeInt64Suite) TestUnmarshalStringErr() {
	Suite.AssertUnmarshalJSON([]byte(`"not_int"`), 0, false,
		errors.New(`strconv.ParseInt: parsing "not_int": invalid syntax`))
}

func (Suite *MaybeInt64Suite) TestSetE() {
	var v MaybeInt64
	Suite.Nil(v.SetE(int64(12345)))
	Suite.Equal(int64(12345), v.Int64)
	Suite.True(v.Valid())
}

func (Suite *MaybeInt64Suite) TestNew() {
	v := Int64(64)
	Suite.True(v.Valid())
	Suite.Equal(int64(64), v.Int64)
}

func (Suite *MaybeInt64Suite) AssertUnmarshalJSON(b []byte, Expected int64, ExpectedValid bool, ExpectedErr error) {
	var v MaybeInt64
	err := json.Unmarshal(b, &v)

	Suite.Equal(Expected, v.Int64)
	Suite.Equal(ExpectedValid, v.Valid())
	if ExpectedErr == nil {
		Suite.Nil(err)
	} else {
		Suite.EqualError(err, ExpectedErr.Error())
	}
}

func TestMaybeInt64Suite(t *testing.T) {
	suite.Run(t, new(MaybeInt64Suite))
}
