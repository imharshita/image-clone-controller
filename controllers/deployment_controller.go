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
	"fmt"
	"github.com/go-logr/logr"
<<<<<<< HEAD

=======
>>>>>>> be3700e39359a52bf864a0da0c947382fc71a6df
	// "github.com/google/go-containerregistry/pkg/authn"
	// "github.com/google/go-containerregistry/pkg/crane"
	// "github.com/google/go-containerregistry/pkg/name"
	"github.com/imharshita/image-controller/pkg/images"
	appsv1 "k8s.io/api/apps/v1"

	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"strings"
)

// var privateRegistry string = "backupregistry"

// func rename(name string) string {
// 	image := strings.Split(name, ":")
// 	img, version := image[0], image[1]
// 	newName := privateRegistry + "/" + img + ":" + version
// 	return newName
// }

// func retag(imgName string) (name.Tag, error) {
// 	tag, err := name.NewTag(imgName)
// 	if err != nil {
// 		return name.Tag{}, err
// 	}
// 	return tag, nil
// }

// func Process(imgName string) (string, error) {
// 	auth := authn.AuthConfig{
// 		Username: "backupregistry",
// 		Password: "mydockerimages",
// 	}
// 	authenticator := authn.FromConfig(auth)
// 	opt := crane.WithAuth(authenticator)
// 	img, err := crane.Pull(imgName, opt)
// 	if err != nil {
// 		return "", err
// 	}
// 	newName := rename(imgName)
// 	tag, err := retag(newName)
// 	if err != nil {
// 		return "", err
// 	}

// 	if err := crane.Push(img, tag.String(), opt); err != nil {
// 		return "", err
// 	}
// 	return newName, nil
// }

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqNamespace := req.NamespacedName.Namespace
	if reqNamespace != "kube-system" && reqNamespace != "system" {
		_ = r.Log.WithValues("deployment", req.NamespacedName)

		deployments := &appsv1.Deployment{}
		err := r.Get(context.TODO(), req.NamespacedName, deployments)
		if err != nil {
			return reconcile.Result{}, err
		}
		// watch namespace
		//namespaces := deployments.Namespace
		//fmt.Println(namespaces)
		containers := deployments.Spec.Template.Spec.Containers
		for i, c := range containers {
<<<<<<< HEAD
			fmt.Println("deployment image", c.Image)
			if !strings.HasPrefix(c.Image, "backupregistry") {
=======
			fmt.Println(c.Image)
			if !strings.HasPrefix(c.Image, "harshitadocker") {
>>>>>>> be3700e39359a52bf864a0da0c947382fc71a6df
				img, err := images.Process(c.Image)
				if err != nil {
					return ctrl.Result{}, err
				}
				fmt.Println("upated image:", img)
				// Update the Deployment
				deployments.Spec.Template.Spec.Containers[i].Image = img
				err = r.Update(context.TODO(), deployments)
				if err != nil {
					return reconcile.Result{}, err
				}
			}

		}
		return ctrl.Result{}, nil
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