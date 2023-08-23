package v1alpha1

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("NamespaceLabel Webhook", func() {
	// Define a NamespaceLabel object
	var namespaceLabel1 *NamespaceLabel

	BeforeEach(func() {
		// Reset the NamespaceLabel object before each test
		namespaceLabel1 = &NamespaceLabel{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "dana.io/v1alpha1",
				Kind:       "NamespaceLabel",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "namespacelabel-webhook-test1",
				Namespace: "default",
			},
			Spec: NamespaceLabelSpec{
				Labels: map[string]string{
					"name":          "namespacelabel1",
					"examplelabel1": "one",
				},
			},
		}
	})

	Context("when validating NamespaceLabel creation", func() {
		It("should prevent creation if a label has a disallowed prefix", func() {
			// Set a disallowed prefix
			namespaceLabel1.Spec.Labels = map[string]string{
				"kubernetes.io/some-label": "value",
			}

			// Try to create the NamespaceLabel
			err := k8sClient.Create(ctx, namespaceLabel1)
			Expect(err).To(HaveOccurred())
			// You can add more specific assertions to check the error message
		})

		It("should allow creation if all labels have allowed prefixes", func() {
			// Set an allowed prefix
			namespaceLabel1.Spec.Labels = map[string]string{
				"some-allowed-prefix/some-label": "value",
			}

			// Try to create the NamespaceLabel
			err := k8sClient.Create(ctx, namespaceLabel1)
			Expect(err).NotTo(HaveOccurred())

			// Cleanup the created NamespaceLabel
			Expect(k8sClient.Delete(ctx, namespaceLabel1)).Should(Succeed())
		})
	})

})
