package resources

import (
	farmv1 "github.com/liuguiyangnwpu/opdemo/api/v1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewDeploy(manor *farmv1.Manor) *v1.Deployment {
	labels := map[string]string{
		"app": manor.Name,
	}
	selector := &metav1.LabelSelector{
		MatchLabels:      labels,
		MatchExpressions: nil,
	}

	deploy := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      manor.Name,
			Namespace: manor.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(manor, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "Manor",
				}),
			},
		},
		Spec: v1.DeploymentSpec{
			Replicas: manor.Spec.Size,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: newContainers(manor),
				},
			},
		},
	}
	return deploy
}

func newContainers(manor *farmv1.Manor) []corev1.Container {
	containerPorts := make([]corev1.ContainerPort, 0)
	for _, svcPort := range manor.Spec.Ports {
		cPort := corev1.ContainerPort{}
		cPort.ContainerPort = svcPort.TargetPort.IntVal
		containerPorts = append(containerPorts, cPort)
	}
	return []corev1.Container{
		corev1.Container{
			Name:            manor.Name,
			Image:           manor.Spec.Image,
			Ports:           containerPorts,
			Env:             manor.Spec.Envs,
			Resources:       manor.Spec.Resources,
			ImagePullPolicy: corev1.PullAlways,
		},
	}
}
