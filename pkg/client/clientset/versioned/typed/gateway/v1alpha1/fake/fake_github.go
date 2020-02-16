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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/n3wscott/gateway/pkg/apis/gateway/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeGitHubs implements GitHubInterface
type FakeGitHubs struct {
	Fake *FakeGatewayV1alpha1
	ns   string
}

var githubsResource = schema.GroupVersionResource{Group: "gateway.n3wscott.com", Version: "v1alpha1", Resource: "githubs"}

var githubsKind = schema.GroupVersionKind{Group: "gateway.n3wscott.com", Version: "v1alpha1", Kind: "GitHub"}

// Get takes name of the gitHub, and returns the corresponding gitHub object, and an error if there is any.
func (c *FakeGitHubs) Get(name string, options v1.GetOptions) (result *v1alpha1.GitHub, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(githubsResource, c.ns, name), &v1alpha1.GitHub{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GitHub), err
}

// List takes label and field selectors, and returns the list of GitHubs that match those selectors.
func (c *FakeGitHubs) List(opts v1.ListOptions) (result *v1alpha1.GitHubList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(githubsResource, githubsKind, c.ns, opts), &v1alpha1.GitHubList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.GitHubList{ListMeta: obj.(*v1alpha1.GitHubList).ListMeta}
	for _, item := range obj.(*v1alpha1.GitHubList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gitHubs.
func (c *FakeGitHubs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(githubsResource, c.ns, opts))

}

// Create takes the representation of a gitHub and creates it.  Returns the server's representation of the gitHub, and an error, if there is any.
func (c *FakeGitHubs) Create(gitHub *v1alpha1.GitHub) (result *v1alpha1.GitHub, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(githubsResource, c.ns, gitHub), &v1alpha1.GitHub{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GitHub), err
}

// Update takes the representation of a gitHub and updates it. Returns the server's representation of the gitHub, and an error, if there is any.
func (c *FakeGitHubs) Update(gitHub *v1alpha1.GitHub) (result *v1alpha1.GitHub, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(githubsResource, c.ns, gitHub), &v1alpha1.GitHub{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GitHub), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeGitHubs) UpdateStatus(gitHub *v1alpha1.GitHub) (*v1alpha1.GitHub, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(githubsResource, "status", c.ns, gitHub), &v1alpha1.GitHub{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GitHub), err
}

// Delete takes name of the gitHub and deletes it. Returns an error if one occurs.
func (c *FakeGitHubs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(githubsResource, c.ns, name), &v1alpha1.GitHub{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGitHubs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(githubsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.GitHubList{})
	return err
}

// Patch applies the patch and returns the patched gitHub.
func (c *FakeGitHubs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GitHub, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(githubsResource, c.ns, name, pt, data, subresources...), &v1alpha1.GitHub{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GitHub), err
}