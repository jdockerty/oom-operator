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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	oomv1alpha1 "github.com/jdockerty/oom-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// OomerReconciler reconciles a Oomer object
type OomerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=jdocklabs.co.uk,resources=oomers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=jdocklabs.co.uk,resources=oomers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=jdocklabs.co.uk,resources=oomers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Oomer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
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

	if oomer.Spec.Replicas == 0 {
		log.Info("0 replicas, nothing to do")
		return ctrl.Result{}, nil
	}

	log.Info("reconciling oomer", "replicas", oomer.Spec.Replicas)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OomerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&oomv1alpha1.Oomer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
