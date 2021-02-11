package marker

import (
	"context"
	"encoding/json"
	"fmt"

	netboxClient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// Marker marks nodes by label
type Marker struct {
	k8sClient    *kubernetes.Clientset
	netboxClient *netboxClient.NetBoxAPI
	config       *Config
}

func (m *Marker) setMarks(ctx context.Context) (markerMap map[string]string, err error) {
	for node, dataCenter := range markerMap {
		labels := map[string]string{
			"DC": dataCenter,
		}
		m.labelNode(ctx, node, labels)
	}
}

func (m *Marker) makeMarkerMap(ctx context.Context) (markerMap map[string]string, err error) {
	nodes, err := m.getNodes(ctx)
	if err != nil {
		return
	}
	for nodeName, nodeAddress := range nodes {
		location, err := m.getLocation(nodeAddress)
		if err != nil {
			return markerMap, err
		}
		markerMap[nodeName] = location
	}
	return
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
	location = m.config.defaultLocation
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
