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
- KernelSpec CRD (new): It is used to configure the kernel specs in runtime
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

> The primary vehicle for indicating a given kernel should be handled in a different manner is the kernel specification, otherwise known as the kernel spec. Enterprise Gateway leverages the natively extensible metadata stanza to introduce a new stanza named process_proxy.
>
> The process_proxy stanza identifies the class that provides the kernel’s process abstraction (while allowing for future extensions). This class then provides the kernel’s lifecycle management operations relative to the managed resource or functional equivalent.
>
> Here’s an example of a kernel specification that uses the DistributedProcessProxy class for its abstraction:
>
```json
{
  "language": "scala",
  "display_name": "Spark - Scala (YARN Client Mode)",
  "metadata": {
    "process_proxy": {
      "class_name": "enterprise_gateway.services.processproxies.distributed.DistributedProcessProxy"
    }
  },
  "env": {
    "SPARK_HOME": "/usr/hdp/current/spark2-client",
    "__TOREE_SPARK_OPTS__": "--master yarn --deploy-mode client --name ${KERNEL_ID:-ERROR__NO__KERNEL_ID}",
    "__TOREE_OPTS__": "",
    "LAUNCH_OPTS": "",
    "DEFAULT_INTERPRETER": "Scala"
  },
  "argv": [
    "/usr/local/share/jupyter/kernels/spark_scala_yarn_client/bin/run.sh",
    "--RemoteProcessProxy.kernel-id",
    "{kernel_id}",
    "--RemoteProcessProxy.response-address",
    "{response_address}",
    "--RemoteProcessProxy.public-key",
    "{public_key}"
  ]
}
```

The kernel specifications are placed in the docker image at build time, which is not easy to maintain on the fly. The kernelspec CRD is defined to support dynamic update. The CRD specification looks like this:

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

When a JupyterKernelSpec CR is created, we will create the corresponding configmap. And the configmap will be used as a mount volume in the gateway.

### JupyterGateway CRD

The specification generation logic needs to be changed to support the new JupyterKernelSpec CRD.

```yaml
spec:
  kernels:
  - python
  - r
  - dask
  ...
```

When `kernels` are defined in the spec, we should get the jupyter kernelspec CRs from the kubernetes api server, then mount the configmaps as volumes into the gateway container.

### Kernel Launcher

### KubeflowProcessProxy

## Reference

- [Jupyter Enterprise Gateway System Architecture](https://jupyter-enterprise-gateway.readthedocs.io/en/latest/system-architecture.html)
