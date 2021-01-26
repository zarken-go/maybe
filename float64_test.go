package maybe

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MaybeFloat64Suite struct {
	MaybeSuite
}

func (Suite *MaybeFloat64Suite) TestSet() {
	var v MaybeFloat64
	Suite.False(v.Valid())
	v.Set(12345.67)
	Suite.True(v.Valid())
	Suite.Equal(12345.67, v.Float64)
}

func (Suite *MaybeFloat64Suite) TestMarshalJSON() {
	Suite.AssertMarshalJSON(MaybeFloat64{}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeFloat64{
		ValidFlag: false,
		Float64:   12345,
	}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeFloat64{
		ValidFlag: true,
		Float64:   12345,
	}, []byte(`12345`), nil)
}

func (Suite *MaybeFloat64Suite) TestUnmarshalJSON() {
	Suite.AssertUnmarshalJSON([]byte(`12345`), 12345, true, nil)
	Suite.AssertUnmarshalJSON(nullBytes, 0, false, nil)
	Suite.AssertUnmarshalJSON([]byte(`12345.11`), 12345.11, true, nil)
	Suite.AssertUnmarshalJSON([]byte(`[]`), 0, false,
		errors.New(`json: cannot unmarshal array into Go value of type float64`))
	Suite.AssertUnmarshalJSON([]byte(`"12345"`), 12345, true, nil)
	Suite.AssertUnmarshalJSON([]byte(`"not_int"`), 0, false,
		errors.New(`strconv.ParseFloat: parsing "not_int": invalid syntax`))
}

func (Suite *MaybeFloat64Suite) TestSetE() {
	var v MaybeFloat64
	Suite.Nil(v.SetE(int64(12345)))
	Suite.Equal(float64(12345), v.Float64)
	Suite.True(v.Valid())

	Suite.Nil(v.SetE(nil))
	Suite.Equal(float64(0), v.Float64)
	Suite.False(v.Valid())

	Suite.Nil(v.SetE(15.00))
	Suite.Equal(float64(15), v.Float64)
	Suite.True(v.Valid())

	var v2 MaybeFloat64
	Suite.EqualError(v2.SetE([]string{}), `unsupported cast from []string to float64`)
	Suite.Equal(float64(0), v2.Float64)
	Suite.False(v2.Valid())
}

func (Suite *MaybeFloat64Suite) TestNew() {
	v := Float64(64)
	Suite.True(v.Valid())
	Suite.Equal(float64(64), v.Float64)
}

func (Suite *MaybeFloat64Suite) AssertUnmarshalJSON(b []byte, Expected float64, ExpectedValid bool, ExpectedErr error) {
	var v MaybeFloat64
	err := json.Unmarshal(b, &v)

	Suite.Equal(Expected, v.Float64)
	Suite.Equal(ExpectedValid, v.Valid())
	if ExpectedErr == nil {
		Suite.Nil(err)
	} else {
		Suite.EqualError(err, ExpectedErr.Error())
	}
}

func TestMaybeFloat64Suite(t *testing.T) {
	suite.Run(t, new(MaybeFloat64Suite))
}
