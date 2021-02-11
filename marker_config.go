package marker

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Config for marker
type Config struct {
	dcLabel         string        `env:"DC_LABEL"`
	defaultLocation string        `env:"DEFAULT_DC"`
	syncInterval    time.Duration `env:"SYNC_INTERVAL_SECONDS"`
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
				duration, err := strconv.Atoi(value)
				durSeconds := time.Second * time.Duration(duration)
				if err == nil {
					fields.Field(i).SetInt(int64(durSeconds))
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

// GetNewConfig returns an exemplar of config
func GetNewConfig() (config Config, err error) {
	err = config.parseFromEnv()
	return
}
