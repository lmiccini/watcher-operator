/*
Copyright 2024.

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

package v1beta1

import (
	topologyv1 "github.com/openstack-k8s-operators/infra-operator/apis/topology/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var watcherapilog = logf.Log.WithName("watcherapi-resource")

func (r *WatcherAPI) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-watcher-openstack-org-v1beta1-watcherapi,mutating=true,failurePolicy=fail,sideEffects=None,groups=watcher.openstack.org,resources=watcherapis,verbs=create;update,versions=v1beta1,name=mwatcherapi.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &WatcherAPI{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *WatcherAPI) Default() {
	watcherapilog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-watcher-openstack-org-v1beta1-watcherapi,mutating=false,failurePolicy=fail,sideEffects=None,groups=watcher.openstack.org,resources=watcherapis,verbs=create;update,versions=v1beta1,name=vwatcherapi.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &WatcherAPI{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *WatcherAPI) ValidateCreate() (admission.Warnings, error) {
	watcherapilog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *WatcherAPI) ValidateUpdate(runtime.Object) (admission.Warnings, error) {
	watcherapilog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *WatcherAPI) ValidateDelete() (admission.Warnings, error) {
	watcherapilog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}

// ValidateTopology validates the referenced TopoRef.Namespace.
func (r *WatcherAPITemplate) ValidateTopology(
	basePath *field.Path,
	namespace string,
) field.ErrorList {
	return topologyv1.ValidateTopologyRef(
		r.TopologyRef,
		*basePath.Child("topologyRef"),
		namespace,
	)
}
