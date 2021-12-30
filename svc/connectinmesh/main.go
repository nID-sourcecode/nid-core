package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	httpUtils "lab.weave.nl/nid/nid-core/pkg/utilities/http"
	"lab.weave.nl/nid/nid-core/pkg/utilities/log/v2"
)

const (
	caCert = `
-----BEGIN CERTIFICATE-----
MIIDHDCCAgQCCQDlLNgEHua6+zANBgkqhkiG9w0BAQsFADBQMR0wGwYDVQQKDBRO
aUQgRXhhbXBsZSBDbHVzdGVyLjEvMC0GA1UEAwwmY29ubmVjdGlubWVzaC5uaWQu
ZXhhbXBsZS5uLWlkLm5ldHdvcmswHhcNMjAwNzI3MDc1MTI0WhcNMjEwNzI3MDc1
MTI0WjBQMR0wGwYDVQQKDBROaUQgRXhhbXBsZSBDbHVzdGVyLjEvMC0GA1UEAwwm
Y29ubmVjdGlubWVzaC5uaWQuZXhhbXBsZS5uLWlkLm5ldHdvcmswggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQC4pUppFAA+9t56Cl0xfn7JK/JHYVwck+NN
Y+K4WH01rE1x2zv9xQ2+UWQnGFvVPWPFLoCVdkLgrsebJifRUz5v9WGYQ4YQZYHd
XjHi8GRg0h6yQz0zOGze8ZgvCbnP62IUdN+U8CaoNEUS7aydMZCnrYokgafAu/rD
lHtS6KzDQzOhfA3viKnqLpxEDKNwWIyiXNrKHUFoQFQdPD2zT7kddzZhVdQT5GHr
Jy2v5fikfZ/pIAUeIEK1RN6+m8Ujdf4u8rs7Uegv+L6QctIEoL8HEl4i+/RalCwS
2ebAmdzGVcdxXr9c/hwnxwLcwQnRlP1wByx2H4Wizbhr35J5WzgfAgMBAAEwDQYJ
KoZIhvcNAQELBQADggEBAJCSZSi0ZHxCKnyst+/2vF8PY2VN/Wihbp1FV6LCIrbO
lR1lKn5FGVDKfrcT/ZCehIPk+NxrlcVFGq+O8795OlBodUF780yVa500g2dyyiGi
XvMkcPdImdAudC9hzy9xvjcHz/3TMQ4503eIc8SKsN+qhSVY/EdP+mVMnMhy/Wys
IkHsaSC6p3svas1dSk6+nMkva59NAIpAzEaFkX9Klrgok1G9lxZNEktO+XTa2WEy
V714YTl1wBwHi5vrAlQziLNiUJ9aaksvYlZATC58s85pBFK2dNBGVHBNwgm0Omg+
e4HNLT0j0vSxPnmfSrFySsVGfqj7SizFWIA9yrkqIEQ=
-----END CERTIFICATE-----
`
	requestTimeOut time.Duration = 10 * time.Second
)

func main() {
	// FIXME: use a pool with only the active client certs https://lab.weave.nl/nid/nid-core/-/issues/63

	go doRequests()

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM([]byte(caCert))

	// nolint: gosec
	tlsConfig := &tls.Config{
		ClientCAs:  caPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      ":80",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/v1/health", healthhandler)
	http.HandleFunc("/hello", worldHandler)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func healthhandler(w http.ResponseWriter, req *http.Request) {
}

func worldHandler(w http.ResponseWriter, req *http.Request) {
	// FIXME: https://lab.weave.nl/nid/nid-core/-/issues/65 we can verify the client certificate here, maybe we can do this in istio gateway
	for key, values := range req.Header {
		for i, val := range values {
			_, err := fmt.Fprintf(w, "%s[%d]:%s\n", key, i, val)
			if err != nil {
				panic("Unexpected error")
			}
		}
	}
}

func doRequests() {
	httpClient := &http.Client{}
	for {
		time.Sleep(requestTimeOut)
		host := os.Getenv("OUTSIDE_MESH")
		if host != "" {
			performRequest(httpClient, host)
		} else {
			log.Info("host is empty")
		}
	}
}

func performRequest(client httpUtils.Client, host string) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://"+host+"/hello", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		return
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.WithError(err).Error("unable to close response body")
		}
	}()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)

		return
	}
	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v", err)

		return
	}
	if string(respString) != "World!" {
		fmt.Printf("Unexpected response")

		return
	}
	log.Infof("Success: %s", host)
}
