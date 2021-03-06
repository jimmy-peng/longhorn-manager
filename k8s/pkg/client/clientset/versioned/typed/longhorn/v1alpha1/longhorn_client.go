/*
Copyright 2017 The Kubernetes Authors.

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

package v1alpha1

import (
	v1alpha1 "github.com/rancher/longhorn-manager/k8s/pkg/apis/longhorn/v1alpha1"
	"github.com/rancher/longhorn-manager/k8s/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type LonghornV1alpha1Interface interface {
	RESTClient() rest.Interface
	ControllersGetter
	NodesGetter
	ReplicasGetter
	SettingsGetter
	VolumesGetter
}

// LonghornV1alpha1Client is used to interact with features provided by the longhorn.rancher.io group.
type LonghornV1alpha1Client struct {
	restClient rest.Interface
}

func (c *LonghornV1alpha1Client) Controllers(namespace string) ControllerInterface {
	return newControllers(c, namespace)
}

func (c *LonghornV1alpha1Client) Nodes(namespace string) NodeInterface {
	return newNodes(c, namespace)
}

func (c *LonghornV1alpha1Client) Replicas(namespace string) ReplicaInterface {
	return newReplicas(c, namespace)
}

func (c *LonghornV1alpha1Client) Settings(namespace string) SettingInterface {
	return newSettings(c, namespace)
}

func (c *LonghornV1alpha1Client) Volumes(namespace string) VolumeInterface {
	return newVolumes(c, namespace)
}

// NewForConfig creates a new LonghornV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*LonghornV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &LonghornV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new LonghornV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *LonghornV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new LonghornV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *LonghornV1alpha1Client {
	return &LonghornV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *LonghornV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
