# elastic-jupyter-operator

Elastic Jupyter Notebooks on Kubernetes

## Motivation

Jupyter is a free, open-source, interactive web tool known as a computational notebook, which researchers can use to combine software code, computational output, explanatory text, and multimedia resources in a single document. 

For data scientists and machine learning engineers, Jupyter has emerged as a de facto standard. At the same time, there has been growing criticism that the way notebooks are being used leads to low resource utilization.

GPU and other hardware resources will be bound to the specified notebooks even if the data scientists do not need them currently. This project proposes some Kubernetes CRDs to solve these problems.

## Introduction

elastic-jupyter-operator provides elastic Jupyter notebook services with these features:

- Provide users the out-of-box Jupyter notebooks on Kubernetes.
- Autoscale Jupyter kernels when the kernels are not used within the given time frame to increase the resource utilization.
- Customize the kernel configuration in runtime without restarting the notebook.

<p align="center"><img src="docs/images/elastic.png" width="300"></p>
Figure 1. elastic-jupyter-operator

<p align="center"><img src="docs/images/jupyter.png" width="300"></p>
Figure 2. Other Jupyter on Kubernetes solutions

## Deploy

```bash
kubectl apply -f ./hack/enterprise_gateway/prepare.yaml
make deploy
```

## Quick start

You can follow the [quick start](./docs/quick-start.md) to create the notebook server and kernel in Kubernetes like this:

```bash
NAME                                                           READY   STATUS    RESTARTS   AGE
jovyan-fd191444-b08c-4668-ba4e-3748a54a0ac1-5789574d66-tb5cm   1/1     Running   0          146m
jupytergateway-sample-858bbc8d5c-xds44                         1/1     Running   0          3h46m
jupyternotebook-sample-5bf7d9d9fb-pdv9b                        1/1     Running   10         77d
```

There are three pods running in the demo:

- `jupyternotebook-sample-5bf7d9d9fb-pdv9b` is the notebook server
- `jupytergateway-sample-858bbc8d5c-xds44` is the jupyter gatrway to support remote kernels
- `jovyan-fd191444-b08c-4668-ba4e-3748a54a0ac1-5789574d66-tb5cm` is the remote kernel

The kernel will be deleted if the notebook does not use it in 10 mins. And it will be recreated if there is any new run in the notebook.

## Design

Please refer to [design doc](docs/design.md)

## API Documentation

Please refer to [API doc](docs/api/generated.asciidoc)

## Special Thanks

- [jupyter/enterprise_gateway](https://github.com/jupyter/enterprise_gateway) which implements the logic to run kernels remotely
