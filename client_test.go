package marker

import (
	"testing"
)

func TestGetNewK8sClient(t *testing.T) {
	_, err := GetNewK8sClient()
	if err != nil {
		t.Error(err)
	}
}
