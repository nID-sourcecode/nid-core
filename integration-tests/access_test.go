// +build integration to files

package integration

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/vrischmann/envconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/portforward"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type AccessTestSuite struct {
	suite.Suite
	clientset *kubernetes.Clientset
	conf      AccessTestSuiteConfig

	systemNamespaceClient  *resty.Client
	testingNamespaceClient *resty.Client
}

type AccessTestSuiteConfig struct {
	Namespace string `envconfig:"NAMESPACE"`
}

func TestAccessTestSuite(t *testing.T) {
	suite.Run(t, &AccessTestSuite{})
}

func (s *AccessTestSuite) SetupSuite() {
	s.conf = AccessTestSuiteConfig{}
	err := envconfig.Init(&s.conf)

	s.Require().NoError(err, "reading env")

	factory := cmdutil.NewFactory(&genericclioptions.ConfigFlags{})
	clientset, err := factory.KubernetesClientSet()
	s.Require().NoError(err)

	s.clientset = clientset

	s.systemNamespaceClient = s.createProxyClient(s.conf.Namespace, "8081", "testproxy")
	s.testingNamespaceClient = s.createProxyClient("testing", "8082", "testproxy-testing")
}

func (s *AccessTestSuite) TestSystemAutopseudoAccessibleFromWithinNamespace() {
	req := s.systemNamespaceClient.NewRequest()
	res, err := req.Get("http://autopseudo." + s.conf.Namespace + "/decryptAndApply")

	s.Require().NoError(err)
	s.Equal(http.StatusOK, res.StatusCode(), "%d: %s", res.StatusCode(), res.String())
}

func (s *AccessTestSuite) TestSystemAutopseudoNotAccessibleFromOutsideNamespace() {
	req := s.testingNamespaceClient.NewRequest()
	res, err := req.Get("http://autopseudo." + s.conf.Namespace + "/decryptAndApply")

	s.Require().NoError(err)
	s.Equal(http.StatusForbidden, res.StatusCode(), "%d: %s", res.StatusCode(), res.String())
}

func (s *AccessTestSuite) TestTestingAutopseudoAccessibleFromWithinNamespace() {
	req := s.testingNamespaceClient.NewRequest()
	res, err := req.Get("http://autopseudo.testing/decryptAndApply")

	s.Require().NoError(err)
	s.Equal(http.StatusOK, res.StatusCode(), "%d: %s", res.StatusCode(), res.String())
}

func (s *AccessTestSuite) TestTestingAutopseudoNotAccessibleFromOutsideNamespace() {
	req := s.systemNamespaceClient.NewRequest()
	res, err := req.Get("http://autopseudo.testing/decryptAndApply")

	s.Require().NoError(err)
	s.Equal(http.StatusForbidden, res.StatusCode(), "%d: %s", res.StatusCode(), res.String())
}

func (s *AccessTestSuite) createProxyClient(namespace string, port string, appname string) *resty.Client {
	podList, err := s.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=" + appname,
	})
	s.Require().NoError(err)

	s.Require().Len(podList.Items, 1)

	podName := podList.Items[0].Name

	factory := cmdutil.NewFactory(&genericclioptions.ConfigFlags{Namespace: &namespace})

	streams, _, out, _ := genericclioptions.NewTestIOStreams()
	portforwarder := portforward.NewCmdPortForward(factory, streams)
	portforwarder.SetArgs([]string{podName, port + ":8081"})

	go func() {
		s.Require().NoError(portforwarder.Execute())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	forwarded := false
	for !forwarded {
		select {
		case <-ctx.Done():
			s.FailNow("port forward timed out")
		default:
			s.Require().NoError(err)
			outBytes, _ := ioutil.ReadAll(out)
			log.Info("Checking port-forward status")
			if strings.Contains(string(outBytes), "Forwarding from") {
				forwarded = true
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	client := resty.New()
	client.SetProxy("http://localhost:" + port)
	return client
}
