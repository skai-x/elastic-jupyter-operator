## Quick start

### Simple deployment

You can create a simple Jupyter notebook with all components in one pod, like this:

```
$ cat ./examples/simple-deployments/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
apiVersion: kubeflow.tkestack.io/v1alpha1
kind: JupyterNotebook
metadata:
  name: jupyternotebook-simple
spec:
  template:
    metadata:
      labels:
        notebook: simple
    spec:
      containers:
        - name: notebook
          image: jupyter/base-notebook:python-3.8.6

$ kubectl apply -f ./examples/simple-deployments/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
```

## 使用

首先，创建一个 Jupyter Gateway CR：

```bash
kubectl apply -f ./config/samples/kubeflow.tkestack.io_v1alpha1_jupytergateway.yaml
```

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

```bash
kubectl apply -f ./config/samples/kubeflow.tkestack.io_v1alpha1_jupyternotebook.yaml
```

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