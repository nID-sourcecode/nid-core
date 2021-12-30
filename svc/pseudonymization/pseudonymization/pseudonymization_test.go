package pseudonymization

//
//import (
//	"crypto/rand"
//	"crypto/sha256"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/suite"
//	"testing"
//)
//
//type PseudonymizationTestSuite struct {
//	suite.Suite
//}
//
//func TestPseudonymizationTestSuite(t *testing.T) {
//	suite.Run(t, new(PseudonymizationTestSuite))
//}
//
//func (suite *PseudonymizationTestSuite) TestDeterminism() {
//
//	serviceHash := sha256.Sum256([]byte("Appeltaart"))
//	internalId := make([]byte, 32)
//	if _, err := rand.Read(internalId); err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//
//	pseudonym1, err := Encode(internalId, serviceHash[:])
//	if err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//	pseudonym2, err := Encode(internalId, serviceHash[:])
//	if err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//	assert.Equal(suite.T(), pseudonym1, pseudonym2)
//}
//
//func (suite *PseudonymizationTestSuite) TestDecode() {
//
//	serviceHash := sha256.Sum256([]byte("Appeltaart"))
//	internalId := make([]byte, 32)
//	if _, err := rand.Read(internalId); err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//	pseudonym, err := Encode(internalId, serviceHash[:])
//	if err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//	decodedInput, err := Decode(string(pseudonym), serviceHash[:])
//	if err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//	assert.Equal(suite.T(), decodedInput, internalId)
//}
//
//func (suite *PseudonymizationTestSuite) TestEncode() {
//
//	serviceId := sha256.Sum256([]byte("Bananen"))
//	internalId := make([]byte, 32)
//	if _, err := rand.Read(internalId); err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//	pseudonym, err := Encode(internalId, serviceId[:])
//	if err != nil {
//		assert.Error(suite.T(), err)
//	}
//
//	println(pseudonym)
//}
//
//func BenchmarkEncode(b *testing.B) {
//	for i := 0;i<b.N;i++ {
//		b.StopTimer()
//		hash := sha256.Sum256([]byte("Peren"))
//		serviceHash := hash[:]
//		internalId := make([]byte, 32)
//		if _, err := rand.Read(internalId); err != nil {
//			panic(err)
//		}
//		b.StartTimer()
//		if _, err := Encode(internalId, serviceHash); err != nil {
//			panic(err)
//		}
//	}
//}
