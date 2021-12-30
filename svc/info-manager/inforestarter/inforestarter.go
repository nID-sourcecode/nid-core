// Package inforestarter restarts info services.
package inforestarter

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubectl/pkg/cmd/rollout"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"lab.weave.nl/nid/nid-core/pkg/utilities/errors"
)

// InfoRestarter restarts info services
type InfoRestarter interface {
	RestartInfoServices() error
}

// K8sRolloutInfoRestarter restarts info services by triggering a `kubectl rollout` for any deployment with the label
// `info_service_automatic_restart=enabled`
type K8sRolloutInfoRestarter struct {
	factory   cmdutil.Factory
	clientset *kubernetes.Clientset
	namespace string
}

// NewK8sRolloutInfoRestarter creates a new K8sRolloutInfoRestarter
func NewK8sRolloutInfoRestarter(namespace string) (*K8sRolloutInfoRestarter, error) {
	factory := cmdutil.NewFactory(&genericclioptions.ConfigFlags{Namespace: &namespace})
	clientset, err := factory.KubernetesClientSet()
	if err != nil {
		return nil, errors.Wrap(err, "creating k8s clientset")
	}

	return &K8sRolloutInfoRestarter{
		factory:   factory,
		clientset: clientset,
		namespace: namespace,
	}, nil
}

// RestartInfoServices restarts info services by triggering a `kubectl rollout` for any deployment with the label
// `info_service_automatic_restart=enabled`
func (r *K8sRolloutInfoRestarter) RestartInfoServices() error {
	ctx := context.Background()

	deployments, err := r.clientset.AppsV1().Deployments(r.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "info_service_automatic_restart=enabled",
	})
	if err != nil {
		return errors.Wrap(err, "getting deployments")
	}

	for i := range deployments.Items {
		name := deployments.Items[i].Name
		streams, _, _, _ := genericclioptions.NewTestIOStreams()
		restarter := rollout.NewCmdRolloutRestart(r.factory, streams)
		restarter.SetArgs([]string{"deployments/" + name})

		err := restarter.Execute()
		if err != nil {
			return errors.Wrapf(err, "restarting deployment %s", name)
		}
	}

	return nil
}
