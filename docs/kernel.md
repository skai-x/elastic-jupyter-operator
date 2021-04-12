# Design Proposal for Jupyter Kernel CRD

Authors:
- Ce Gao <cegao@tencent.com>

## Background

There are two CustomResourceDefinitions `JupyterNotebook` and `JupyterGateway` for users. They are used to create Notebook instances and [Jupyter Enterprise Gateways](https://github.com/jupyter/enterprise_gateway). But it is still hard to configure the kernels.

## Motivation

When we configure kernels with Jupyter Enterprise Gateway, we need to provide the kernel json format configburation files. Currently, these configurations is packaged with the gateway image. The kernel specification is shown here:

```
kernel.json
logo-64x64.png
scripts/
  kernel-pod.yaml.j2
  launch_kubernetes.py
```

The scripts directory is copied from the source code of enterprise gateway, and the kernel.json looks like:

```json
{
  "language": "python",
  "display_name": "Python on Kubernetes with Tensorflow",
  "metadata": {
    "process_proxy": {
      "class_name": "enterprise_gateway.services.processproxies.k8s.KubernetesProcessProxy",
      "config": {
        "image_name": "elyra/kernel-tf-py:VERSION"
      }
    }
  },
  "env": {
  },
  "argv": [
    "python",
    "/usr/local/share/jupyter/kernels/python_tf_kubernetes/scripts/launch_kubernetes.py",
    "--RemoteProcessProxy.kernel-id",
    "{kernel_id}",
    "--RemoteProcessProxy.response-address",
    "{response_address}"
  ]
}
```

It is hand-written by enterprise-gateway maintainers. All kernel specs are copied into the enterprise-gateway docker image:

```dockerfile
ADD jupyter_enterprise_gateway_kernelspecs*.tar.gz /usr/local/share/jupyter/kernels/
```

Thus it is hard to update them on the fly.

## Goals

This proposal is to allow users to:

- Configure kernels on the fly
- Manage kernels on Kubernetes easily

## Non-goals

This proposal is not to:

- Make the design the default implementaion, which means that the current design and implementation will not change.

## Implementation

The design proposal involves four main changes to elastic jupyter operator and enterprise gateway:

- Kernel CRD (new): It is used to manage kernels on Kubernetes
- JupyterGateway CRD: Changes is made to support updating kernels on the fly
- Kubeflow Kernel Launcher (new): It is used to launch Kernel CRD inside the gateway
- Kubeflow Process Proxy (new): It is used to manage kernels in enterprise gateway

### Kernel CRD

```yaml
spec:
  ID: "{{ kernel_id }}"
  restartPolicy: Never
  serviceAccountName: "{{ kernel_service_account_name }}"
  securityContext: ...
  environments:
    respondAddress: "{{ eg_response_address }}"
    language: "{{ kernel_language }}"
```

### KernelSpec CRD

cluster scoped.

```yaml
spec:
  language: Python
  displayName: "Python on Kubernetes with Tensorflow"
  image: elyra/kernel-tf-py:VERSION
  envs: ...
  command: 
  - "python",
  - "/usr/local/share/jupyter/scripts/launch_kubernetes.py",
  - "--RemoteProcessProxy.kernel-id",
  - "{kernel_id}",
  - "--RemoteProcessProxy.response-address",
  - "{response_address}"
```

- Create configmap for default specs and scripts when initialization
- Create configmap for new specs

### JupyterGateway CRD

```yaml
spec:
  kernels:
  - python
  - r
  - dask
  ...
```

### Kernel Launcher

### KubeflowProcessProxy
