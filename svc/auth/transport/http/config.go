package http

type HttpConfig struct {
	Port string `envconfig:"default=8080,PORT"`
	// CertificateHeader is the header that contains the identity.
	CertificateHeader string `envconfig:"CERTIFICATE_HEADER"`
}
