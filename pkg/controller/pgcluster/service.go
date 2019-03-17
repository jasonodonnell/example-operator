package pgcluster

import (
	pgclusterv1alpha1 "github.com/jasonodonnell/example-operator/pkg/apis/pgcluster/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcilePGCluster) newServiceForPGCluster(cr *pgclusterv1alpha1.PGCluster) *corev1.Service {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:     "postgres",
					Port:     int32(5432),
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"name": cr.Name,
			},
		},
	}

	controllerutil.SetControllerReference(cr, service, r.scheme)
	return service
}
