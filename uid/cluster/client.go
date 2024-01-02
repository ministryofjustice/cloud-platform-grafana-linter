package cluster

import (
	"fmt"
	"os"

	"github.com/gophercloud/utils/env"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfigPath = env.Getenv("KUBE_CONFIG")
)

func ClientSet() *kubernetes.Clientset {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	return clientSet
}
