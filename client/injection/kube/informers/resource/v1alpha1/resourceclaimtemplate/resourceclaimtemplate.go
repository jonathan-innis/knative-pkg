/*
Copyright 2022 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package resourceclaimtemplate

import (
	context "context"

	apiresourcev1alpha1 "k8s.io/api/resource/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	v1alpha1 "k8s.io/client-go/informers/resource/v1alpha1"
	kubernetes "k8s.io/client-go/kubernetes"
	resourcev1alpha1 "k8s.io/client-go/listers/resource/v1alpha1"
	cache "k8s.io/client-go/tools/cache"
	client "knative.dev/pkg/client/injection/kube/client"
	factory "knative.dev/pkg/client/injection/kube/informers/factory"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Resource().V1alpha1().ResourceClaimTemplates()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx), resourceVersion: injection.GetResourceVersion(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1alpha1.ResourceClaimTemplateInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch k8s.io/client-go/informers/resource/v1alpha1.ResourceClaimTemplateInformer from context.")
	}
	return untyped.(v1alpha1.ResourceClaimTemplateInformer)
}

type wrapper struct {
	client kubernetes.Interface

	namespace string

	resourceVersion string
}

var _ v1alpha1.ResourceClaimTemplateInformer = (*wrapper)(nil)
var _ resourcev1alpha1.ResourceClaimTemplateLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apiresourcev1alpha1.ResourceClaimTemplate{}, 0, nil)
}

func (w *wrapper) Lister() resourcev1alpha1.ResourceClaimTemplateLister {
	return w
}

func (w *wrapper) ResourceClaimTemplates(namespace string) resourcev1alpha1.ResourceClaimTemplateNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace, resourceVersion: w.resourceVersion}
}

// SetResourceVersion allows consumers to adjust the minimum resourceVersion
// used by the underlying client.  It is not accessible via the standard
// lister interface, but can be accessed through a user-defined interface and
// an implementation check e.g. rvs, ok := foo.(ResourceVersionSetter)
func (w *wrapper) SetResourceVersion(resourceVersion string) {
	w.resourceVersion = resourceVersion
}

func (w *wrapper) List(selector labels.Selector) (ret []*apiresourcev1alpha1.ResourceClaimTemplate, err error) {
	lo, err := w.client.ResourceV1alpha1().ResourceClaimTemplates(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector:   selector.String(),
		ResourceVersion: w.resourceVersion,
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apiresourcev1alpha1.ResourceClaimTemplate, error) {
	return w.client.ResourceV1alpha1().ResourceClaimTemplates(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		ResourceVersion: w.resourceVersion,
	})
}
