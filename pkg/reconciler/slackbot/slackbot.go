/*
Copyright 2020 The Knative Authors

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

package slackbot

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"

	"github.com/n3wscott/gateway/pkg/apis/gateway/v1alpha1"
	versioned "github.com/n3wscott/gateway/pkg/client/clientset/versioned"
	reconcilerslack "github.com/n3wscott/gateway/pkg/client/injection/reconciler/gateway/v1alpha1/slackbot"
	listers "github.com/n3wscott/gateway/pkg/client/listers/gateway/v1alpha1"
	"github.com/n3wscott/gateway/pkg/reconciler/slackbot/resources"
	corev1listers "k8s.io/client-go/listers/core/v1"
	eventinglisters "knative.dev/eventing/pkg/client/listers/eventing/v1alpha1"
)

// newReconciledNormal makes a new reconciler event with event type Normal, and
// reason SlackReconciled.
func newReconciledNormal(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, "SlackReconciled", "Slackbot reconciled: \"%s/%s\"", namespace, name)
}

// newDeploymentCreated makes a new reconciler event with event type Normal, and
// reason SlackDeploymentCreated.
func newDeploymentCreated(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, "SlackDeploymentCreated", "Slackbot created deployment: \"%s/%s\"", namespace, name)
}

// newDeploymentFailed makes a new reconciler event with event type Warning, and
// reason SlackDeploymentFailed.
func newDeploymentFailed(namespace, name string, err error) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeWarning, "SlackDeploymentFailed", "Slackbot failed to create deployment: \"%s/%s\", %w", namespace, name, err)
}

// newDeploymentUpdated makes a new reconciler event with event type Normal, and
// reason SlackDeploymentUpdated.
func newDeploymentUpdated(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, "SlackDeploymentUpdated", "Slackbot updated deployment: \"%s/%s\"", namespace, name)
}

// Reconciler reconciles a Slackbot object
type Reconciler struct {
	// KubeClientSet allows us to talk to the k8s for core APIs
	kubeClientSet kubernetes.Interface

	ReceiveAdapterImage string `envconfig:"SLACKBOT_IMAGE" required:"true"`

	// listers index properties about resources
	slackbotLister        listers.SlackbotLister
	deploymentLister      appsv1listers.DeploymentLister
	eventTypeLister       eventinglisters.EventTypeLister
	serviceLister         corev1listers.ServiceLister
	samplesourceClientSet versioned.Interface

	sinkResolver *resolver.URIResolver
}

// Check that our Reconciler implements Interface
var _ reconcilerslack.Interface = (*Reconciler)(nil)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, source *v1alpha1.Slackbot) pkgreconciler.Event {
	source.Status.InitializeConditions()
	source.Status.ObservedGeneration = source.Generation

	dest := source.Spec.Sink.DeepCopy()
	if dest.Ref != nil {
		// To call URIFromDestination(), dest.Ref must have a Namespace. If there is
		// no Namespace defined in dest.Ref, we will use the Namespace of the source
		// as the Namespace of dest.Ref.
		if dest.Ref.Namespace == "" {
			//TODO how does this work with deprecated fields
			dest.Ref.Namespace = source.GetNamespace()
		}
	}

	sinkURI, err := r.sinkResolver.URIFromDestinationV1(*dest, source)
	if err != nil {
		source.Status.MarkNoSink("NotFound", "")
		return err
	}
	source.Status.MarkSink(sinkURI)

	ra, event := r.createReceiveAdapter(ctx, source, sinkURI)
	// Update source status
	if ra != nil {
		source.Status.PropagateDeploymentAvailability(ra)
	}
	if event != nil {
		return event
	}

	logging.FromContext(ctx).Infow("Reconciling Service")
	if err := r.reconcileService(ctx, source); err != nil {
		return err
	}

	return newReconciledNormal(source.Namespace, source.Name)
}

func (r *Reconciler) createReceiveAdapter(ctx context.Context, src *v1alpha1.Slackbot, sinkURI *apis.URL) (*appsv1.Deployment, pkgreconciler.Event) {
	eventSource := r.makeEventSource(src)
	logging.FromContext(ctx).Debug("event source", zap.Any("source", eventSource))

	adapterArgs := resources.ReceiveAdapterArgs{
		EventSource: eventSource,
		Image:       r.ReceiveAdapterImage,
		Source:      src,
		Labels:      resources.Labels(src.Name),
		SinkURI:     sinkURI,
	}
	expected := resources.MakeReceiveAdapter(&adapterArgs)

	ra, err := r.kubeClientSet.AppsV1().Deployments(src.Namespace).Get(expected.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		ra, err = r.kubeClientSet.AppsV1().Deployments(src.Namespace).Create(expected)
		if err != nil {
			return nil, newDeploymentFailed(expected.Namespace, expected.Name, err)
		}
		return ra, newDeploymentCreated(ra.Namespace, ra.Name)
	} else if err != nil {
		return nil, fmt.Errorf("error getting receive adapter: %v", err)
	} else if !metav1.IsControlledBy(ra, src) {
		return nil, fmt.Errorf("deployment %q is not owned by Slackbot %q", ra.Name, src.Name)
	} else if r.podSpecChanged(ra.Spec.Template.Spec, expected.Spec.Template.Spec) {
		ra.Spec.Template.Spec = expected.Spec.Template.Spec
		if ra, err = r.kubeClientSet.AppsV1().Deployments(src.Namespace).Update(ra); err != nil {
			return ra, err
		}
		return ra, nil //newDeploymentUpdated(ra.Namespace, ra.Name)
	} else {
		logging.FromContext(ctx).Debugw("Reusing existing receive adapter", zap.Any("receiveAdapter", ra))
	}
	return ra, nil
}

func (r *Reconciler) podSpecChanged(oldPodSpec corev1.PodSpec, newPodSpec corev1.PodSpec) bool {
	if !equality.Semantic.DeepDerivative(newPodSpec, oldPodSpec) {
		return true
	}
	if len(oldPodSpec.Containers) != len(newPodSpec.Containers) {
		return true
	}
	for i := range newPodSpec.Containers {
		if !equality.Semantic.DeepEqual(newPodSpec.Containers[i].Env, oldPodSpec.Containers[i].Env) {
			return true
		}
	}
	return false
}

func (r *Reconciler) getReceiveAdapter(ctx context.Context, src *v1alpha1.Slackbot) (*appsv1.Deployment, error) {
	dl, err := r.kubeClientSet.AppsV1().Deployments(src.Namespace).List(metav1.ListOptions{
		LabelSelector: r.getLabelSelector(src).String(),
	})
	if err != nil {
		logging.FromContext(ctx).Errorw("Unable to list deployments: %v", zap.Error(err))
		return nil, err
	}
	for _, dep := range dl.Items {
		if metav1.IsControlledBy(&dep, src) {
			return &dep, nil
		}
	}
	return nil, apierrors.NewNotFound(schema.GroupResource{}, "")
}

func (r *Reconciler) getLabelSelector(src *v1alpha1.Slackbot) labels.Selector {
	return labels.SelectorFromSet(resources.Labels(src.Name))
}

func (r *Reconciler) reconcileService(ctx context.Context, sb *v1alpha1.Slackbot) error {
	svc, err := r.getService(ctx, sb, labels.SelectorFromSet(resources.Labels(sb.Name)))

	if apierrors.IsNotFound(err) {
		svc = resources.MakeService(resources.ReceiveAdapterArgs{
			Source: sb,
		})

		var err error
		svc, err = r.kubeClientSet.CoreV1().Services(sb.Namespace).Create(svc)
		if err != nil || svc == nil {
			msg := "Failed to make Service."
			if err != nil {
				msg = msg + " " + err.Error()
			}
			sb.Status.MarkAddress(nil)
			return fmt.Errorf("failed to create Job: %s", err)
		}
	} else if err != nil {
		sb.Status.MarkAddress(nil)
		return err
	}

	url := &apis.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s.%s.svc.cluster.local", svc.Name, svc.Namespace),
	}

	sb.Status.MarkAddress(url)
	return nil
}

func (r *Reconciler) getService(ctx context.Context, owner metav1.Object, ls labels.Selector) (*corev1.Service, error) {
	list, err := r.serviceLister.List(ls)
	if err != nil {
		return nil, err
	}

	for _, i := range list {
		if metav1.IsControlledBy(i, owner) {
			return i, nil
		}
	}

	return nil, apierrors.NewNotFound(schema.GroupResource{}, "")
}

// makeEventSource computes the Cloud Event source attribute for the given source
func (r *Reconciler) makeEventSource(src *v1alpha1.Slackbot) string {
	return src.Namespace + "/" + src.Name
}
