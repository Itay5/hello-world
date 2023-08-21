package controller_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("NamespacelabelController", func() {
	ctx := context.Background()
	ns := &corev1.Namespace{}

	Context("When working with NamespaceLabel", func() {
		It("Should handle NamespaceLabel 1 correctly", func() {

			// check namespacelabel 1 reconcile
			By("Waiting for namespacelabel 1 to be reconcile")
			Eventually(func() bool {
				// Get the latest state of the namespace
				err := k8sClient.Get(ctx, client.ObjectKey{Name: NamespaceLabelNamespace}, ns)
				if err != nil {
					return false
				}

				// Check if the namespace labels match the NamespaceLabel spec labels
				for key, value := range namespaceLabel1.Spec.Labels {
					if nsValue, exists := ns.Labels[key]; !exists || nsValue != value {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue(), "Namespace label 1 should be set correctly")
		})

		// check namespacelabel 2 reconcile
		It("Should handle NamespaceLabel 2 correctly", func() {
			// check namespacelabel 2 reconcile
			By("Waiting for namespacelabel 2 to reconcile")
			Eventually(func() bool {

				// Get the latest state of the namespace
				err := k8sClient.Get(ctx, client.ObjectKey{Name: NamespaceLabelNamespace}, ns)
				if err != nil {
					return false
				}

				// Check if the namespace labels match the NamespaceLabel spec labels
				// also considering the first namespacelabel priority
				for key, value := range namespaceLabel2.Spec.Labels {
					_, firstsNsExists := namespaceLabel1.Spec.Labels[key]
					if nsValue, exists := ns.Labels[key]; !exists || nsValue == value && firstsNsExists {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue(), "Namespace label 2 should be set correctly")

		})

		It("Should edit NamespaceLabel 1 correctly", func() {
			newLabels := map[string]string{
				"newkey": "newvalue",
			}

			// Edit namespacelabel 1
			By("Editing the NamespaceLabel 1")
			Eventually(func() error {
				// Fetch the latest version of the object
				err := k8sClient.Get(ctx, client.ObjectKey{Name: namespaceLabel1.Name, Namespace: namespaceLabel1.Namespace}, namespaceLabel1)
				if err != nil {
					return err
				}

				// Apply the changes
				namespaceLabel1.Spec.Labels = newLabels

				// Attempt to update
				return k8sClient.Update(ctx, namespaceLabel1)
			}, timeout, interval).Should(Succeed(), "NamespaceLabel 1 should be updated")

			// Check namespacelabel 1 reconcile after editing
			By("Waiting for namespacelabel 1 to reconcile after editing")
			Eventually(func() bool {
				// Re-fetch the namespace to check the current labels
				err := k8sClient.Get(ctx, client.ObjectKey{Name: NamespaceLabelNamespace}, ns)
				if err != nil {
					return false
				}

				// Check if the namespace labels match the edited NamespaceLabel spec labels
				for key, value := range newLabels {
					if nsValue, exists := ns.Labels[key]; !exists || nsValue != value {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue(), "Namespace label 1 should be updated correctly after editing")
		})

		It("Should enforce original labels after manual deletion", func() {
			// Manually delete a label from the namespace
			By("Manually deleting a label from the namespace")
			deletedKey := "newkey"
			delete(ns.Labels, deletedKey)
			Expect(k8sClient.Update(ctx, ns)).Should(Succeed())

			// Check that the controller enforces the original labels from namespacelabel 1
			By("Waiting for the controller to enforce the original labels from namespacelabel 1")
			Eventually(func() bool {
				// Re-fetch the namespace to check the current labels
				err := k8sClient.Get(ctx, client.ObjectKey{Name: NamespaceLabelNamespace}, ns)
				if err != nil {
					return false
				}

				// Check if the namespace labels match the NamespaceLabel spec labels
				for key, value := range namespaceLabel1.Spec.Labels {
					if nsValue, exists := ns.Labels[key]; !exists || nsValue != value {
						return false
					}
				}

				// Check that the deleted label has been re-added
				if _, exists := ns.Labels[deletedKey]; !exists {
					return false
				}

				return true
			}, timeout, interval).Should(BeTrue(), "Controller should enforce the original labels from namespacelabel 1, even after deletion")
		})

		It("Should delete NamespaceLabel 1 correctly", func() {
			// delete namespacelabel 1
			By("By deleting the NamespaceLabel 1")
			Expect(k8sClient.Delete(ctx, namespaceLabel1)).Should(Succeed())

			// check namespace labels after the deletion of namespacelabel 1
			By("Waiting for namespace labels to be delete")
			Eventually(func() bool {
				// Re-fetch the namespace to check the current labels
				err := k8sClient.Get(ctx, client.ObjectKey{Name: NamespaceLabelNamespace}, ns)
				if err != nil {
					return false
				}

				// Check if the namespace labels were removed
				for key, _ := range namespaceLabel1.Spec.Labels {
					secondNsValue, secondNsExists := namespaceLabel2.Spec.Labels[key]
					if value, exists := ns.Labels[key]; exists && !secondNsExists || secondNsExists && secondNsValue != value {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue(), "Namespace label 1 should be deleted correctly")
		})

		It("Should delete NamespaceLabel 2 correctly", func() {
			// delete namespacelabel 2
			By("By deleting the NamespaceLabel 2")
			Expect(k8sClient.Delete(ctx, namespaceLabel2)).Should(Succeed())

			// check namespace labels after the deletion of namespacelabel 2
			By("Waiting for namespace labels to be delete")
			Eventually(func() bool {
				// Re-fetch the namespace to check the current labels
				err := k8sClient.Get(ctx, client.ObjectKey{Name: NamespaceLabelNamespace}, ns)
				if err != nil {
					return false
				}

				// Check if the namespace labels were removed
				for key, _ := range namespaceLabel2.Spec.Labels {
					if _, exists := ns.Labels[key]; exists {
						return false
					}
				}
				return true
			}, timeout, interval).Should(BeTrue(), "Namespace label 2 should be deleted correctly")
		})

	})
})
