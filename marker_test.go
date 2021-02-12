package marker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/netbox-community/go-netbox/netbox"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var marker *Marker

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
	netBoxClient := netbox.NewNetboxWithAPIKey(
		"your.netbox.host:8000",
		"your_netbox_token",
	)
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

	jsonLabels, err := json.Marshal(labels)
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}

	patchJson, err := json.Marshal(
		map[string]string{
			"op":    "repalce",
			"path":  "metadata/labels",
			"value": fmt.Sprint(jsonLabels),
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

	_, err = marker.k8sClient.CoreV1().Nodes().Patch(
		context.TODO(),
		nodeName,
		types.JSONPatchType,
		patchJson,
		v1.PatchOptions{},
	)
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}

}
