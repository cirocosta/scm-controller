package status

import (
	"context"
	"fmt"

	"github.com/cirocosta/scm-controller/git"
	"github.com/cirocosta/scm-controller/reconciler"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Reconciler struct {
	Client        reconciler.Client
	SecretFetcher SecretFetcher
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . SecretFetcher

type SecretFetcher interface {
	SecretFromServiceAccount(
		ctx context.Context,
		annotationKey string,
		serviceAccount, namespace string,
	) (*corev1.Secret, error)
}

const (
	annotationKey = "experimental.kontinue.io/git"
)

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	commitStatus, found, err := r.Client.CommitStatus(ctx, req.Name, req.Namespace)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("get %s: %w", req.NamespacedName, err)
	}

	if !found {
		return ctrl.Result{
			Requeue: false,
		}, nil
	}

	// even before starting - check my own status and see if there's
	// anything else left to be done.

	secret, err := r.SecretFetcher.SecretFromServiceAccount(
		ctx, annotationKey,
		commitStatus.Spec.ServiceAccountName, req.Namespace,
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("secret from svc account: %w", err)
	}

	token, found := secret.Data["password"]
	if !found {
		return ctrl.Result{}, fmt.Errorf("secret '%s' has no password entry")
	}

	remote, err := git.New(
		commitStatus.Spec.Repository,
		string(token),
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("git new: %w", err)
	}

	err = remote.UpdateStatus(
		ctx,
		commitStatus.Spec.Revision,
		commitStatus.Spec.Label,
		commitStatus.Spec.Description,
		commitStatus.Spec.Target,
	)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("")
	}

	// grab the username and password
	// initialize a remote
	// create the status

	return ctrl.Result{}, nil
}
