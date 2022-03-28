package util

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// kubectlGetNodesJSON parses the output of the "kubectl get nodes -o json" command.
type kubectlGetNodesJSON struct {
	Items []struct {
		Status struct {
			Addresses []struct {
				Address string `json:"address"`
				Type    string `json:"type"`
			} `json:"addresses"`
		} `json:"status"`
	} `json:"items"`
}

func parseControlPlaneNodeIPs(jsonOutput []byte) ([]string, error) {
	var response kubectlGetNodesJSON
	if err := json.Unmarshal(jsonOutput, &response); err != nil {
		return nil, fmt.Errorf("failed to parse kubectl command output: %w", err)
	}

	nodes := make([]string, 0, len(response.Items))
	for _, item := range response.Items {
		for _, address := range item.Status.Addresses {
			if address.Type == "InternalIP" {
				nodes = append(nodes, address.Address)
			}
		}
	}

	return nodes, nil
}

// ListControlPlaneNodeIPs returns the internal IPs of the control plane nodes of the MicroK8s cluster.
func ListControlPlaneNodeIPs(ctx context.Context) ([]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", SnapDataPath("credentials", "client.config"))
	if err != nil {
		return nil, fmt.Errorf("failed to read load kubeconfig: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize kubernetes client: %w", err)
	}

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{
		LabelSelector: "node.kubernetes.io/microk8s-controlplane=microk8s-controlplane",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	addresses := make([]string, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeInternalIP {
				addresses = append(addresses, address.Address)
			}
		}
	}

	return addresses, nil
}
