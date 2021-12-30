package extproc

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, &ServiceTestSuite{})
}
