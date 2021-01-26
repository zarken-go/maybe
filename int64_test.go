package maybe

import (
	"encoding/json"
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/suite"
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
	Suite.AssertMarshalJSON(MaybeInt64{}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeInt64{
		ValidFlag: false,
		Int64:     12345,
	}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeInt64{
		ValidFlag: true,
		Int64:     12345,
	}, []byte(`12345`), nil)
}

func (Suite *MaybeInt64Suite) TestUnmarshalJSON() {
	Suite.AssertUnmarshalJSON([]byte(`12345`), 12345, true, nil)
	Suite.AssertUnmarshalJSON([]byte(`12345.00`), 12345, true, nil)
	Suite.AssertUnmarshalJSON(nullBytes, 0, false, nil)
	Suite.AssertUnmarshalJSON([]byte(`12345.11`), 0, false,
		errors.New(`unsafe float64 to int64 cast: 12345.110000`))
	Suite.AssertUnmarshalJSON([]byte(`[]`), 0, false,
		errors.New(`json: cannot unmarshal array into Go value of type int64`))
	Suite.AssertUnmarshalJSON([]byte(`"12345"`), 12345, true, nil)
	Suite.AssertUnmarshalJSON([]byte(`"not_int"`), 0, false,
		errors.New(`strconv.ParseInt: parsing "not_int": invalid syntax`))
}

func (Suite *MaybeInt64Suite) TestSetE() {
	var v MaybeInt64
	Suite.Nil(v.SetE(int64(12345)))
	Suite.Equal(int64(12345), v.Int64)
	Suite.True(v.Valid())

	Suite.Nil(v.SetE(nil))
	Suite.Equal(int64(0), v.Int64)
	Suite.False(v.Valid())

	Suite.Nil(v.SetE(15.00))
	Suite.Equal(int64(15), v.Int64)
	Suite.True(v.Valid())

	var v2 MaybeInt64
	Suite.EqualError(v2.SetE([]string{}), `unsupported cast from []string to int64`)
	Suite.Equal(int64(0), v2.Int64)
	Suite.False(v2.Valid())

	Suite.EqualError(v2.SetE(math.Inf(1)), `unsafe float64 to int64 cast: +Inf`)
	Suite.EqualError(v2.SetE(math.Inf(-1)), `unsafe float64 to int64 cast: -Inf`)
	Suite.EqualError(v2.SetE(math.NaN()), `unsafe float64 to int64 cast: NaN`)
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
