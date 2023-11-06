package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/machinebox/graphql"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	security "istio.io/api/security/v1beta1"
	typeb1 "istio.io/api/type/v1beta1"
	networkingv1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	"istio.io/client-go/pkg/apis/security/v1beta1"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/nID-sourcecode/nid-core/pkg/utilities/errors"
	grpcerrors "github.com/nID-sourcecode/nid-core/pkg/utilities/grpcserver/errors"
	"github.com/nID-sourcecode/nid-core/pkg/utilities/log/v2"
	"github.com/nID-sourcecode/nid-core/svc/dashboard/kubeutil"
	"github.com/nID-sourcecode/nid-core/svc/dashboard/proto"
)

// ErrInvalidNamespace is returned when an operation on the specified namespace is not allowed (operations on kube-system for example)
var ErrInvalidNamespace = errors.New("invalid namespace")

const audienceQuery = `query audiences($filter: AudienceFilterInput!) {
  audiences(filter: $filter) {
    id
  }
}`

const createAudienceMutation = `
mutation createAudience($audience: CreateAudience!) {
  createAudience(input: $audience) {
    id
  }
}
`

// CreateAudienceInput contains the information that is needed to create an audience
type CreateAudienceInput struct {
	Audience  string `json:"audience"`
	Namespace string `json:"namespace"`
}

type audienceResponse struct {
	Audiences []struct {
		ID string
	}
}

const uriScheme = "http://{{service}}.{{namespace}}{{port}}{{gqlUri}}"
const (
	httpPort = 80
)

// DashboardServiceServer implementation of the dashboard server
type DashboardServiceServer struct {
	config             *DashBoardConfig
	kubeClientSet      *kubernetes.Clientset
	istioClientSet     *istio.Clientset
	registrySecretJSON []byte
	kubeutil           kubeutil.Interface
}

// NewDashboardServiceServer creates a new dashboard service server
func NewDashboardServiceServer(config *DashBoardConfig) (*DashboardServiceServer, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error reading in-cluster config")
	}
	// creates the clientset
	kubeClientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error creating kubernetes clientset")
	}

	istioClientset, err := istio.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error creating istio clientset")
	}

	kubeUtil := kubeutil.NewKubeUtil(kubeClientset, istioClientset)

	var registrySecretJSON []byte
	registrySecretJSON, err = ioutil.ReadFile(config.RegistrySecretPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.Wrapf(err, "could find registrysecret in %s (REGISTRY_SECRET_PATH)", config.RegistrySecretPath)
		}
		return nil, errors.Wrapf(err, "error reading registry secret from %s (REGISTRY_SECRET_PATH)", config.RegistrySecretPath)
	}

	return &DashboardServiceServer{
		kubeClientSet:      kubeClientset,
		istioClientSet:     istioClientset,
		registrySecretJSON: registrySecretJSON,
		kubeutil:           kubeUtil,
		config:             config,
	}, nil
}

func (s *DashboardServiceServer) createAutoPseudoIfNotExists(ctx context.Context, namespace string) error {
	existing, err := s.kubeClientSet.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=autopseudo",
	})
	if err != nil {
		return errors.Wrap(err, "unable to list services for auto pseudo")
	}
	if len(existing.Items) > 0 {
		// skip creation, already exists
		return nil
	}

	err = s.kubeutil.CreateService(ctx, "", "autopseudo", namespace, 80)
	if err != nil {
		return errors.Wrap(err, "unable to create auto pseudo service")
	}

	err = s.kubeutil.CreateDeployment(ctx, "", "autopseudo", namespace, s.config.AutopseudoImage, 80, map[string]string{"NAMESPACE": namespace})
	if err != nil {
		return errors.Wrap(err, "unable to create deployment for auto pseudo")
	}

	authPolicy := &v1beta1.AuthorizationPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind: "AuthorizationPolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "autopseudo",
		},
		Spec: security.AuthorizationPolicy{
			Selector: &typeb1.WorkloadSelector{
				MatchLabels: map[string]string{"app": "autopseudo"},
			},
			Action: security.AuthorizationPolicy_ALLOW,
			Rules: []*security.Rule{
				{
					From: []*security.Rule_From{{
						Source: &security.Source{
							Namespaces: []string{namespace},
						},
					}},
					To: []*security.Rule_To{{
						Operation: &security.Operation{
							Paths:   []string{"/decryptAndApply"},
							Methods: []string{"GET", "POST", "PUT", "DELETE"},
						},
					}},
				},
				{
					To: []*security.Rule_To{{
						Operation: &security.Operation{
							Paths:   []string{"/jwks"},
							Methods: []string{"GET"},
						},
					}},
				},
			},
		},
	}
	if _, err := s.istioClientSet.SecurityV1beta1().AuthorizationPolicies(namespace).Create(ctx, authPolicy, metav1.CreateOptions{}); err != nil {
		return errors.Wrap(err, "unable to create authorization policy")
	}

	return nil
}

// ListNamespaces lists all namespaces for current kube client set
func (s *DashboardServiceServer) ListNamespaces(ctx context.Context, _ *emptypb.Empty) (*proto.NamespaceList, error) {
	namespaces, err := s.kubeutil.ListNamespace(ctx)
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to list namesapces")

		return nil, grpcerrors.ErrInternalServer()
	}

	// FIXME we might want to filter `kube-...` and `isito-system` namespaces for example if we ever expose this to end users
	return &proto.NamespaceList{
		Items: namespaces,
	}, nil
}

func (s *DashboardServiceServer) createNamespaceIfNotExists(ctx context.Context, namespace string) error {
	log.WithField("namespace", namespace).Info("Checking if ns exists")
	existing, err := s.kubeClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", namespace),
	})
	if err != nil {
		return errors.Wrap(err, "unable to list namespaces")
	}
	if len(existing.Items) > 0 {
		log.WithField("namespace", namespace).Info("Namespace exists skipping creation")
		// Already exists, skip namespace creation
		return nil
	}

	ns := v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind: "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
			Labels: map[string]string{
				"name":            namespace,
				"istio-injection": "enabled",
			},
		},
	}
	if _, err := s.kubeClientSet.CoreV1().Namespaces().Create(ctx, &ns, metav1.CreateOptions{}); err != nil {
		return errors.Wrap(err, "unable to create namespace")
	}

	registrySecret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind: "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "registry-weave",
			Namespace: namespace,
		},
		Data: map[string][]byte{
			".dockerconfigjson": s.registrySecretJSON,
		},
		Type: "kubernetes.io/dockerconfigjson",
	}
	if _, err := s.kubeClientSet.CoreV1().Secrets(namespace).Create(ctx, &registrySecret, metav1.CreateOptions{}); err != nil {
		return errors.Wrap(err, "unable to create registry secret")
	}

	gateway := networkingv1.Gateway{
		TypeMeta: metav1.TypeMeta{
			Kind: "Gateway",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "external-access-gateway",
		},
		Spec: networkingv1beta1.Gateway{
			Selector: map[string]string{
				"istio": "ingressgateway",
			},
			Servers: []*networkingv1beta1.Server{
				{
					Port: &networkingv1beta1.Port{
						Number:   httpPort,
						Name:     "http",
						Protocol: "HTTP",
					},
					Hosts: []string{fmt.Sprintf("*.%s.%s", namespace, s.config.BaseDomain)},
				},
			},
		},
	}

	_, err = s.istioClientSet.NetworkingV1beta1().Gateways(namespace).Create(ctx, &gateway, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to create gatewayc")
	}

	return nil
}

// DeployService deploys a docker image in given namespace
func (s *DashboardServiceServer) DeployService(ctx context.Context, req *proto.DeployServiceRequest) (*proto.DeployServiceResponse, error) {
	if err := s.createNamespaceIfNotExists(ctx, req.GetNamespace()); err != nil {
		log.Extract(ctx).WithError(err).Error("unable to create namespace")

		return nil, grpcerrors.ErrInternalServer()
	}

	// Apply config
	hostname := fmt.Sprintf("%s.%s.%s", req.GetServiceName(), req.GetNamespace(), s.config.BaseDomain)
	if err := s.kubeutil.CreateServiceAndDeployment(ctx, req.GetNamespace(), req.GetServiceName(), hostname, req.GetServicePort(), req.GetDockerImage(), req.GetEnv()); err != nil {
		log.Extract(ctx).WithError(err).Error("unable to create service and deployment")

		return nil, grpcerrors.ErrInternalServer()
	}

	portURIString := ""
	if req.GetServicePort() != httpPort {
		portURIString = fmt.Sprintf(":%d", req.GetServicePort())
	}
	uri := strings.Replace(strings.Replace(strings.Replace(strings.Replace(uriScheme,
		"{{service}}", req.GetServiceName(), 1),
		"{{namespace}}", req.GetNamespace(), 1),
		"{{port}}", portURIString, 1),
		"{{gqlUri}}", req.GetGqlUri(), 1)

	// Check whether audience exists
	c := graphql.NewClient(fmt.Sprintf("%s/gql", s.config.AuthorizationURI)) // FIXME: https://lab.weave.nl/nid/nid-core/-/issues/20
	r := graphql.NewRequest(audienceQuery)
	r.Var("filter", map[string]interface{}{"audience": map[string]string{"eq": uri}})
	audienceResp := audienceResponse{}
	if err := c.Run(ctx, r, &audienceResp); err != nil {
		log.Extract(ctx).WithError(err).Error("unable to query audience")

		return nil, grpcerrors.ErrInternalServer()
	}
	audienceExists := len(audienceResp.Audiences) > 0

	// Create audience if not exists
	if !audienceExists {
		r := graphql.NewRequest(createAudienceMutation)
		r.Var("audience", CreateAudienceInput{
			Namespace: req.GetNamespace(),
			Audience:  uri,
		},
		)
		err := c.Run(ctx, r, make(map[string]interface{}))
		if err != nil {
			log.Extract(ctx).WithError(err).Error("unable to create audience mutation")

			return nil, grpcerrors.ErrInternalServer()
		}
	} else {
		log.Infof("audience %s already exists, skipping creation\n", uri)
	}

	// Deploy autopseudo if not yet present
	if err := s.createAutoPseudoIfNotExists(ctx, req.GetNamespace()); err != nil {
		log.Extract(ctx).WithError(err).Error("unable to create auto pseudo")

		return nil, grpcerrors.ErrInternalServer()
	}

	return &proto.DeployServiceResponse{ClusterUri: uri, KubernetesOutput: ""}, nil
}

// DeleteService deletes a service in a given namespace
func (s *DashboardServiceServer) DeleteService(ctx context.Context, req *proto.DeleteServiceRequest) (*emptypb.Empty, error) {
	err := s.validateNamespace(req.GetNamespace(), true)
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument(cases.Title(language.English, cases.NoLower).String(err.Error()))
	}

	err = s.kubeutil.DeleteService(ctx, req.GetNamespace(), req.GetName())
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to delete services")

		return nil, grpcerrors.ErrInternalServer()
	}

	return &emptypb.Empty{}, nil
}

// ListServices lists services for a given namespace
func (s *DashboardServiceServer) ListServices(ctx context.Context, req *proto.ListServiceRequest) (*proto.ServiceList, error) {
	err := s.validateNamespace(req.GetNamespace(), false)
	if err != nil {
		return nil, grpcerrors.ErrInvalidArgument(cases.Title(language.English, cases.NoLower).String(err.Error()))
	}

	svcs, err := s.kubeutil.ListService(ctx, req.GetNamespace())
	if err != nil {
		log.Extract(ctx).WithError(err).Error("unable to list services")

		return nil, grpcerrors.ErrInternalServer()
	}
	resp := &proto.ServiceList{
		Items: make([]*proto.Service, len(svcs.Items)),
	}
	for i := range svcs.Items {
		resp.Items[i] = &proto.Service{
			Name:      svcs.Items[i].GetName(),
			Namespace: svcs.Items[i].GetNamespace(),
			Age: &duration.Duration{
				Seconds: int64(time.Now().UTC().Sub(svcs.Items[i].GetCreationTimestamp().UTC()).Seconds()),
			},
		}
	}

	return resp, nil
}

func (s *DashboardServiceServer) validateNamespace(namespace string, deleteService bool) error {
	if strings.HasPrefix(namespace, "kube") || strings.HasPrefix(namespace, "istio") || namespace == "default" || (deleteService && namespace == s.config.Namespace) {
		return errors.Wrap(ErrInvalidNamespace, "unexpected namespace")
	}

	return nil
}
