package reconciler

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type secretFetcher struct {
	client Client
}

func NewSecretFetcher(
	client Client,
) *secretFetcher {
	return &secretFetcher{
		client,
	}
}

// SecretFromServiceAccount retrieves a secret (with a given annotation) out of
// a set of secrets linked to a service account.
//
func (f *secretFetcher) SecretFromServiceAccount(
	ctx context.Context,
	annotationKey string,
	serviceAccount, namespace string,
) (*corev1.Secret, error) {
	account, found, err := f.client.ServiceAccount(ctx, serviceAccount, namespace)
	if err != nil {
		return nil, fmt.Errorf("serviceaccount: %s/%s: %w", serviceAccount, namespace, err)
	}

	if !found {
		return nil, fmt.Errorf("service account %s/%s not found", serviceAccount, namespace)
	}

	for _, secretRef := range account.Secrets {
		secret, _, err := f.client.Secret(ctx, secretRef.Name, namespace)
		if err != nil {
			return nil, fmt.Errorf("get secret %s/%s: %w", secretRef.Name, namespace, err)
		}

		_, found := secret.GetAnnotations()[annotationKey]
		if !found {
			continue
		}

		return secret, nil
	}

	return nil, fmt.Errorf("no secret found")
}
