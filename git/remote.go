package git

import (
	"context"
	"fmt"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
)

type remote struct {
	client     *scm.Client
	repository string
	token      string

	serverUrl string
	driver    string
}

func WithServerUrl(v string) opt {
	return func(r *remote) {
		r.serverUrl = v
	}
}

func WithDriver(v string) opt {
	return func(r *remote) {
		r.driver = v
	}
}

type opt func(*remote)

func New(repository, token string, opts ...opt) (*remote, error) {
	r := &remote{
		token:      token,
		repository: repository,
	}

	for _, opt := range opts {
		opt(r)
	}

	client, err := factory.NewClient(r.driver, r.serverUrl, token)
	if err != nil {
		return nil, fmt.Errorf("new client from env: %w", err)
	}

	r.client = client

	return r, nil
}

func (r *remote) ListStatuses(
	ctx context.Context,
	sha string,
) error {
	statuses, resp, err := r.client.Repositories.ListStatus(ctx, r.repository, sha, scm.ListOptions{})
	if err != nil {
		return fmt.Errorf("list status: %w", err)
	}

	_, _ = statuses, resp

	return nil
}

func (r *remote) UpdateStatus(
	ctx context.Context,
	revision, label, desc, target string,
) error {
	status, resp, err := r.client.Repositories.CreateStatus(
		ctx, r.repository, revision,
		&scm.StatusInput{
			State:  scm.StateSuccess,
			Label:  label,
			Desc:   desc,
			Target: target,
		},
	)
	if err != nil {
		return fmt.Errorf("create status: %w", err)
	}

	_, _ = status, resp

	return nil
}
