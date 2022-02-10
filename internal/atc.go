package internal

import (
	"fmt"
	"reflect"
)

type Config struct {
	Cmd      string
	Lang     string
	FileName string
	Template string
}

func (c *Config) Validate() error {
	if err := c.errorEmptyValue("Cmd"); err != nil {
		return err
	}

	if err := c.errorEmptyValue("Lang"); err != nil {
		return err
	}

	if err := c.errorEmptyValue("FileName"); err != nil {
		return err
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
