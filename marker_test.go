package marker

import (
	"context"
	"fmt"
	"os"
	"testing"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/netbox-community/go-netbox/netbox/client"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	marker      *Marker
	testAddress = os.Getenv("TEST_NODE_ADDRESS")
	netboxHost  = os.Getenv("NETBOX_HOST")
	netboxToken = os.Getenv("NETBOX_TOKEN")
)

const (
	nodeName = "minikube"
)

func init() {
	testEnvs = map[string]string{
		"DC_LABEL":              "DC",
		"DEFAULT_DC":            "Unknown",
		"SYNC_INTERVAL_SECONDS": "10",
	}
	for k, v := range testEnvs {
		os.Setenv(k, v)
	}
	config, err := GetNewConfig()
	if err != nil {
		fmt.Println("Error:", err)
	}
	clientSet, err := GetNewK8sClient()
	if err != nil {
		fmt.Println("Error:", err)
	}
	transport := httptransport.New(
		netboxHost,
		client.DefaultBasePath,
		[]string{"https"})
	transport.DefaultAuthentication = httptransport.APIKeyAuth(
		"Authorization",
		"header",
		fmt.Sprintf("Token %s", netboxToken))
	netBoxClient := client.New(transport, nil)

	if err != nil {
		fmt.Println("Error", err)
	}

	marker = &Marker{
		clientSet,
		netBoxClient,
		&config,
	}

}

func TestGetNodes(t *testing.T) {
	nodes, err := marker.getNodes(context.TODO())
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}
	if len(nodes) < 1 {
		t.Fail()
		t.Log("There are no Nodes")
	}

}

func TestGetLocation(t *testing.T) {
	location, err := marker.getLocation(testAddress)
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}
	if len(location) < 1 {
		t.Fail()
		t.Log("Empty location")
	}
}

func TestLabelNode(t *testing.T) {
	err := marker.labelNode(
		context.TODO(),
		nodeName,
		map[string]string{
			"MARKER": "TEST",
		},
	)
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}
	testNode, err := marker.k8sClient.CoreV1().Nodes().Get(
		context.TODO(),
		nodeName,
		v1.GetOptions{},
	)
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}
	labels := testNode.GetLabels()
	if value, ok := labels["MARKER"]; ok {
		if value != "TEST" {
			t.Fail()
			t.Log("Incorrect value of test label")
		}
	} else {
		t.Fail()
		t.Log("There is no test label")
	}
	delete(labels, "MARKER")
	testNode.SetLabels(labels)
	_, err = marker.k8sClient.CoreV1().Nodes().Update(context.TODO(), testNode, v1.UpdateOptions{})
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}
}
