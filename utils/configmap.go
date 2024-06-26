package utils

import (
	"context"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NamespaceConfigMaps(clientset *kubernetes.Clientset) ([]*v1.ConfigMapList, error) {
	var configMaps []*v1.ConfigMapList

	namespaces, err := getNamespaces(clientset)
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces {
		configMap, err := clientset.CoreV1().ConfigMaps(namespace).List(context.Background(), meta_v1.ListOptions{
			LabelSelector: "grafana_dashboard",
		})
		if err != nil {
			return nil, err
		}
		configMaps = append(configMaps, configMap)
	}
	return configMaps, nil
}

func getNamespaces(clientset *kubernetes.Clientset) ([]string, error) {
	var namespaces []string

	ns, err := clientset.CoreV1().Namespaces().List(ctx, meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, n := range ns.Items {
		namespaces = append(namespaces, n.Name)
	}

	return namespaces, nil
}
