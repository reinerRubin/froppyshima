package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type (
	Config struct {
		BoltDB BoltDB
		Server Server
	}

	BoltDB struct {
		Path string `env:"BOLT_DB_PATH" default:"/tmp/froppyshimaDB"`
	}

	Server struct {
		Port int `env:"PORT" default:"8080"`
	}
)

// New returns an "env" config
func New() (*Config, error) {
	cfg := &Config{}
	if err := initCfgWithEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// I am aware about the viper and other libraries and probably would use them in real projects
// but here I wanted to try .reflect library
func initCfgWithEnv(c interface{}) error {
	cPtrValue := reflect.ValueOf(c)
	if cPtrValue.Kind() != reflect.Ptr || cPtrValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("initWithEnv: value must be struct ptr")
	}
	cValue := cPtrValue.Elem()

	nameValueMap := make(map[string]interface{}, cValue.NumField())
	for i := 0; i < cValue.NumField(); i++ {
		var (
			fieldMeta  = cValue.Type().Field(i)
			fieldValue = cValue.Field(i)
			fieldKind  = fieldMeta.Type.Kind()
		)

		if fieldKind == reflect.Struct {
			if err := initCfgWithEnv(fieldValue.Addr().Interface()); err != nil {
				return err
			}

			continue
		}
		if fieldKind == reflect.Ptr && fieldValue.Type().Elem().Kind() == reflect.Struct {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))

			if err := initCfgWithEnv(fieldValue.Interface()); err != nil {
				return err
			}

			continue
		}

		if env := fieldMeta.Tag.Get("env"); env != "" {
			value := ""

			defaultValue, defaultPresent := fieldMeta.Tag.Lookup("default")
			if defaultPresent {
				value = defaultValue
			}

			envValue, found := os.LookupEnv(env)
			if found {
				value = envValue
			}

			if !defaultPresent && !found {
				return fmt.Errorf("ENV.%s is not found", env)
			}

			nameValueMap[fieldMeta.Name] = value
		}
	}

	return mapstructure.WeakDecode(nameValueMap, cPtrValue.Interface())
}
