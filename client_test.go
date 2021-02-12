package marker

import (
	"context"
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetNewK8sClient(t *testing.T) {
	client, err := GetNewK8sClient()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{}))
}
