/*
Copyright 2021.

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
	"os"
	"strings"

	"github.com/go-logr/logr"
	"github.com/imharshita/image-controller/pkg/images"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	// KubeNs Namespace to exclude in Reconiler
	KubeNs = "kube-system"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// DaemonSetReconciler reconciles a DaemonSet object
type DaemonSetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func isImagePresent(image string) bool {
	Registry := os.Getenv("REGISTRY")
	if len(Registry) == 0 {
		return false
	} else if !strings.HasPrefix(image, Registry) {
		return true
	}
	return false
}

func isDaemonSetReady(daemonsets *appsv1.DaemonSet) bool {
	status := daemonsets.Status
	desired := status.DesiredNumberScheduled
	ready := status.NumberReady
	if desired > 0 && ready > 0 && desired == ready {
		return true
	}
	return false
}

// Reconcile recociles DaemonSet to uodate the image
func (r *DaemonSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqNamespace := req.NamespacedName.Namespace
	if reqNamespace != KubeNs {
		_ = r.Log.WithValues("daemonset", req.NamespacedName)
		daemonsets := &appsv1.DaemonSet{}
		err := r.Get(context.TODO(), req.NamespacedName, daemonsets)
		if err != nil {
			return reconcile.Result{}, err
		}
		if isDaemonSetReady(daemonsets) {
			containers := daemonsets.Spec.Template.Spec.Containers
			for i, c := range containers {
				if isImagePresent(c.Image) {
					r.Log.WithValues("Retagging image", c.Image, "of daemonset:", daemonsets.Name)
					img, err := images.Process(c.Image)
					if err != nil {
						return reconcile.Result{}, err
					}
					r.Log.WithValues("Updating image", c.Image, "of daemonset:", daemonsets.Name)
					daemonsets.Spec.Template.Spec.Containers[i].Image = img
					err = r.Update(context.TODO(), daemonsets)
					if err != nil {
						return reconcile.Result{}, err
					}
					r.Log.WithValues("Updated image:", c.Image, " -> ", img)
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

func isDeploymentReady(deployments *appsv1.Deployment) bool {
	status := deployments.Status
	desired := status.Replicas
	ready := status.ReadyReplicas
	if desired > 0 && ready > 0 && desired == ready {
		return true
	}
	return false
}

// Reconcile reconciles Deployment to update the image
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqNamespace := req.NamespacedName.Namespace
	if reqNamespace != KubeNs {
		_ = r.Log.WithValues("deployment", req.NamespacedName)
		deployments := &appsv1.Deployment{}
		err := r.Get(context.TODO(), req.NamespacedName, deployments)
		if err != nil {
			return reconcile.Result{}, err
		}
		if isDeploymentReady(deployments) {
			containers := deployments.Spec.Template.Spec.Containers
			for i, c := range containers {
				if isImagePresent(c.Image) {
					r.Log.WithValues("Retagging image", c.Image, "of deployment:", deployments.Name)
					img, err := images.Process(c.Image)
					if err != nil {
						return reconcile.Result{}, err
					}
					// Update image
					r.Log.WithValues("Updating image", c.Image, "of deployment:", deployments.Name)
					deployments.Spec.Template.Spec.Containers[i].Image = img
					err = r.Update(context.TODO(), deployments)
					if err != nil {
						return reconcile.Result{}, err
					}
					r.Log.WithValues("Updated image:", c.Image, " -> ", img)
				}
			}
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Owns(&corev1.Pod{}).
		Watches(&source.Kind{Type: &appsv1.Deployment{}},
			&handler.EnqueueRequestForObject{}).
		Complete(r)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DaemonSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		Owns(&corev1.Pod{}).
		Watches(&source.Kind{Type: &appsv1.DaemonSet{}},
			&handler.EnqueueRequestForObject{}).
		Complete(r)
}
