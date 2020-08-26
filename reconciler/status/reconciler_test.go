package status_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/cirocosta/scm-controller/api/v1"
	"github.com/cirocosta/scm-controller/reconciler/reconcilerfakes"
	"github.com/cirocosta/scm-controller/reconciler/status"
	"github.com/cirocosta/scm-controller/reconciler/status/statusfakes"
	ctrl "sigs.k8s.io/controller-runtime"
)

var _ = Describe("StatusReconciler", func() {
	var (
		reconciler status.Reconciler
		result     ctrl.Result
		err        error

		secretFetcher *statusfakes.FakeSecretFetcher
		client        *reconcilerfakes.FakeClient

		emptyResult = ctrl.Result{}
	)

	BeforeEach(func() {
		client = new(reconcilerfakes.FakeClient)
		secretFetcher = new(statusfakes.FakeSecretFetcher)

		reconciler = status.Reconciler{
			Client:        client,
			SecretFetcher: secretFetcher,
		}
	})

	JustBeforeEach(func() {
		result, err = reconciler.Reconcile(ctrl.Request{})
	})

	Context("erroring retrieving commit status", func() {
		var expectedErr error

		BeforeEach(func() {
			expectedErr = fmt.Errorf("failed")
			client.CommitStatusReturns(nil, false, expectedErr)
		})

		It("fails", func() {
			Expect(err).To(MatchError(expectedErr))
		})
	})

	Context("not finding commitstatus", func() {
		BeforeEach(func() {
			client.CommitStatusReturns(nil, false, nil)
		})

		It("doesn't fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("doesn't requeue", func() {
			Expect(result).To(Equal(emptyResult))
		})
	})

	Context("finding commitstatus", func() {
		BeforeEach(func() {
			client.CommitStatusReturns(&v1.CommitStatus{}, true, nil)
		})

		Context("failing to find secret for serviceaccount", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = fmt.Errorf("failed")
				secretFetcher.SecretFromServiceAccountReturns(nil, expectedErr)
			})

			It("fails", func() {
				Expect(err).To(MatchError(expectedErr))
			})
		})
	})
})
