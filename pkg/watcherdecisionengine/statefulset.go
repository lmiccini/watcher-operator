package watcherdecisionengine

import (
	"path/filepath"

	topologyv1 "github.com/openstack-k8s-operators/infra-operator/apis/topology/v1beta1"
	"github.com/openstack-k8s-operators/lib-common/modules/common"
	"github.com/openstack-k8s-operators/lib-common/modules/common/affinity"
	"github.com/openstack-k8s-operators/lib-common/modules/common/env"
	watcherv1beta1 "github.com/openstack-k8s-operators/watcher-operator/api/v1beta1"
	"github.com/openstack-k8s-operators/watcher-operator/pkg/watcher"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

const (
	// KollaServiceCommand -
	KollaServiceCommand = "/usr/local/bin/kolla_start"

	// ComponentName -
	ComponentName = watcher.ServiceName + "-decision-engine"
)

func StatefulSet(
	instance *watcherv1beta1.WatcherDecisionEngine,
	configHash string,
	prometheusCaCertSecret map[string]string,
	labels map[string]string,
	topology *topologyv1.Topology,
) *appsv1.StatefulSet {

	// This allows the pod to start up slowly. The pod will only be killed
	// if it does not succeed a probe in 60 seconds.
	startupProbe := &corev1.Probe{
		FailureThreshold: 6,
		PeriodSeconds:    10,
	}
	// After the first successful startupProbe, livenessProbe takes over
	livenessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds: 30,
		PeriodSeconds:  30,
	}
	readinessProbe := &corev1.Probe{
		// TODO might need tuning
		TimeoutSeconds: 30,
		PeriodSeconds:  30,
	}

	envVars := map[string]env.Setter{}
	envVars["KOLLA_CONFIG_STRATEGY"] = env.SetValue("COPY_ALWAYS")
	envVars["CONFIG_HASH"] = env.SetValue(configHash)
	args := []string{"-c", KollaServiceCommand}

	startupProbe.Exec = &corev1.ExecAction{
		Command: []string{
			"/usr/bin/pgrep", "-f", "-r", "DRST", ComponentName,
		},
	}

	livenessProbe.Exec = &corev1.ExecAction{
		Command: []string{
			"/usr/bin/pgrep", "-f", "-r", "DRST", ComponentName,
		},
	}

	readinessProbe.Exec = &corev1.ExecAction{
		Command: []string{
			"/usr/bin/pgrep", "-f", "-r", "DRST", ComponentName,
		},
	}
	var config0644AccessMode int32 = 0644
	volumes := append(watcher.GetLogVolume(),
		corev1.Volume{
			Name: "config-data-custom",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &config0644AccessMode,
					SecretName:  instance.Name + "-config-data",
				},
			},
		},
	)

	volumeMounts := []corev1.VolumeMount{
		watcher.GetKollaConfigVolumeMount(ComponentName),
	}
	volumeMounts = append(volumeMounts, watcher.GetLogVolumeMount()...)

	// Create mount for bundle CA if defined in TLS.CaBundleSecretName
	if instance.Spec.TLS.CaBundleSecretName != "" {
		volumes = append(volumes, instance.Spec.TLS.CreateVolume())
		volumeMounts = append(volumeMounts, instance.Spec.TLS.CreateVolumeMounts(nil)...)
	}

	if len(prometheusCaCertSecret) != 0 {
		volumes = append(volumes,
			corev1.Volume{
				Name: "custom-prometheus-ca",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: prometheusCaCertSecret["casecret_name"],
					},
				},
			},
		)
		volumeMounts = append(volumeMounts,
			corev1.VolumeMount{
				Name:      "custom-prometheus-ca",
				MountPath: filepath.Join(watcher.PrometheusCaCertFolderPath, prometheusCaCertSecret["casecret_key"]),
				SubPath:   prometheusCaCertSecret["casecret_key"],
				ReadOnly:  true,
			},
		)
	}

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: instance.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: instance.Spec.ServiceAccount,
					Containers: []corev1.Container{
						{
							Name: ComponentName,
							Command: []string{
								"/bin/bash",
							},
							Args:  args,
							Image: instance.Spec.ContainerImage,
							SecurityContext: &corev1.SecurityContext{
								RunAsUser: ptr.To(watcher.WatcherUserID),
							},
							Env: env.MergeEnvs([]corev1.EnvVar{}, envVars),
							VolumeMounts: append(watcher.GetVolumeMounts(
								[]string{}),
								volumeMounts...,
							),
							Resources:      instance.Spec.Resources,
							StartupProbe:   startupProbe,
							ReadinessProbe: readinessProbe,
							LivenessProbe:  livenessProbe,
						},
					},
				},
			},
		},
	}

	if topology != nil {
		topology.ApplyTo(&statefulset.Spec.Template)
	} else {
		// If possible two pods of the same service should not
		// run on the same worker node. If this is not possible
		// the get still created on the same worker node.
		statefulset.Spec.Template.Spec.Affinity = affinity.DistributePods(
			common.AppSelector,
			[]string{
				instance.Name,
			},
			corev1.LabelHostname,
		)
	}

	statefulset.Spec.Template.Spec.Volumes = append(watcher.GetVolumes(
		instance.Name,
		[]string{}),
		volumes...)

	if instance.Spec.NodeSelector != nil {
		statefulset.Spec.Template.Spec.NodeSelector = *instance.Spec.NodeSelector
	}

	return statefulset

}
