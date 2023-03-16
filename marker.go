package marker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	netboxClient "github.com/netbox-community/go-netbox/netbox/client"
	"github.com/netbox-community/go-netbox/netbox/client/ipam"
	"github.com/rs/zerolog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

// Marker marks nodes by label
type Marker struct {
	k8sClient    *kubernetes.Clientset
	netboxClient *netboxClient.NetBoxAPI
	config       *Config
	logger       *zerolog.Logger
}

// Run runs marker loop
func (m *Marker) Run(ctx context.Context) error {
	ticker := time.NewTicker(
		time.Duration(m.config.SyncInterval) * time.Second,
	)
	defer ticker.Stop()
	m.logger.Info().
		Str("Settings", fmt.Sprintf("%s", m.config)).
		Msg("Marker has started")
	m.markNodesByDC(ctx)
	for {
		select {
		case <-ctx.Done():
			m.logger.Warn().
				Msg("Context was canceled, waiting for graceful shutdown")
			// Graceful shutdown
			return nil

		case <-ticker.C:
			m.markNodesByDC(ctx)
		}
	}
}

func (m *Marker) markNodesByDC(ctx context.Context) (err error) {
	mm, err := m.makeMarkerMap(ctx)
	if err != nil {
		m.logger.Error().
			Err(err).
			Msg("Error occurred while making the markerMap")
	}
	m.logger.Debug().
		Str("marker_map", fmt.Sprint(mm)).
		Msg("markerMap was made")
	err = m.setMarks(ctx, mm)
	if err != nil {
		m.logger.Error().
			Err(err).
			Msg("Error occurred while setting the markerMap")
	} else {
		m.logger.Info().
			Msg("Nodes were labeled")
	}
	return
}

func (m *Marker) setMarks(ctx context.Context, markerMap map[string]string) (err error) {
	for node, dataCenter := range markerMap {
		labels := map[string]string{
			m.config.DCLabel: dataCenter,
		}
		err = m.labelNode(ctx, node, labels)
		if err != nil {
			return
		}
	}
	return
}

func (m *Marker) makeMarkerMap(ctx context.Context) (map[string]string, error) {
	markerMap := make(map[string]string)
	nodes, err := m.getNodes(ctx)
	if err != nil {
		return markerMap, err
	}
	for nodeName, nodeAddress := range nodes {
		location, err := m.getLocation(nodeAddress)
		if err != nil {
			return markerMap, err
		}
		markerMap[nodeName] = location
	}
	return markerMap, err
}

func (m *Marker) getNodes(ctx context.Context) (nodeMap map[string]string, err error) {
	nodeMap = make(map[string]string)
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
	location = m.config.DefaultLocation
	status := "active"
	// TODO: Check if nodeAddress is IP
	params := ipam.IpamPrefixesListParams{
		Q: &nodeAddress,
	}
	params.WithTimeout(time.Second * 5)
	prefixList, err := m.netboxClient.Ipam.IpamPrefixesList(
		&params,
		nil,
	)
	if err != nil {
		return
	}

	for _, prefix := range prefixList.Payload.Results {
		if *prefix.Status.Value == status {
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

// GetNewMarker returns new exemplar of marker
func GetNewMarker(k8sClientset *kubernetes.Clientset, netboxClient *netboxClient.NetBoxAPI,
	config *Config, logger *zerolog.Logger) *Marker {
	marker := Marker{
		k8sClientset,
		netboxClient,
		config,
		logger,
	}
	return &marker
}
