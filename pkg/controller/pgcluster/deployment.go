package pgcluster

import (
	"fmt"

	pgclusterv1alpha1 "github.com/jasonodonnell/example-operator/pkg/apis/pgcluster/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcilePGCluster) newDeploymentForPGCluster(cr *pgclusterv1alpha1.PGCluster) *appsv1.Deployment {
	image := fmt.Sprintf("crunchydata/crunchy-postgres:%s", cr.Spec.ImageTag)
	ls := labelsForPGCluster(cr.Name)
	replicas := int32(1)

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
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
					Containers: []corev1.Container{
						{
							Name:  "postgres",
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "postgres",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 5432,
								},
							},
							Env: setEnvs(cr),
						},
					},
				},
			},
		},
	}
	controllerutil.SetControllerReference(cr, dep, r.scheme)
	return dep
}

func labelsForPGCluster(name string) map[string]string {
	return map[string]string{
		"app":          "postgres",
		"pgcluster_cr": name,
		"name":         name,
	}
}

func setEnvs(cr *pgclusterv1alpha1.PGCluster) []corev1.EnvVar {
	defaults := getDefaultVars(cr.Spec.Mode, cr.Spec.Database, cr.Spec.Username)
	if cr.Spec.Mode == "replica" {
		vars := setReplicaVars(cr.Spec.PrimaryHost, cr.Spec.PrimaryPort)
		defaults = append(defaults, vars...)
	}
	fmt.Println(defaults)
	return defaults
}

func setReplicaVars(host, port string) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "PG_PRIMARY_HOST",
			Value: host,
		},
		{
			Name:  "PG_PRIMARY_PORT",
			Value: port,
		},
	}
}

func getDefaultVars(mode, db, user string) []corev1.EnvVar {
	vars := []corev1.EnvVar{
		{
			Name:  "PG_MODE",
			Value: mode,
		},
		{
			Name:  "PG_DATABASE",
			Value: db,
		},
		{
			Name:  "PG_USER",
			Value: user,
		},
		{
			Name:  "PG_PRIMARY_PORT",
			Value: "5432",
		},
		{
			Name:  "PGDATA_PATH_OVERRIDE",
			Value: "primary",
		},
		{
			Name:  "PGHOST",
			Value: "/tmp",
		},
		{
			Name:  "PG_PASSWORD",
			Value: "password",
		},
		{
			Name:  "PG_PRIMARY_USER",
			Value: "replication",
		},
		{
			Name:  "PG_PRIMARY_PASSWORD",
			Value: "password",
		},
		{
			Name:  "PG_ROOT_PASSWORD",
			Value: "password",
		},
	}
	return vars
}
