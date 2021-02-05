package marker

import (
	netboxClient "github.com/netbox-community/go-netbox/netbox/client"
	"k8s.io/client-go/kubernetes"
)

type Marker struct {
	k8sClient    *kubernetes.Clientset
	netboxClient *netboxClient.NetBoxAPI
}

func (m *Marker) makeMarkerMap() (markerMap map[string]string, err error) {

}

func (m *Marker) getNodes() (nodes []string, err error) {

}

func (m *Marker) getLocation(node string) (location string, err error) {

}

func (m *Marker) labelNode(node string, label map[string]string) (err error) {

}
