## Quickstart

### Simple deployment

You can create a simple Jupyter notebook with all components in one pod, like this:

```yaml
$ cat ./examples/simple-deployments/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterNotebook
metadata:
  name: jupyternotebook-simple
spec:
  auth:
    mode: disable
  template:
    metadata:
      labels:
        notebook: simple
    spec:
      containers:
        - name: notebook
          image: jupyter/base-notebook:python-3.8.6
          command: ["tini", "-g", "--", "start-notebook.sh"]

$ kubectl apply -f ./examples/simple-deployments/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
$ kubectl port-forward deploy/jupyternotebook-simple 8888:8888
```

Then you can open the URL `http://127.0.0.1:8888/` to get the simple Jupyter notebook instance. The deployment follows the architecture below:

<p align="center"><img src="./images/kubeflow.png" width="300"></p>


## Elastic deployment

elastic-jupyter-operator supports running Jupyter kernels in separate pods. In this example, we will create the notebook and gateway. 

```yaml
$ cat ./examples/elastic/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterNotebook
metadata:
  name: jupyternotebook-elastic
spec:
  gateway:
    name: jupytergateway-elastic
    namespace: default
  auth:
    mode: disable

$  cat ./examples/elastic/kubeflow.tkestack.io_v1alpha1_jupytergateway.yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterGateway
metadata:
  name: jupytergateway-elastic
spec:
  cullIdleTimeout: 3600
  image: ccr.ccs.tencentyun.com/kubeflow-oteam/enterprise-gateway:2.5.0

$ kubectl apply -f ./examples/elastic/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
$ kubectl apply -f ./examples/elastic/kubeflow.tkestack.io_v1alpha1_jupytergateway.yaml
$ kubectl port-forward deploy/jupyternotebook-elastic 8888:8888
```

When users run the code in the browser, there will be a new kernel pod created in the cluster.

```
NAME                                          READY   STATUS    RESTARTS   AGE
jovyan-219cfd49-89ad-428c-8e0d-3e61e15d79a7   1/1     Running   0          170m
jupytergateway-elastic-868d8f465c-8mg44       1/1     Running   0          3h
jupyternotebook-elastic-787d94bb4b-xdwnc      1/1     Running   0          3h10m
```

### Elastic deployment with custom kernel

If you want to custom the kernel deployment, for example. you want to update the resource requirements of the python kernel or use different images for the kernel, you can deploy the jupyter notebooks and gateways with custom kernels.

First, you need to create the JupyterKernelSpec CR, which is used to generate the [Jupyter kernelspec](https://jupyter-client.readthedocs.io/en/stable/kernels.html).

```yaml
$ cat examples/elastic-with-custom-kernels/kubeflow.tkestack.io_v1alpha1_jupyterkernelspec.yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterKernelSpec
metadata:
  name: python-kubernetes
spec:
  language: Python
  displayName: "Python on Kubernetes as a JupyterKernelSpec"
  image: ccr.ccs.tencentyun.com/kubeflow-oteam/jupyter-kernel-py:2.5.0
  className: enterprise_gateway.services.processproxies.kubeflow.KubeflowProcessProxy
  # Use the template defined in JupyterKernelTemplate CR.
  template:
    namespace: default
    name: jupyterkerneltemplate-elastic-with-custom-kernels
  command: 
  # Use the default scripts to launch the kernel.
  - "kubeflow-launcher"
  - "--verbose"
  - "--RemoteProcessProxy.kernel-id"
  - "{kernel_id}"
  - "--RemoteProcessProxy.port-range"
  - "{port_range}"
  - "--RemoteProcessProxy.response-address"
  - "{response_address}"

$ cat examples/elastic-with-custom-kernels/kubeflow.tkestack.io_v1alpha1_jupyterkerneltemplate.yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterKernelTemplate
metadata:
  name: jupyterkerneltemplate-elastic-with-custom-kernels
spec:
  template:
    metadata: 
      app: enterprise-gateway
      component: kernel
    spec:
      restartPolicy: Always
      containers:
        - name: kernel

$ kubectl apply -f  ./examples/elastic-with-custom-kernels/kubeflow.tkestack.io_v1alpha1_jupyterkernelspec.yaml
$ kubectl apply -f ./examples/elastic-with-custom-kernels/kubeflow.tkestack.io_v1alpha1_jupyterkerneltemplate.yaml
```

There will be a configmap created with the given CR, and it will be mounted into the gateway.

```yaml
$ cat examples/elastic-with-custom-kernels/kubeflow.tkestack.io_v1alpha1_jupytergateway.yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterGateway
metadata:
  name: jupytergateway-elastic-with-custom-kernels
spec:
  cullIdleTimeout: 10
  cullInterval: 10
  logLevel: DEBUG
  image: ccr.ccs.tencentyun.com/kubeflow-oteam/enterprise-gateway:dev
  # Use the kernel which is defined in JupyterKernelSpec CR.
  kernels: 
  - python-kubernetes

$ kubectl apply -f ./examples/elastic/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
$ kubectl apply -f ./examples/elastic/kubeflow.tkestack.io_v1alpha1_jupytergateway.yaml
$ kubectl port-forward deploy/jupyternotebook-elastic-with-custom-kernels 8888:8888
```
