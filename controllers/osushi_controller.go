/*


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
	"strings"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	cachev1alpha1 "github.com/iaoiui/osushi/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OsushiReconciler reconciles a Osushi object
type OsushiReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const esdressOsushi string = "endressOsushi"
const traditionalKaitenSushi string = "traditionalKaitenSushi"

// +kubebuilder:rbac:groups=cache.my.domain,resources=osushis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.my.domain,resources=osushis/status,verbs=get;update;patch
// 任意の長さのosushiの絵文字を生成する
func generateOsushiEmoji(lengthOfOsushiLane int) string {
	return strings.Repeat(` `, lengthOfOsushiLane-1) + `🍣`
}

// Reconcile はメインのReconcillation Loopです
func (r *OsushiReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("osushi", req.NamespacedName)

	instance := &cachev1alpha1.Osushi{}
	err := r.Get(ctx, req.NamespacedName, instance)
	instanceTmp := instance.DeepCopy()

	// For the issue: https://github.com/kubernetes/kubernetes/issues/84430#issuecomment-638376994
	// time.Sleep(30)

	if err != nil {
		fmt.Println(err)
		return ctrl.Result{}, nil
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	result, err := r.reconcileDeployment(ctx, instanceTmp, req)
	if err != nil {
		fmt.Println(err)
		return reconcile.Result{}, err
	}
	if result.Requeue == true {
		return reconcile.Result{Requeue: true}, nil
	}
	return r.reconcileOsushi(ctx, instanceTmp, req)
}

type osushiAnnotator struct {
	Client  client.Client
	decoder *admission.Decoder
}

var defaultLengthOfOsushiLane int32 = 20

// reconcile Deployment
func (r *OsushiReconciler) reconcileDeployment(ctx context.Context, instance *cachev1alpha1.Osushi, req ctrl.Request) (reconcile.Result, error) {
	logger := r.Log.WithValues("osushi", req.NamespacedName)
	// deploymentの存在をチェックして、既に存在していれば生成しない
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// 新deploymentの定義
		dep := r.deploymentForOsushi(instance)
		logger.Info(fmt.Sprintln("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name))
		err = r.Create(ctx, dep)
		if err != nil {
			logger.Error(err, fmt.Sprintln("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name))
			return reconcile.Result{}, err
		}
		// Deployment can be generated
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		logger.Error(err, "Failed to get Deployment.")
		return reconcile.Result{}, err
	}

	// デプロイメントのサイズをspecと同じになるように調整する
	size := instance.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			logger.Error(err, fmt.Sprintln("Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name))
			return reconcile.Result{}, err
		}
		// Specを更新
		return reconcile.Result{Requeue: true}, nil
	}
	return reconcile.Result{}, nil
}

//お寿司に関するreconcileを行う
func (r *OsushiReconciler) reconcileOsushi(ctx context.Context, instance *cachev1alpha1.Osushi, req ctrl.Request) (reconcile.Result, error) {
	logger := r.Log.WithValues("osushi", req.NamespacedName)

	//回転寿司の速さ
	// TODO instance.Spec.osushiSpeed, instance.Spec.lengthOfOsushiLaneが空かどうかチェック
	osushiSpeed := int(instance.Spec.OsushiSpeed)

	//その場合、モードの併用はできなくなるがいいか
	switch instance.Spec.Mode {
	case esdressOsushi:
		//エンドレスお寿司モードの場合
		// instance.Spec.Emoji = strings.Repeat(`🍣`, int(instance.Spec.Size))
		instance.Spec.Emoji = strings.Repeat(`🍣`, int(instance.Spec.Size))
		// err := r.Update(ctx, instance)
		err := r.Update(ctx, instance)

		if err != nil {
			logger.Error(err, fmt.Sprintf("Failed to update Osushi %v \n", string(instance.Name)))
			return reconcile.Result{}, err
		}
		// Specを更新
		return reconcile.Result{Requeue: true}, nil
	case traditionalKaitenSushi:
		// 回転寿司モード
		// TODO 2行分を表示できるか試す
		tmp := instance.Spec.Emoji

		if strings.Contains(tmp, `🍣`) {
			// 2 全ての文字列を1つ左にずらす
			tmpp := string(tmp[osushiSpeed:]) + string(tmp[0:osushiSpeed])
			tmp = tmpp
			time.Sleep(500 * time.Millisecond)
			index := strings.Index(tmp, `🍣`)
			if index < 0 {
				//reset Osushi position
				tmp = generateOsushiEmoji(int(instance.Spec.LengthOfOsushiLane))
				logger.Info("お寿司をリセット")
			}

			logger.Info(fmt.Sprintf("お寿司の位置は %v 番目", index))
			//お寿司が届く範囲にあるかどうか確認
			if index < osushiSpeed {
				instance.Status.Reachable = true
				logger.Info("とどく！")
			} else {
				instance.Status.Reachable = false
				logger.Info("とどかない...")
			}

		} else {
			// 1 お寿司を含んでいなければ、初期状態のお寿司をセットする
			tmp = generateOsushiEmoji(int(instance.Spec.LengthOfOsushiLane))
			logger.Info("お寿司をセット")
		}
		instance.Status.Reachable = true
		// 3 instance.Spec.Emoji にその文字列を突っ込む
		instance.Spec.Emoji = tmp

		a := instance.DeepCopy()
		err := r.Update(ctx, a)
		if err != nil {
			logger.Error(err, fmt.Sprintln("Failed to update Osushi.", string(instance.Name)))
			return reconcile.Result{RequeueAfter: time.Second * 5}, nil
		}
		// Specを更新
		return reconcile.Result{RequeueAfter: time.Second * 5}, nil

	default:
		// どのモードでもない場合
		if instance.Spec.Emoji != `🍣` {
			//お寿司絵文字を1つセットする
			instance.Spec.Emoji = `🍣`
			err := r.Update(ctx, instance)
			if err != nil {
				logger.Error(err, fmt.Sprintln("Failed to update Osushi.", string(instance.Name)))
				return reconcile.Result{}, err
			}
			// Specを更新
			return reconcile.Result{Requeue: true}, nil
		}

	}

	// TODO お寿司の新鮮さを更新する
	// if instance.get
	// instance.Status.Freshness = "新鮮"

	// ポッド名と共にOsushiの状態を更新
	// Osushi's deployment のポッドをリストする
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForOsushi(instance.Name))
	listOps := &client.ListOptions{
		Namespace:     instance.Namespace,
		LabelSelector: labelSelector,
	}
	err := r.List(ctx, podList, listOps)
	if err != nil {
		logger.Error(err, fmt.Sprintln("Failed to list pods.", "Osushi.Namespace", instance.Namespace, "Osushi.Name", instance.Name))
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

// osushi デプロイメントを生成する
func (r *OsushiReconciler) deploymentForOsushi(m *cachev1alpha1.Osushi) *appsv1.Deployment {
	ls := labelsForOsushi(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "busybox",
						Name:    "osushi",
						Command: []string{"sleep", "3600"},
					}},
				},
			},
		},
	}
	// Osushiインスタンスをオーナーとコントローラとしてセットする
	controllerutil.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// labelsForOsushi は、セレクター用ラベルを返す
// カスタムリソースの一員として与えられたラベルも返す
func labelsForOsushi(name string) map[string]string {
	return map[string]string{"app": "osushi", "osushi_cr": name}
}

// getPodNames ポッドの名前の配列を返す
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

func (r *OsushiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Osushi{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
