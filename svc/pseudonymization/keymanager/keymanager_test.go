package keymanager

//
//import (
//	"fmt"
//	"github.com/stretchr/testify/suite"
//	"lab.weave.nl/nid/nid-core/services/pseudonymization/test"
//	"testing"
//	"time"
//)
//
//type KeymanagerTestSuite struct {
//	suite.Suite
//}
//
//func TestKeymanagerTestSuite(t *testing.T) {
//	suite.Run(t, new(KeymanagerTestSuite))
//}
//
//func (suite *KeymanagerTestSuite) TestGetKey() {
//	port := "54445"
//	s := test.RunJwkServer(port)
//
//	k := NewKeyManager(fmt.Sprintf("http://localhost:%s/{{namespace}}",port), 24 * time.Hour)
//
//	suite.Equal(0,s.ReceivedRequests)
//	key1, err := k.GetKey("alice")
//	if err != nil {
//		suite.Fail(err.Error())
//		return
//	}
//
//	suite.Equal(1,s.ReceivedRequests)
//
//	key2, err := k.GetKey("alice")
//	if err != nil {
//		suite.Fail(err.Error())
//		return
//	}
//
//	suite.Equal(key1, key2)
//
//	suite.Equal(1,s.ReceivedRequests)
//
//	_, err = k.GetKey("bob")
//	if err != nil {
//		suite.Fail(err.Error())
//		return
//	}
//
//	suite.Equal(2,s.ReceivedRequests)
//
//	k.Cleanup()
//}
