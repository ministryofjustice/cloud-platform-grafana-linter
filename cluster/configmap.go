package cluster

import (
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Dashboard struct {
	Title,
	Namespace string
	UID string `json:"uid"`
}

func SearchConfigMapsData(clientset *kubernetes.Clientset, configMap *v1.ConfigMapList, uid string) string {
	var dashboard Dashboard
	var d string

	for _, items := range configMap.Items {
		for _, value := range items.Data {
			Data := []byte(value)
			json.Unmarshal(Data, &dashboard)
			dashboard.Title = items.Name
			dashboard.Namespace = items.Namespace

			if dashboard.UID != "" {
				b := uid == dashboard.UID
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
