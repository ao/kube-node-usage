# kube-node-usage
Tell me about the usage of all the nodes in my active Kubernetes cluster

## Requirements

If you get an error similar to this:

```bash
Error getting metrics for node ip-192-168-123-456.us-west-2.compute.internal: the server could not find the requested resource (get nodes.metrics.k8s.io ip-192-168-123-456.us-west-2.compute.internal)
+------+-----------+--------------+
| NODE | CPU USAGE | MEMORY USAGE |
+------+-----------+--------------+
+------+-----------+--------------+
```

Then you need to make sure that the Metrics server is installed as follows:

To ensure the Metrics Server is installed and running, you can follow these steps:

1. Install the Metrics Server:

If the Metrics Server is not installed, you can deploy it using the following command:

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

2. Verify Metrics Server Deployment:

Check if the Metrics Server is running properly:

```bash
kubectl get deployment metrics-server -n kube-system
```

Ensure that the Metrics Server pods are running without any issues.

3. Check Metrics API Availability:

Verify that the Metrics API is available:

```bash
kubectl get --raw "/apis/metrics.k8s.io/v1beta1/nodes" | jq .
```
