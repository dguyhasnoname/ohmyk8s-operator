package controller

import (
	"context"
	"io/ioutil"
	"path"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	namespaceconfigv1 "github.com/dguyhasnoname/ohmyk8s-operator/api/v1"
)

func ReadNamespaceConfigFromFile(filename string) (*namespaceconfigv1.Namespaceconfig, error) {
	testfile := path.Join("..", "..", "config", "samples", filename)
	content, err := ioutil.ReadFile(filepath.Clean(testfile))
	if err != nil {
		return nil, err
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(content, nil, nil)
	if err != nil {
		return nil, err
	}
	return obj.(*namespaceconfigv1.Namespaceconfig), nil
}

var _ = Describe("Namespace controller", func() {

	Context("NamespaceConfig with given size", func() {
		It("Should create namespace successfully", func() {
			By("Creating a new Namespace")
			ctx := context.Background()
			nc, err1 := ReadNamespaceConfigFromFile("namespaceconfig_v1_namespaceconfig.yaml")
			Expect(err1).ToNot(HaveOccurred())
			Expect(k8sClient.Create(ctx, nc)).Should(Succeed())
			// Verify Namespace
			ns := &v1.Namespace{}
			nsname := "apr-dev"
			Eventually(func() bool {
				_ = k8sClient.Get(ctx, types.NamespacedName{Name: nsname, Namespace: ""}, ns)
				return ns.ObjectMeta.UID != ""
			}, time.Second*30, time.Second*3).Should(BeTrue())
			Expect(ns.ObjectMeta.Labels).To(HaveKeyWithValue("env", "dev"))
			Expect(ns.ObjectMeta.Labels).To(HaveKeyWithValue("owner", "mukund"))
			Expect(ns.Status.Phase).To(Equal(v1.NamespaceActive))
			//Delete namespace
			By("Deleting an existing namespacedef")
			Eventually(func() bool {
				errdel := k8sClient.Delete(ctx, nc)
				return errdel == nil
			}, time.Second*30, time.Second*3).Should(BeTrue())
		})
	})
})
