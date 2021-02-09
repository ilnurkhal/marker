package marker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	netboxClient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type Marker struct {
	k8sClient    *kubernetes.Clientset
	netboxClient *netboxClient.NetBoxAPI
	syncInterval time.Duration
}

func (m *Marker) makeMarkerMap() (markerMap map[string]string, err error) {

}

func (m *Marker) getNodes(ctx context.Context) (nodeMap map[string]string, err error) {
	nodes, err := m.k8sClient.CoreV1().Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return
	}
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == "InternalIP" {
				nodeMap[node.GetName()] = address.Address
			}
		}
	}
	return
}

func (m *Marker) getLocation(nodeAddress string) (location string, err error) {
	location = "Unknown" // Default location
	status := "active"
	// TODO: Check if nodeAddress is IP
	params := ipam.IpamPrefixesListParams{
		Q: &nodeAddress,
	}
	prefixList, err := m.netboxClient.Ipam.IpamPrefixesList(
		&params,
		nil,
	)
	if err != nil {
		return
	}
	for _, prefix := range prefixList.Payload.Results {
		if prefix.Status.Value == &status {
			location = *prefix.Site.Name
		}
	}
	return
}

func (m *Marker) labelNode(ctx context.Context, nodeName string, labels map[string]string) (err error) {
	node, err := m.k8sClient.CoreV1().Nodes().Get(ctx, nodeName, v1.GetOptions{})
	if err != nil {
		return
	}
	nodeLabels := node.GetLabels()
	for key, value := range labels {
		nodeLabels[key] = value
	}
	json, err := json.Marshal(nodeLabels)
	if err != nil {
		return
	}
	patch := []byte(
		fmt.Sprintf(
			`{"metadata":{"labels":%s}}`,
			json,
		),
	)
	_, err = m.k8sClient.CoreV1().Nodes().Patch(
		ctx,
		nodeName,
		types.StrategicMergePatchType,
		patch,
		v1.PatchOptions{},
	)
	return

}
