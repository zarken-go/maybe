package maybe

import (
	"encoding/json"

	"github.com/stretchr/testify/suite"
)

type MaybeSuite struct {
	suite.Suite
}

func (Suite *MaybeSuite) AssertMarshalJSON(v interface{}, Expected []byte, ExpectedErr error) {
	b, err := json.Marshal(v)
	Suite.Equal(Expected, b)
	Suite.Equal(ExpectedErr, err)
}
