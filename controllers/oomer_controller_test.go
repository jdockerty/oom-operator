package controllers

import (
	"context"
	"time"

	oomv1alpha1 "github.com/jdockerty/oom-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	oomerApiVersion = "jdocklabs.co.uk/v1alpha1"
	oomerKind       = "Oomer"
)

var _ = Describe("Oomer Operator", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		operatorName   = "test-oomer"
		oomerNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating the object", func() {
		It("Should create an underlying deployment object", func() {

			var replicas int32 = 1

			ctx := context.Background()
			oom := &oomv1alpha1.Oomer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: oomerApiVersion,
					Kind:       oomerKind,
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      operatorName,
					Namespace: oomerNamespace,
				},
				Spec: oomv1alpha1.OomerSpec{
					Replicas: &replicas,
				},
			}

			Expect(k8sClient.Create(ctx, oom)).Should(Succeed())

			lookupOomer := types.NamespacedName{Name: operatorName, Namespace: oomerNamespace}
			createdOomer := &oomv1alpha1.Oomer{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupOomer, createdOomer)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			By("checking the underlying deployment exists")
			d := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupOomer, d)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(d.ObjectMeta.Name).Should(Equal(oom.ObjectMeta.Name))
			Expect(d.Spec.Replicas).Should(Equal(oom.Spec.Replicas))
		})
	})
})
