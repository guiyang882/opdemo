package resources

import (
	farmv1 "github.com/liuguiyangnwpu/opdemo/api/v1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewService(manor *farmv1.Manor) *corev1.Service {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
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
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeNodePort,
			Ports: manor.Spec.Ports,
			Selector: map[string]string{
				"app": manor.Name,
			},
		},
	}
	return service
}
