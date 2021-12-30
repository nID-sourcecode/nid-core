package keyutil

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

type KeyutilTestSuite struct {
	suite.Suite
}

func (s *KeyutilTestSuite) TestErrorOn512BitRSAKey() {
	smallKey := `-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBALeUvqLo1dKixuDHDdazf2IpQD/SqWn42IXTq1E8kWdYOqrUgEHj
Q20fdHoI23QEH48SS83k/fcARfrMs9tSm+kCAwEAAQJBAKvsKggo4y5K8NNKtyQN
n8sO9LOQlsW+nQ/fZf5DKazMNothn5IAjnEDxEnmFca4Ki18rVJszw7cvxwDWm37
RgECIQDf4KNlOsAsW26M9271ryXhMF02IzeGQWZ87M/blmafqQIhANHr9Lzc94iq
g6wAtUC5cZzKyOkbTOWSo7mkvKVX1sJBAiAoUdG7moAfvPvFAY8HSlr9GnO/G0qV
sFOf7hplRsoGuQIgEqh3Q0YclkAZnfMeKReSeo4nl1h+2DTVao2y2rtY8kECIQDB
P27WykDtc6w1Vbx2LjC0NgAHq8nGNPBlaX5aoPR8JA==
-----END RSA PRIVATE KEY-----
`

	_, err := ParseKeypair(smallKey)

	s.Require().Error(err)
	s.Require().True(errors.Is(err, ErrInvalidRsaKeyLength))
}

func (s *KeyutilTestSuite) TestNoErrorOn2048BitRSAKey() {
	bigKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAr1mGVndB0gOaZ+t56/H7Eq5ye/SngEXlXj86vlG5krlstaxR
sqHjW6kOriMLh3W9rPRftBsUPfOOLyhARagKUJI0/rSnmiU+WKFdDCCVtooxitiE
LeqdLL32HPGUUtQN0wFS9mmrk+r1yQJH/QGnjhx6C5tyNfQkK/cZZGLDR+OPFbYr
h1hFxwrpI61OFM5DvvSsOzZQ3MWYjB0X3KQopbXMS+nU2L+DLFKErgg5/qoGV1tX
swlpYdZFMAABruljQ3NISTxv8HDjc2xwfM46oNLukpiLGkR2lU/33pAVEgM3p54j
FON4vJtxSJczrEG3Io4Z8MZ4PJl7fRCn3DIl4wIDAQABAoIBAQChIxLpaIRK+1as
Qt9yrJc+TqMUN/qpTRH/rvlLpgxzwgQdWzJkhpJJTC0aZ0gT/mYEhzlfaDcMy0GR
IlsV59s6uXEL03XlmG4XwomgOF18NrhaUKf/DgfL8vE1Hedgyk+95QyZNvgeYR/m
zTrcTOXuGUsqWXn2DNoksNlbv8qWxwjf+TynvLUi63cTYJhCyS0HXuoiGPbn6KHP
DNxy4JOnCj15YeIbgg87cNOq0aC5M4xxzWEteuVDYYEQWQmbAiYzuj+AA1XG9ivX
nRTesE9gvmstHoaG0lwQXKuOAty7zEJIwx7/jteVB/UerDxeHDO/+ciN54ZOjGmG
GpdE8YkBAoGBANfxkHiRTEBz4+o7K3D0Uejt3UPIw6gWHYixKE0mtmwPaqKLcU6r
3DkZL1CUi+XMV7EW2nAkWUOuOCuvWsCbjBggRFFV2ReRJnWYy/iTVH9dBLYA2OFQ
p+qrhQbRff1oUKBJBteckAQWSQXtcqf+KRnRIDAXR5Jw7WEOX5BoGJBjAoGBAM/g
St1f/qlLPlCzjZoAhuW4PNikA+O77xJPL2uESd+HmWvIwy+g/uhNOknrh6EXyqsY
CZDCFVL7wwuH8HiTwje1ODEnomsamnz2RhznqfftP7YiLUhQCo0z0ObHRaAE5aEh
CwJUQl8T6oftPjLWogYrhOQBlGfJaGfnaFwvfEyBAoGATOKfU643/gLFNVKX5wG2
YD7Aty+2KhSls1OQS9fqv5LFntYTI7WhFVtYM1KQdONKnazLXX4zohtXuIYYw9ce
DEEA0gzE3NU7Ykdi6EBcp3RRBxRKI/75ql4jYQgZ2a3YdxlJLF98D1h363pdhl7B
94Uz9qtzOjqm6hWaBOprRI0CgYEApZ0naAmTxWLqCbeTaA9lad2HtH2vj59pz+eA
eyNRC6Jny+SOBQM6Mu9cMgpQ6zoeQIONE2RdQtjLwwMRxa7KvEFHvHm8P6JZVJeM
snirBJhi+wNtmkASt/6BP2uhf+SG4gGGWNuyaTdf0d1kgXJYcZv4awMLLkjbQnSt
w0wdtoECgYAkWepoAcSCiK7rofHo0qzAuGa3A+DSpNBVFVPYtmtEZB4+0ygZZcRH
U79zztnL6QPVPSCh+/3ZN9clcyIwgwqFHEKzQSAs6gq+pSttYx2cE3uBlUGR/0f2
/BIm+h+Xn50Sol/efVaGg3novBRC0Nz2eYvQvFWrTmr4dllLDNX3pQ==
-----END RSA PRIVATE KEY-----
`

	_, err := ParseKeypair(bigKey)

	s.Require().NoError(err)
}

func TestKeyutilTestSuite(t *testing.T) {
	suite.Run(t, &KeyutilTestSuite{})
}
