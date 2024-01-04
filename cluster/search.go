package cluster

import (
	"context"
	"encoding/json"
	"fmt"

	utils "github.com/ministryofjustice/cloud-platform-grafana-linter/utils"

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

func SearchConfigMapsData(clientset *kubernetes.Clientset, configMap *v1.ConfigMapList, uid string) string {
	var dashboard utils.Dashboard
	var d string

	for _, items := range configMap.Items {
		for _, value := range items.Data {
			Data := []byte(value)
			json.Unmarshal(Data, &dashboard)
			dashboard.Title = items.Name
			dashboard.Namespace = items.Namespace

			if dashboard.UID != "" {
				b := DuplicateUID(uid, dashboard.UID)
				if !b {
					continue
				}
				if b {
					d = fmt.Sprintf("Found duplicate UID %v in namespace %v, title %v\n", dashboard.UID, dashboard.Namespace, dashboard.Title)
					continue
				}
			}
		}
	}
	return d
}
