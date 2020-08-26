package reconciler

import (
	"context"
	"fmt"

	v1 "github.com/cirocosta/scm-controller/api/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Client

// Client represents a client that's able to retrieve objects from
// Kubernetes.
//
type Client interface {
	ServiceAccount(ctx context.Context, name, namespace string) (*corev1.ServiceAccount, bool, error)
	Secret(ctx context.Context, name, namespace string) (*corev1.Secret, bool, error)
	CommitStatus(ctx context.Context, name, namespace string) (*v1.CommitStatus, bool, error)
}

func NewClient(client client.Client) *getter {
	return &getter{
		client: client,
	}
}

type getter struct {
	client client.Client
}

func (g *getter) ServiceAccount(ctx context.Context, name, namespace string) (*corev1.ServiceAccount, bool, error) {
	obj := new(corev1.ServiceAccount)
	found, err := g.get(ctx, name, namespace, obj)
	return obj, found, err
}

func (g *getter) Secret(ctx context.Context, name, namespace string) (*corev1.Secret, bool, error) {
	obj := new(corev1.Secret)
	found, err := g.get(ctx, name, namespace, obj)
	return obj, found, err
}

func (g *getter) CommitStatus(ctx context.Context, name, namespace string) (*v1.CommitStatus, bool, error) {
	obj := new(v1.CommitStatus)
	found, err := g.get(ctx, name, namespace, obj)
	return obj, found, err
}

func (g *getter) get(ctx context.Context, name, namespace string, object runtime.Object) (bool, error) {
	err := g.client.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, object)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}

		return false, fmt.Errorf("get %s: %w", object.GetObjectKind().GroupVersionKind(), err)
	}

	return true, nil
}
