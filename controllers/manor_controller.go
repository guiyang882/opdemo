/*
Copyright 2020 wuming.lgy@alibaba-inc.com.

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
	"encoding/json"
	"github.com/liuguiyangnwpu/opdemo/controllers/resources"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	farmv1 "github.com/liuguiyangnwpu/opdemo/api/v1"
)

// ManorReconciler reconciles a Manor object
type ManorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=farm.example.com,resources=manors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=farm.example.com,resources=manors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=farm.example.com,resources=manors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Manor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ManorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("manor", req.NamespacedName)

	// your logic here
	instance := &farmv1.Manor{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	if instance.DeletionTimestamp != nil {
		return ctrl.Result{}, err
	}

	// 如果不存在，则创建关联资源
	// 如果存在，判断是否需要更新
	deploy := &v1.Deployment{}
	if err := r.Client.Get(context.TODO(), req.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		// 创建关联资源
		deploy = resources.NewDeploy(instance)
		if errCreate := r.Client.Create(context.TODO(), deploy); errCreate != nil {
			return ctrl.Result{}, errCreate
		}
		// 创建Service
		service := resources.NewService(instance)
		if errCreate := r.Client.Create(context.TODO(), service); errCreate != nil {
			return ctrl.Result{}, errCreate
		}
		// 关联Annotations
		data, _ := json.Marshal(instance.Spec)
		if instance.Annotations != nil {
			instance.Annotations["spec"] = string(data)
		} else {
			instance.Annotations = map[string]string{
				"spec": string(data),
			}
		}
		// 更新对应的资源
		if errUpdate := r.Client.Update(context.TODO(), instance); errUpdate != nil {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, nil
	}

	oldSpec := farmv1.ManorSpec{}
	if err := json.Unmarshal([]byte(instance.Annotations["spec"]), &oldSpec); err != nil {
		return ctrl.Result{}, err
	}

	if !reflect.DeepEqual(instance.Spec, oldSpec) {
		// 更新关联资源
		newDeploy := resources.NewDeploy(instance)
		oldDeploy := &v1.Deployment{}
		if errGet := r.Client.Get(context.TODO(), req.NamespacedName, oldDeploy); errGet != nil {
			return ctrl.Result{}, errGet
		}
		oldDeploy.Spec = newDeploy.Spec
		if errUpdate := r.Client.Update(context.TODO(), oldDeploy); errUpdate != nil {
			return ctrl.Result{}, errUpdate
		}
		newService := resources.NewService(instance)
		oldService := &corev1.Service{}
		if errGet := r.Client.Get(context.TODO(), req.NamespacedName, oldService); errGet != nil {
			return ctrl.Result{}, errGet
		}
		oldService.Spec = newService.Spec
		if errUpdate := r.Client.Update(context.TODO(), oldService); errUpdate != nil {
			return ctrl.Result{}, errUpdate
		}
		return ctrl.Result{}, nil
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ManorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&farmv1.Manor{}).
		Complete(r)
}
