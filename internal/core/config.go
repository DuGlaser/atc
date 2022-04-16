package core

import (
	"fmt"
	"reflect"
)

type Config struct {
	RunCmd   string
	BuildCmd string
	Lang     string
	FileName string
	Template string
}

func (c *Config) Validate() error {
	keys := []string{"RunCmd", "Lang", "FileName"}

	for _, key := range keys {
		if err := c.errorEmptyValue(key); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) errorEmptyValue(key string) error {
	v := reflect.ValueOf(c)
	f := reflect.Indirect(v).FieldByName(key)

	if !f.IsValid() {
		return fmt.Errorf("Config.%s is invalid key.", key)
	}

	if f.String() == "" {
		return fmt.Errorf("Config.%s is empty.", key)
	}

	return nil
}
