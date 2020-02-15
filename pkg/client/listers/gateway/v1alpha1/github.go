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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/n3wscott/gateway/pkg/apis/gateway/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GitHubLister helps list GitHubs.
type GitHubLister interface {
	// List lists all GitHubs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.GitHub, err error)
	// GitHubs returns an object that can list and get GitHubs.
	GitHubs(namespace string) GitHubNamespaceLister
	GitHubListerExpansion
}

// gitHubLister implements the GitHubLister interface.
type gitHubLister struct {
	indexer cache.Indexer
}

// NewGitHubLister returns a new GitHubLister.
func NewGitHubLister(indexer cache.Indexer) GitHubLister {
	return &gitHubLister{indexer: indexer}
}

// List lists all GitHubs in the indexer.
func (s *gitHubLister) List(selector labels.Selector) (ret []*v1alpha1.GitHub, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GitHub))
	})
	return ret, err
}

// GitHubs returns an object that can list and get GitHubs.
func (s *gitHubLister) GitHubs(namespace string) GitHubNamespaceLister {
	return gitHubNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GitHubNamespaceLister helps list and get GitHubs.
type GitHubNamespaceLister interface {
	// List lists all GitHubs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.GitHub, err error)
	// Get retrieves the GitHub from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.GitHub, error)
	GitHubNamespaceListerExpansion
}

// gitHubNamespaceLister implements the GitHubNamespaceLister
// interface.
type gitHubNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GitHubs in the indexer for a given namespace.
func (s gitHubNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GitHub, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GitHub))
	})
	return ret, err
}

// Get retrieves the GitHub from the indexer for a given namespace and name.
func (s gitHubNamespaceLister) Get(name string) (*v1alpha1.GitHub, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("github"), name)
	}
	return obj.(*v1alpha1.GitHub), nil
}
