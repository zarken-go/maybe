package maybe

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MaybeBoolSuite struct {
	MaybeSuite
}

func (Suite *MaybeBoolSuite) TestSet() {
	var v MaybeBool
	Suite.False(v.Valid())
	v.Set(true)
	Suite.True(v.Valid())
	Suite.Equal(true, v.Bool)
}

func (Suite *MaybeBoolSuite) TestMarshalJSON() {
	Suite.AssertMarshalJSON(MaybeBool{}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeBool{
		ValidFlag: false,
		Bool:      true,
	}, nullBytes, nil)
	Suite.AssertMarshalJSON(MaybeBool{
		ValidFlag: true,
		Bool:      true,
	}, []byte(`true`), nil)
}

func (Suite *MaybeBoolSuite) TestUnmarshalJSON() {
	Suite.AssertUnmarshalJSON([]byte(`true`), true, true, nil)
	Suite.AssertUnmarshalJSON([]byte(`"0"`), false, true, nil)
	Suite.AssertUnmarshalJSON([]byte(`1`), true, true, nil)
	Suite.AssertUnmarshalJSON(nullBytes, false, false, nil)
	Suite.AssertUnmarshalJSON([]byte(`2`), false, false,
		errors.New(`unsafe int64 to bool cast: 2`))
	Suite.AssertUnmarshalJSON([]byte(`1.2`), false, false,
		errors.New(`json: cannot unmarshal number into Go value of type bool`))
	Suite.AssertUnmarshalJSON([]byte(`"not_int"`), false, false,
		errors.New(`strconv.ParseBool: parsing "not_int": invalid syntax`))
}

func (Suite *MaybeBoolSuite) TestSetE() {
	var v MaybeBool
	Suite.Nil(v.SetE(true))
	Suite.Equal(true, v.Bool)
	Suite.True(v.Valid())

	Suite.Nil(v.SetE(nil))
	Suite.Equal(false, v.Bool)
	Suite.False(v.Valid())

	Suite.Nil(v.SetE(int64(1)))
	Suite.Equal(true, v.Bool)
	Suite.True(v.Valid())

	var v2 MaybeBool
	Suite.EqualError(v2.SetE([]string{}), `unsupported cast from []string to bool`)
	Suite.Equal(false, v2.Bool)
	Suite.False(v2.Valid())
}

func (Suite *MaybeBoolSuite) TestNew() {
	v := Bool(true)
	Suite.True(v.Valid())
	Suite.Equal(true, v.Bool)
}

func (Suite *MaybeBoolSuite) AssertUnmarshalJSON(b []byte, Expected bool, ExpectedValid bool, ExpectedErr error) {
	var v MaybeBool
	err := json.Unmarshal(b, &v)

	Suite.Equal(Expected, v.Bool)
	Suite.Equal(ExpectedValid, v.Valid())
	if ExpectedErr == nil {
		Suite.Nil(err)
	} else {
		Suite.EqualError(err, ExpectedErr.Error())
	}
}

func TestMaybeBoolSuite(t *testing.T) {
	suite.Run(t, new(MaybeBoolSuite))
}
