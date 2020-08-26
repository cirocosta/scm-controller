package reconciler_test

import (
	"context"
	"fmt"

	"github.com/cirocosta/scm-controller/reconciler"
	"github.com/cirocosta/scm-controller/reconciler/reconcilerfakes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {
	Context("SecretFromServiceAccount", func() {
		var (
			ctx    context.Context
			client *reconcilerfakes.FakeClient

			annotationKey  string
			serviceAccount string
			namespace      string

			secret *corev1.Secret
			err    error
		)

		BeforeEach(func() {
			ctx = context.Background()
			namespace = "namespace"
			serviceAccount = "serviceaccount"

			client = new(reconcilerfakes.FakeClient)
		})

		JustBeforeEach(func() {
			secret, err = reconciler.
				NewSecretFetcher(client).
				SecretFromServiceAccount(
					ctx, annotationKey, serviceAccount, namespace,
				)
		})

		Context("failing to retrieve service account", func() {
			var expectedErr error

			BeforeEach(func() {
				expectedErr = fmt.Errorf("failed")
				client.ServiceAccountReturns(nil, false, expectedErr)
			})

			It("fails", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(expectedErr))
			})
		})

		Context("not finding the serviceaccount", func() {
			BeforeEach(func() {
				client.ServiceAccountReturns(nil, false, nil)
			})

			It("fails", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})

		Context("succeeding getting the service account", func() {
			var serviceAccountObject *corev1.ServiceAccount

			BeforeEach(func() {
				serviceAccountObject = &corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name:      serviceAccount,
						Namespace: serviceAccount,
					},
					Secrets: []corev1.ObjectReference{},
				}

				client.ServiceAccountReturns(serviceAccountObject, true, nil)
			})

			Context("having no secrets linked to it", func() {
				It("fails", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("no secret found"))
				})
			})

			Context("having secrets linked to it", func() {
				BeforeEach(func() {
					serviceAccountObject.Secrets = []corev1.ObjectReference{
						{Name: "secret1"},
						{Name: "secret2"},
					}
				})

				Context("failing to retrieve secret", func() {
					var expectedErr error

					BeforeEach(func() {
						expectedErr = fmt.Errorf("failed")
						client.SecretReturns(nil, false, expectedErr)
					})

					It("fails", func() {
						Expect(err).To(HaveOccurred())
						Expect(err).To(MatchError(expectedErr))
					})
				})

				Context("secrets being retrievable", func() {
					secrets := []*corev1.Secret{
						{ObjectMeta: metav1.ObjectMeta{
							Name: "secret1",
						}},

						{ObjectMeta: metav1.ObjectMeta{
							Name: "secret2",
						}},
					}

					BeforeEach(func() {
						client.SecretReturnsOnCall(0, secrets[0], true, nil)
						client.SecretReturnsOnCall(1, secrets[1], true, nil)
					})

					Context("none being annotated", func() {
						It("fails", func() {
							Expect(err).To(HaveOccurred())
							Expect(err).To(MatchError("no secret found"))
						})
					})

					Context("multiple secrets matching annotation", func() {
						BeforeEach(func() {
							secrets[0].Annotations = map[string]string{
								annotationKey: "foo",
							}
							secrets[1].Annotations = map[string]string{
								annotationKey: "foo",
							}
						})

						It("works", func() {
							Expect(err).ToNot(HaveOccurred())
						})

						It("picks first", func() {
							Expect(secret.Name).To(Equal(secrets[0].Name))
						})
					})
				})
			})
		})
	})
})
