/*
Copyright 2024.

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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	namespaceconfigv1 "github.com/dguyhasnoname/ohmyk8s-operator/api/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	k8sClient   client.Client
	testEnv     *envtest.Environment
	testCluster string
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
	testCluster := os.ExpandEnv("${TEST_CLUSTER}")
	if testCluster != "" {
		fmt.Println("Running test on", testCluster, " cluster.")
		var err error
		cfg := ctrl.GetConfigOrDie()
		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
		Expect(k8sClient).NotTo(BeNil())

		err = namespaceconfigv1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())
		err = admissionv1beta1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())
		err = apiextensions.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())

	} else {
		By("bootstrapping test environment")
		testEnv = &envtest.Environment{
			CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
			ErrorIfCRDPathMissing: true,

			// The BinaryAssetsDirectory is only required if you want to run the tests directly
			// without call the makefile target test. If not informed it will look for the
			// default path defined in controller-runtime which is /usr/local/kubebuilder/.
			// Note that you must have the required binaries setup under the bin directory to perform
			// the tests directly. When we run make test it will be setup and used automatically.
			BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
				fmt.Sprintf("1.28.3-%s-%s", runtime.GOOS, runtime.GOARCH)),
		}

		var err error
		// cfg is defined in this file globally.
		cfg, err := testEnv.Start()
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())

		err = namespaceconfigv1.AddToScheme(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())

		//+kubebuilder:scaffold:scheme

		k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())
		Expect(k8sClient).NotTo(BeNil())
	}

})

var _ = AfterSuite(func() {
	if os.ExpandEnv("${TEST_CLUSTER}") == "" {
		By("tearing down the test environment")
		err := testEnv.Stop()
		Expect(err).NotTo(HaveOccurred())
	}
})
