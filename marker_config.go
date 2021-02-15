package marker

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// Config for marker
type Config struct {
	DCLabel         string `env:"DC_LABEL"`
	DefaultLocation string `env:"DEFAULT_DC"`
	SyncInterval    int    `env:"SYNC_INTERVAL_SECONDS"`
}

// ParseFromEnv parses configs from env
func (conf *Config) parseFromEnv() error {
	fields := reflect.ValueOf(conf).Elem()
	for i := 0; i < fields.NumField(); i++ {
		fieldInfo := fields.Type().Field(i)
		env := fieldInfo.Tag.Get("env")
		if value, ok := os.LookupEnv(env); ok {
			switch fieldInfo.Type.Kind() {
			case reflect.String:
				fields.Field(i).SetString(value)
			case reflect.Int: // Case for syncInterval
				syncInterval, err := strconv.Atoi(value)
				if err == nil {
					fields.Field(i).SetInt(int64(syncInterval))
				} else {
					fields.Field(i).SetInt(60)
				}
			}
		} else {
			return fmt.Errorf("The environmet %s isn't set", env)
		}
	}
	return nil
}

func (conf *Config) String() string {
	return fmt.Sprintf(
		"DCLabel: %s DefaultLocation: %s, SyncInterval: %v",
		conf.DCLabel,
		conf.DefaultLocation,
		conf.SyncInterval,
	)
}

// GetNewConfig returns an exemplar of config
func GetNewConfig() (config Config, err error) {
	err = config.parseFromEnv()
	return
}
