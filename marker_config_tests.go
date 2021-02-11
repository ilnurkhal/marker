package marker

import (
	"os"
	"strconv"
	"testing"
	"time"
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
	config, err := GetNewConfig()
	if err != nil {
		t.Fail()
		t.Logf(
			"%s failed. Error: %s",
			t.Name(),
			err,
		)
	}
	if config.dcLabel != testEnvs["DC_LABEL"] {
		t.Fail()
		t.Logf(
			"%s failed. Reason: %s",
			t.Name(),
			"Value of DC_LABEL was incorrect parsed",
		)
	}
	if config.defaultLocation != testEnvs["DEFAULT_DC"] {
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
	if config.syncInterval != time.Duration(seconds)*time.Second {
		t.Fail()
		t.Logf(
			"%s failed. Reason: %s",
			t.Name(),
			"Value of SYNC_INTERVAL_SECONDS was incorrect parsed")
	}

}
