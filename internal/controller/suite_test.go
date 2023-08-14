/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	danaiodanaiov1alpha1 "dana.io/hello-world/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

var namespaceLabel1 *danaiodanaiov1alpha1.NamespaceLabel
var namespaceLabel2 *danaiodanaiov1alpha1.NamespaceLabel

const (
	FirstNamespaceLabelName = "test-namespacelabel1"

	NamespaceLabelNamespace = "default"

	SecondNamespaceLabelName = "test-namespacelabel2"

	timeout  = time.Second * 10
	duration = time.Second * 10
	interval = time.Millisecond * 250
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {

	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	useExistingCluster := true
	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		UseExistingCluster: &useExistingCluster,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = danaiodanaiov1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	ctx := context.Background()

	By("By creating a second NamespaceLabel")
	namespaceLabel1 = &danaiodanaiov1alpha1.NamespaceLabel{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dana.io/v1alpha1",
			Kind:       "NamespaceLabel",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      FirstNamespaceLabelName,
			Namespace: NamespaceLabelNamespace,
		},
		Spec: danaiodanaiov1alpha1.NamespaceLabelSpec{
			Labels: map[string]string{
				"name":          "namespacelabel1",
				"examplelabel1": "one",
			},
		},
	}

	Expect(k8sClient.Create(ctx, namespaceLabel1)).Should(Succeed())

	By("By creating a second NamespaceLabel")
	namespaceLabel2 = &danaiodanaiov1alpha1.NamespaceLabel{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dana.io/v1alpha1",
			Kind:       "NamespaceLabel",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecondNamespaceLabelName,
			Namespace: NamespaceLabelNamespace,
		},
		Spec: danaiodanaiov1alpha1.NamespaceLabelSpec{
			Labels: map[string]string{
				"name":          "namespacelabel2",
				"examplelabel2": "one",
			},
		},
	}

	Expect(k8sClient.Create(ctx, namespaceLabel2)).Should(Succeed())

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())

	ctx := context.Background()

	// Check if namespaceLabel1 exists and delete it if it does
	if err := k8sClient.Get(ctx, client.ObjectKey{Name: namespaceLabel1.Name, Namespace: namespaceLabel1.Namespace}, &danaiodanaiov1alpha1.NamespaceLabel{}); err == nil {
		Expect(k8sClient.Delete(ctx, namespaceLabel1)).Should(Succeed())
	}

	// Check if namespaceLabel2 exists and delete it if it does
	if err := k8sClient.Get(ctx, client.ObjectKey{Name: namespaceLabel2.Name, Namespace: namespaceLabel2.Namespace}, &danaiodanaiov1alpha1.NamespaceLabel{}); err == nil {
		Expect(k8sClient.Delete(ctx, namespaceLabel2)).Should(Succeed())
	}
})
