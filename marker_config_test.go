package marker

import (
	"os"
	"strconv"
	"testing"
)

var testEnvs map[string]string

func init() {
	testEnvs = map[string]string{
		"DC_LABEL":              "DC",
		"DEFAULT_DC":            "Unknown",
		"SYNC_INTERVAL_SECONDS": "10",
	}
	for k, v := range testEnvs {
		os.Setenv(k, v)
	}
}

func TestGetNewConfig(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}
	if config.DCLabel != testEnvs["DC_LABEL"] {
		t.Fail()
		t.Logf(
			"%s failed. Reason: %s",
			t.Name(),
			"Value of DC_LABEL was incorrect parsed",
		)
	}
	if config.DefaultLocation != testEnvs["DEFAULT_DC"] {
		t.Fail()
		t.Logf(
			"%s failed. Reason: %s",
			t.Name(),
			"Value of DEFAULT_DC was incorrect parsed",
		)
	}
	seconds, err := strconv.Atoi(testEnvs["SYNC_INTERVAL_SECONDS"])
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	if config.SyncInterval != seconds {
		t.Fail()
		t.Logf(
			"%s failed. Reason: %s",
			t.Name(),
			"Value of SYNC_INTERVAL_SECONDS was incorrect parsed")
	}

}
