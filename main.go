package main

import (
    "context"
    "flag"
    "log"
    "os"
    "path/filepath"
    "strconv"

    v1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/metrics/pkg/client/clientset/versioned"
    "github.com/olekukonko/tablewriter"
)

func main() {
    // Parse command-line flags
    kubeconfig := flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    flag.Parse()

    // Build the Kubernetes client config
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        log.Fatalf("Error building kubeconfig: %s", err)
    }

    // Create the Kubernetes clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating Kubernetes clientset: %s", err)
    }

    // Create the Metrics clientset
    metricsClientset, err := versioned.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating Metrics clientset: %s", err)
    }

    // List all nodes
    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Error listing nodes: %s", err)
    }

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Node", "CPU Usage (m)", "Memory Usage (MiB)"})

    for _, node := range nodes.Items {
        // Get node metrics
        nodeMetrics, err := metricsClientset.MetricsV1beta1().NodeMetricses().Get(context.TODO(), node.Name, metav1.GetOptions{})
        if err != nil {
            log.Printf("Error getting metrics for node %s: %s\n", node.Name, err)
            continue
        }

        cpuUsage := nodeMetrics.Usage[v1.ResourceCPU]
        memoryUsage := nodeMetrics.Usage[v1.ResourceMemory]

        cpuUsageMillicores := float64(cpuUsage.MilliValue())
        memoryUsageMiB := float64(memoryUsage.Value()) / (1024 * 1024)

        table.Append([]string{
            node.Name,
            strconv.FormatFloat(cpuUsageMillicores, 'f', 2, 64) + " m",
            strconv.FormatFloat(memoryUsageMiB, 'f', 2, 64) + " MiB",
        })
    }

    table.Render()
}

