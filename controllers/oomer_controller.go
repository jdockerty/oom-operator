/*
Copyright 2023 Jack Dockerty.

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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	oomv1alpha1 "github.com/jdockerty/oom-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const (
	defaultImage           = "jdockerty/oomer:v0.0.1"
	terminationMessagePath = "/tmp/oomed-pod.log"
)

// OomerReconciler reconciles a Oomer object
type OomerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *OomerReconciler) createOrUpdateDeployment(ctx context.Context, req ctrl.Request, o *oomv1alpha1.Oomer, log logr.Logger) error {

	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        req.Name,
			Namespace:   req.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: make(map[string]string),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      make(map[string]string),
					Annotations: make(map[string]string),
				},
				Spec: corev1.PodSpec{
					Containers: make([]corev1.Container, 1),
				},
			},
		},
	}
	if err := r.Get(ctx, req.NamespacedName, d); err != nil {

		// If not found, create the Deployment
		if apierrors.IsNotFound(err) {

			d.Spec.Replicas = o.Spec.Replicas

			if o.Spec.Labels != nil {
				d.Spec.Selector.MatchLabels = o.Spec.Labels
				d.Spec.Template.ObjectMeta.Labels = o.Spec.Labels
			} else {
				d.Spec.Selector.MatchLabels["app"] = "oomer"
				d.Spec.Template.ObjectMeta.Labels["app"] = "oomer"
			}

			if o.Spec.Image != nil {
				d.Spec.Template.Spec.Containers[0].Image = *o.Spec.Image
			} else {
				d.Spec.Template.Spec.Containers[0].Image = defaultImage
			}

			log.Info("set image", "image", d.Spec.Template.Spec.Containers[0].Image)
			d.Spec.Template.Spec.Containers[0].Name = "oomer"
			d.Spec.Template.Spec.Containers[0].TerminationMessagePath = terminationMessagePath

			log.Info("underlying deployment not found, creating...")

			if err := r.Create(ctx, d); err != nil {
				return err
			}

			log.Info("updating oomer observed replicas status", "replicas", o.Spec.Replicas)

			// Update the status of observed replicas to those which are
			// provided in the spec/to the deployment
			o.Status.ObservedReplicas = o.Spec.Replicas
			if err := r.Status().Update(ctx, o); err != nil {
				log.Error(err, "unable to update oomer status observed replicas", "ObservedReplicas", o.Status.ObservedReplicas, "Spec.Replicas", o.Spec.Replicas)
				return err
			}
		}

	}

	return nil
}

//+kubebuilder:rbac:groups=jdocklabs.co.uk,resources=oomers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=jdocklabs.co.uk,resources=oomers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=jdocklabs.co.uk,resources=oomers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OomerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var oomer oomv1alpha1.Oomer
	if err := r.Get(ctx, req.NamespacedName, &oomer); err != nil {
		if apierrors.IsNotFound(err) {
			// we'll ignore not-found errors, since they can't be fixed by an immediate
			// requeue (we'll need to wait for a new notification), and we can get them
			// on deleted requests.
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Oomer")
		return ctrl.Result{}, err
	}

	if *oomer.Spec.Replicas == int32(0) {
		log.Info("0 replicas, no creation")
		return ctrl.Result{}, nil
	}

	log.Info("reconciling oomer", "replicas", oomer.Spec.Replicas)

	if err := r.createOrUpdateDeployment(ctx, req, &oomer, log); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OomerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oomv1alpha1.Oomer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
