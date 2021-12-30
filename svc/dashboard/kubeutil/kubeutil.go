// Package kubeutil provides kubernetes utility functionality
package kubeutil

import (
	"context"
	"fmt"

	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	istio "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

const (
	defaultVirtualServicePort = 80
)

var errEmptyExternalHostname error = fmt.Errorf("parameter ExternalHostname is an empty string")

// Client hosts the kubernetes and istio clients
type Client struct {
	KubernetesClient kubernetes.Interface
	IstioClient      istio.Interface
}

// NewKubeUtil initialises the kubeutil client
func NewKubeUtil(kubernetesClient kubernetes.Interface, istioClient istio.Interface) Interface {
	return &Client{
		KubernetesClient: kubernetesClient,
		IstioClient:      istioClient,
	}
}

// Interface is the kubeutil interface
type Interface interface {
	CreateServiceAndDeployment(ctx context.Context, namespace, serviceName, externalHostname string, servicePort int32, dockerImage string, env map[string]string) error
	CreateVirtualService(ctx context.Context, externalHostname, serviceName, namespace string, servicePort int32) error
	CreateService(ctx context.Context, externalHostname, serviceName, namespace string, servicePort int32) error
	CreateDeployment(ctx context.Context, externalHostname, serviceName, namespace, dockerImage string, servicePort int32, env map[string]string) error
	ListNamespace(ctx context.Context) ([]string, error)
	DeleteService(ctx context.Context, namespace, name string) error
	ListService(ctx context.Context, namespace string) (*corev1.ServiceList, error)
}

// CreateServiceAndDeployment creates service, virtual service and deployment
func (c *Client) CreateServiceAndDeployment(ctx context.Context, namespace, serviceName, externalHostname string,
	servicePort int32, dockerImage string, env map[string]string) error {
	err := c.CreateService(ctx, externalHostname, serviceName, namespace, servicePort)
	if err != nil {
		return err
	}

	err = c.CreateVirtualService(ctx, externalHostname, serviceName, namespace, servicePort)
	if err != nil {
		return err
	}

	err = c.CreateDeployment(ctx, externalHostname, serviceName, namespace, dockerImage, servicePort, env)
	if err != nil {
		return err
	}

	return nil
}

// CreateVirtualService creates a virtual service
func (c *Client) CreateVirtualService(ctx context.Context, externalHostname, serviceName, namespace string, servicePort int32) error {
	if externalHostname == "" {
		return errEmptyExternalHostname
	}
	virtualService := v1beta1.VirtualService{
		TypeMeta: metav1.TypeMeta{
			Kind: "VirtualService",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: networkingv1beta1.VirtualService{
			Hosts:    []string{externalHostname},
			Gateways: []string{"external-access-gateway"},
			Http: []*networkingv1beta1.HTTPRoute{
				{
					Route: []*networkingv1beta1.HTTPRouteDestination{
						{
							Destination: &networkingv1beta1.Destination{
								Host: fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace),
								Port: &networkingv1beta1.PortSelector{
									Number: defaultVirtualServicePort,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := c.IstioClient.NetworkingV1beta1().VirtualServices(namespace).Create(ctx, &virtualService, metav1.CreateOptions{})
	if err != nil {
		return errors.Wrap(err, "unable ot create virtual service")
	}

	return nil
}

// CreateService creates a service for given name port and namespace
func (c *Client) CreateService(ctx context.Context, externalHostname, serviceName, namespace string, servicePort int32) error {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: map[string]string{"app": serviceName, "origin": "dashboard"},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       servicePort,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: servicePort},
				},
			},
			Selector: map[string]string{"app": serviceName},
		},
	}
	if _, err := c.KubernetesClient.CoreV1().Services(namespace).Create(ctx, service, metav1.CreateOptions{}); err != nil {
		return errors.Wrap(err, "unable to create service")
	}

	return nil
}

// CreateDeployment creates a deployment
func (c *Client) CreateDeployment(ctx context.Context, externalHostname, serviceName, namespace, dockerImage string, servicePort int32, env map[string]string) error {
	envVars := make([]corev1.EnvVar, 0)
	for k, v := range env {
		envVars = append(envVars, corev1.EnvVar{
			Name:  k,
			Value: v,
		})
	}

	var replicas int32 = 1
	deployment := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: map[string]string{"origin": "dashboard"},
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": serviceName},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": serviceName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image:           dockerImage,
							ImagePullPolicy: "Always",
							Name:            serviceName,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: servicePort,
								},
							},
							Env: envVars,
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{
						{
							Name: "registry-weave",
						},
					},
				},
			},
		},
	}
	if _, err := c.KubernetesClient.AppsV1().Deployments(namespace).Create(ctx, deployment, metav1.CreateOptions{}); err != nil {
		return errors.Wrap(err, "unable to create deployment")
	}

	return nil
}

// ListNamespace list namespaces for given client set
func (c *Client) ListNamespace(ctx context.Context) ([]string, error) {
	namespaces, err := c.KubernetesClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to list namespaces")
	}
	namespaceNames := make([]string, len(namespaces.Items))
	for i := range namespaces.Items {
		namespaceNames[i] = namespaces.Items[i].GetName()
	}

	return namespaceNames, nil
}

// DeleteService deletes a service with given name in given namespace
func (c *Client) DeleteService(ctx context.Context, namespace, name string) error {
	err := c.KubernetesClient.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to get delete service")
	}

	err = c.KubernetesClient.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to get delete deployment")
	}

	vs, err := c.IstioClient.NetworkingV1beta1().VirtualServices(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to get virtual service")
	}
	if vs != nil {
		err = c.IstioClient.NetworkingV1beta1().VirtualServices(namespace).Delete(ctx, name, metav1.DeleteOptions{})
		if err != nil {
			return errors.Wrap(err, "unable to delete virtual service")
		}
	}

	err = c.IstioClient.SecurityV1beta1().RequestAuthentications(namespace).Delete(ctx, name+"-jwtauth", metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to delete request authentication")
	}

	err = c.IstioClient.SecurityV1beta1().AuthorizationPolicies(namespace).Delete(ctx, name+"-jwtcheck", metav1.DeleteOptions{})
	if err != nil {
		return errors.Wrap(err, "unable to delete authorization policy")
	}

	return nil
}

// ListService list services in given namespace
func (c *Client) ListService(ctx context.Context, namespace string) (*corev1.ServiceList, error) {
	services, err := c.KubernetesClient.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to list kubernetes services")
	}

	return services, nil
}
