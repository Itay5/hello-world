package utils_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"dana.io/hello-world/internal/controller/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLabels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Labels Suite")
}

var _ = Describe("Labels", func() {
	It("should add and remove labels appropriately", func() {
		namespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"existingLabel": "value",
				},
			},
		}

		labelsToAdd := map[string]string{
			"labelToAdd": "newValue",
		}

		labelsToRemove := map[string]struct{}{
			"existingLabel": {},
		}

		utils.UpdateNamespaceLabels(namespace, labelsToAdd, labelsToRemove)

		Expect(namespace.ObjectMeta.Labels).To(HaveKeyWithValue("labelToAdd", "newValue"))
		Expect(namespace.ObjectMeta.Labels).NotTo(HaveKey("existingLabel"))
	})
})
