package marker

import (
	"testing"
)

func TestGetNewK8sClient(t *testing.T) {
	_, err := NewK8sClient()
	if err != nil {
		t.Error(err)
	}
}
