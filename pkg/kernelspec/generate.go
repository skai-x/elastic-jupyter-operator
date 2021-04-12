// Tencent is pleased to support the open source community by making TKEStack
// available.
//
// Copyright (C) 2012-2020 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License. You may obtain a copy of the
// License at
//
// https://opensource.org/licenses/Apache-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.
// Tencent is pleased to support the open source community by making TKEStack
// available.
//
// Copyright (C) 2012-2020 Tencent. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License. You may obtain a copy of the
// License at
//
// https://opensource.org/licenses/Apache-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OF ANY KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations under the License.

package kernelspec

import (
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tkestack/elastic-jupyter-operator/api/v1alpha1"
)

const (
	LabelKernelSpec = "kernelspec"
	LabelNS         = "namespace"

	fileName         = "kernel.json"
	defaultClassName = "enterprise_gateway.services.processproxies.k8s.KubernetesProcessProxy"
)

// generator defines the generator which is used to generate
// desired specs.
type generator struct {
	kernelSpec *v1alpha1.JupyterKernelSpec
}

// newGenerator creates a new Generator.
func newGenerator(kernelSpec *v1alpha1.JupyterKernelSpec) (
	*generator, error) {
	if kernelSpec == nil {
		return nil, fmt.Errorf("Got nil when initializing Generator")
	}
	g := &generator{
		kernelSpec: kernelSpec,
	}

	return g, nil
}

func (g generator) DesiredConfigmapWithoutOwner() (*v1.ConfigMap, error) {
	labels := g.labels()

	jsonConfig, err := g.desiredJSON()
	if err != nil {
		return nil, err
	}

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: g.kernelSpec.Namespace,
			Name:      g.kernelSpec.Name,
			Labels:    labels,
		},
		Data: map[string]string{
			fileName: jsonConfig,
		},
	}

	return cm, nil
}

func (g generator) desiredJSON() (string, error) {
	c := &kernelConfig{
		Language:    g.kernelSpec.Spec.Language,
		DisplayName: g.kernelSpec.Spec.DisplayName,
		Metadata: metadata{
			Config: config{
				ImageName: g.kernelSpec.Spec.Image,
			},
			ProcessProxy: processProxy{
				ClassName: defaultClassName,
			},
		},
		Argv: g.kernelSpec.Spec.Command,
	}
	v, err := json.Marshal(c)
	return string(v), err
}

func (g generator) labels() map[string]string {
	return map[string]string{
		LabelNS:         g.kernelSpec.Namespace,
		LabelKernelSpec: g.kernelSpec.Name,
	}
}
