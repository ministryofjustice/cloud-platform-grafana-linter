package cluster

import (
	"context"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func SearchNamespacesForConfigMaps(clientset *kubernetes.Clientset) ([]*v1.ConfigMapList, error) {
	var configMaps []*v1.ConfigMapList
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, namespace := range namespaces.Items {
		configMap, err := clientset.CoreV1().ConfigMaps(namespace.Name).List(context.Background(), meta_v1.ListOptions{
			LabelSelector: "grafana_dashboard",
		})
		if err != nil {
			return nil, err
		}
		configMaps = append(configMaps, configMap)
	}
	return configMaps, nil
}
