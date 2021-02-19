# jupyter-operator: Elastic Jupyter on Kubernetes

Kubernetes 原生的弹性 Jupyter 即服务

## 介绍

Jupyter 已经成为了数据科学家和算法工程师最常用的 IDE 类应用之一。但目前大规模运维 Jupyter Notebook 实例仍然存在许多问题：

- 资源利用率低，Notebook 实例常常会占据 GPU 资源但并没有真正进行计算。这会导致虽然资源已经被分配，但没有得到合理地使用。
- 管理困难，不同的实例需要访问不同的存储、利用不同规格的硬件资源完成数据科学家的数据分析和模型训练任务。这对集群的运维和管理都带来了挑战。
- 缺乏弹性，当集群上混合运行分布式模型训练、模型离线批处理推理和 Notebook 时，难以进行资源的回收再分配，影响高优先级工作负载的运行。

Jupyter operator 主要面向这三个问题，设计了两个 CRD，实现了在 Kubernetes 上的弹性部署。

## 部署

```bash
$ kubectl apply -f ./hack/enterprise_gateway/prepare.yaml
$ make deploy
```

## 架构



## 使用

首先，创建一个 Jupyter Gateway CR：

```yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterGateway
metadata:
  name: jupytergateway-sample
spec:
  cullIdleTimeout: 3600
```

其中 `cullIdleTimeout` 是一个配置项，在 Kernel 空闲指定 `cullIdleTimeout` 秒内，会由 Gateway 回收对应 Kernel 以释放资源。

其次需要创建一个 Jupyter Notebook CR 实例，并且指定对应的 Gateway CR：

```yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterNotebook
metadata:
  name: jupyternotebook-sample
spec:
  gateway:
    name: jupytergateway-sample
    namespace: default
```

集群上所有资源如下所示：

```
NAME                                          READY   STATUS    RESTARTS   AGE
pod/jupytergateway-sample-6d5d97949c-p8bj6    1/1     Running   2          11d
pod/jupyternotebook-sample-5bf7d9d9fb-nq9b8   1/1     Running   2          11d

NAME                            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service/jupytergateway-sample   ClusterIP   10.96.138.111   <none>        8888/TCP   11d
service/kubernetes              ClusterIP   10.96.0.1       <none>        443/TCP    31d

NAME                                     READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/jupytergateway-sample    1/1     1            1           11d
deployment.apps/jupyternotebook-sample   1/1     1            1           11d

NAME                                                DESIRED   CURRENT   READY   AGE
replicaset.apps/jupytergateway-sample-6d5d97949c    1         1         1       11d
replicaset.apps/jupyternotebook-sample-5bf7d9d9fb   1         1         1       11d
```

随后可以通过 NodePort、`kubectl port-forward`、ingress 等方式将 Notebook CR 对外暴露提供服务，这里以 `kubectl port-forward` 为例：

```
kubectl port-forward jupyternotebook-sample-5bf7d9d9fb-nq9b8 8888
```

