package kubeutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	istio "istio.io/client-go/pkg/clientset/versioned"
	fakeIstio "istio.io/client-go/pkg/clientset/versioned/fake"
	fakeVirtualService "istio.io/client-go/pkg/clientset/versioned/typed/networking/v1beta1/fake"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	fakeservice "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stesting "k8s.io/client-go/testing"
)

// Error definitions
var (
	ErrTest error = fmt.Errorf("alles is lek")
)

type KubeUtilServiceTestSuite struct {
	suite.Suite
}

func (s *KubeUtilServiceTestSuite) TestCreateService() {
	data := []struct {
		testName              string
		inputExternalHostname string
		inputServiceName      string
		inputNamespace        string
		inputServicePort      int32
		err                   bool
	}{
		{
			testName:              "originalCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputServicePort:      80,
			err:                   false,
		},
		{
			testName:              "duplicateCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputServicePort:      80,
			err:                   true,
		},
	}

	client := NewKubeUtil(fake.NewSimpleClientset(), fakeIstio.NewSimpleClientset())
	for _, test := range data {
		s.Run(test.testName, func() {
			err := client.CreateService(context.Background(), test.inputExternalHostname, test.inputServiceName, test.inputNamespace, test.inputServicePort)
			if test.err {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KubeUtilServiceTestSuite) TestCreateServiceAndDeployment() {
	data := []struct {
		testName              string
		inputExternalHostname string
		inputServiceName      string
		inputNamespace        string
		inputDockerImage      string
		inputServicePort      int32
		inputEnv              map[string]string
		err                   bool
	}{
		{
			testName:              "originalCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputDockerImage:      "testDockerImage",
			inputEnv:              map[string]string{"VAR": "TEST"},
			inputServicePort:      80,
			err:                   false,
		},
		{
			testName:              "duplicateCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputDockerImage:      "testDockerImage",
			inputEnv:              map[string]string{"VAR": "TEST"},
			inputServicePort:      80,
			err:                   true,
		},
	}

	client := NewKubeUtil(fake.NewSimpleClientset(), fakeIstio.NewSimpleClientset())
	for _, test := range data {
		s.Run(test.testName, func() {
			err := client.CreateServiceAndDeployment(context.Background(), test.inputNamespace, test.inputServiceName, test.inputExternalHostname, test.inputServicePort, test.inputDockerImage, test.inputEnv)
			if test.err {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KubeUtilServiceTestSuite) TestCreateVirtualService() {
	data := []struct {
		testName              string
		inputExternalHostname string
		inputServiceName      string
		inputNamespace        string
		inputServicePort      int32
		err                   bool
	}{
		{
			testName:              "originalCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputServicePort:      80,
			err:                   false,
		},
		{
			testName:              "duplicateCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputServicePort:      80,
			err:                   true,
		},
		{
			testName:              "emptyExternalHostname",
			inputExternalHostname: "",
			inputServiceName:      "testServiceName2",
			inputNamespace:        "testNamespace2",
			inputServicePort:      81,
			err:                   true,
		},
	}

	client := NewKubeUtil(fake.NewSimpleClientset(), fakeIstio.NewSimpleClientset())
	for _, test := range data {
		s.Run(test.testName, func() {
			err := client.CreateVirtualService(context.Background(), test.inputExternalHostname, test.inputServiceName, test.inputNamespace, test.inputServicePort)
			if test.err {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KubeUtilServiceTestSuite) TestCreateDeployment() {
	data := []struct {
		testName              string
		inputExternalHostname string
		inputServiceName      string
		inputNamespace        string
		inputDockerImage      string
		inputServicePort      int32
		inputEnv              map[string]string
		err                   bool
	}{
		{
			testName:              "originalCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputDockerImage:      "testDockerImage",
			inputServicePort:      80,
			inputEnv:              map[string]string{"VAR": "TEST"},
			err:                   false,
		},
		{
			testName:              "duplicateCreate",
			inputExternalHostname: "testExternalHostname",
			inputServiceName:      "testServiceName",
			inputNamespace:        "testNamespace",
			inputDockerImage:      "testDockerImage",
			inputServicePort:      80,
			err:                   true,
		},
	}

	client := NewKubeUtil(fake.NewSimpleClientset(), fakeIstio.NewSimpleClientset())
	for _, test := range data {
		s.Run(test.testName, func() {
			err := client.CreateDeployment(context.Background(), test.inputExternalHostname, test.inputServiceName, test.inputNamespace, test.inputDockerImage, test.inputServicePort, test.inputEnv)
			if test.err {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KubeUtilServiceTestSuite) TestListNamespace() {
	data := []struct {
		clientset               kubernetes.Interface
		countExpectedNamespaces int
		err                     error
	}{
		{
			clientset:               fake.NewSimpleClientset(),
			countExpectedNamespaces: 0,
		},
		{
			clientset: fake.NewSimpleClientset(&v1.NamespaceList{
				Items: []v1.Namespace{{ObjectMeta: v12.ObjectMeta{Name: "test"}}},
			}),
			countExpectedNamespaces: 1,
		},
	}

	for _, test := range data {
		s.Run("", func() {
			client := NewKubeUtil(test.clientset, fakeIstio.NewSimpleClientset())
			namespaces, err := client.ListNamespace(context.Background())
			s.Equal(test.err, err)
			s.Len(namespaces, test.countExpectedNamespaces)
		})
	}
}

func (s *KubeUtilServiceTestSuite) TestDeleteService() {
	data := []struct {
		getKubernetesClient func() kubernetes.Interface
		getIstioClient      func() istio.Interface
		inputNamespace      string
		inputName           string
		err                 bool
	}{
		{
			getKubernetesClient: func() kubernetes.Interface { return fake.NewSimpleClientset() },
			getIstioClient:      func() istio.Interface { return fakeIstio.NewSimpleClientset() },
			inputNamespace:      "",
			inputName:           "",
			err:                 true,
		},
		{
			getKubernetesClient: func() kubernetes.Interface {
				fakeKubeClient := fake.NewSimpleClientset()
				fakeKubeClient.CoreV1().Services("").(*fakeservice.FakeServices).Fake.PrependReactor("*", "*", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, nil
				})

				return fakeKubeClient
			},
			getIstioClient: func() istio.Interface {
				fakeIstioClient := fakeIstio.NewSimpleClientset()
				fakeIstioClient.NetworkingV1beta1().VirtualServices("").(*fakeVirtualService.FakeVirtualServices).Fake.PrependReactor("*", "*", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, nil
				})

				return fakeIstioClient
			},
			inputNamespace: "",
			inputName:      "",
		},
		{
			getKubernetesClient: func() kubernetes.Interface {
				fakeKubeClient := fake.NewSimpleClientset()
				fakeKubeClient.CoreV1().Services("").(*fakeservice.FakeServices).Fake.PrependReactor("*", "*", func(action k8stesting.Action) (bool, runtime.Object, error) {
					return true, nil, nil
				})

				return fakeKubeClient
			},
			getIstioClient: func() istio.Interface {
				fakeIstioClient := fakeIstio.NewSimpleClientset()
				fakeIstioClient.NetworkingV1beta1().VirtualServices("").(*fakeVirtualService.FakeVirtualServices).Fake.PrependReactor("*", "*", func(action k8stesting.Action) (bool, runtime.Object, error) {
					if action.GetResource().Resource == "virtualservices" {
						return true, nil, ErrTest
					}

					return true, nil, nil
				})

				return fakeIstioClient
			},
			inputNamespace: "",
			inputName:      "",
			err:            true,
		},
	}

	for _, test := range data {
		s.Run("", func() {
			client := NewKubeUtil(test.getKubernetesClient(), test.getIstioClient())
			err := client.DeleteService(context.Background(), test.inputNamespace, test.inputName)
			if test.err {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *KubeUtilServiceTestSuite) TestListService() {
	data := []struct {
		clientset             kubernetes.Interface
		countExpectedServices int
		inputNamespace        string
		err                   error
	}{
		{
			clientset: fake.NewSimpleClientset(),
		},
		{
			clientset: fake.NewSimpleClientset(&v1.ServiceList{
				Items: []v1.Service{{ObjectMeta: v12.ObjectMeta{Name: "test", Namespace: "namespace1"}}},
			}),
			countExpectedServices: 1,
		},
	}

	for _, test := range data {
		s.Run("", func() {
			client := NewKubeUtil(test.clientset, fakeIstio.NewSimpleClientset())
			services, err := client.ListService(context.Background(), test.inputNamespace)
			s.Equal(test.err, err)
			s.Len(services.Items, test.countExpectedServices)
		})
	}
}

func TestKubeUtilServiceTestSuite(t *testing.T) {
	suite.Run(t, &KubeUtilServiceTestSuite{})
}
