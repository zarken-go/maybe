package maybe

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MaybeStringSuite struct {
	MaybeSuite
}

func (Suite *MaybeStringSuite) TestSet() {
	var v MaybeString
	Suite.False(v.Valid())
	v.Set(`str`)
	Suite.True(v.Valid())
	Suite.Equal(`str`, v.String)
}

func (Suite *MaybeStringSuite) TestMarshalJSON() {
	Suite.AssertMarshalJSON(MaybeString{}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeString{
		ValidFlag: false,
		String:    `str`,
	}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeString{
		ValidFlag: true,
		String:    `str`,
	}, []byte(`"str"`), nil)
}

func (Suite *MaybeStringSuite) TestUnmarshalJSON() {
	Suite.AssertUnmarshalJSON([]byte(`"str"`), `str`, true, nil)
	Suite.AssertUnmarshalJSON(nullBytes, ``, false, nil)
	Suite.AssertUnmarshalJSON([]byte(`12345.11`), ``, false,
		errors.New(`json: cannot unmarshal number into Go value of type string`))
	Suite.AssertUnmarshalJSON([]byte(`"12345"`), `12345`, true, nil)
}

func (Suite *MaybeStringSuite) TestSetE() {
	var v MaybeString
	Suite.Nil(v.SetE(`str`))
	Suite.Equal(`str`, v.String)
	Suite.True(v.Valid())

	Suite.Nil(v.SetE(nil))
	Suite.Equal(``, v.String)
	Suite.False(v.Valid())

	Suite.EqualError(v.SetE(1234.23), `unsupported cast from float64 to string`)
}

func (Suite *MaybeStringSuite) TestNew() {
	v := String(`str`)
	Suite.True(v.Valid())
	Suite.Equal(`str`, v.String)
}

func (Suite *MaybeStringSuite) AssertUnmarshalJSON(b []byte, Expected string, ExpectedValid bool, ExpectedErr error) {
	var v MaybeString
	err := json.Unmarshal(b, &v)

	Suite.Equal(Expected, v.String)
	Suite.Equal(ExpectedValid, v.Valid())
	if ExpectedErr == nil {
		Suite.Nil(err)
	} else {
		Suite.EqualError(err, ExpectedErr.Error())
	}
}

func TestMaybeStringSuite(t *testing.T) {
	suite.Run(t, new(MaybeStringSuite))
}
