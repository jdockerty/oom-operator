package controllers

import (
	"context"
	"time"

	oomv1alpha1 "github.com/jdockerty/oom-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Oomer Operator", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		operatorName    = "test-oomer"
		oomerApiVersion = "jdocklabs.co.uk/v1alpha1"
		oomerKind       = "Oomer"
		oomerNamespace  = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

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

	Context("When creating the object", func() {
		It("Should create an underlying deployment object", func() {

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

		It("Should update the status to reflect the observed replicas", func() {

			lookupOomer := types.NamespacedName{Name: operatorName, Namespace: oomerNamespace}
			createdOomer := &oomv1alpha1.Oomer{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupOomer, createdOomer)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			d := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupOomer, d)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(d.Spec.Replicas).Should(Equal(createdOomer.Spec.Replicas))
			Expect(createdOomer.Status.ObservedReplicas).Should(Equal(createdOomer.Spec.Replicas))
		})
	})

	Context("When deleting the object", func() {
		It("should delete the underlying deployment", func() {

			lookupOomer := types.NamespacedName{Name: operatorName, Namespace: oomerNamespace}
			createdOomer := &oomv1alpha1.Oomer{}

			By("ensuring finalizers exist")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, lookupOomer, createdOomer)
				if err != nil {
					return false
				}
				return len(createdOomer.ObjectMeta.Finalizers) > 0
			}, timeout, interval).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, oom)).Should(Succeed())

			Consistently(func() bool {
				err := k8sClient.Get(ctx, lookupOomer, &appsv1.Deployment{})
				if err != nil {
					if apierrors.IsNotFound(err) {
						return true
					}
					return false
				}
				return false
			})

		})
	})
})
