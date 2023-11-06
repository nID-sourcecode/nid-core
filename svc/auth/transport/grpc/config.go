package grpc

type GrpcConfig struct {
	Port int `envconfig:"default=8081,PORT"`
	// CertificateHeader is the header that contains the identity.
	CertificateHeader string `envconfig:"CERTIFICATE_HEADER"`
}
