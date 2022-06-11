package core

import (
	"bytes"
	"fmt"
	"html/template"
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

func (c *Config) GenerateCmd(path, fileName string) error {
	rcTmpl, err := template.New("runCmd").Parse(c.RunCmd)
	if err != nil {
		return err
	}

	var rc bytes.Buffer
	err = rcTmpl.Execute(&rc, map[string]interface{}{
		"file": path,
		"dir":  path[0 : len(path)-len(fileName)-1],
	})

	c.RunCmd = rc.String()

	if c.BuildCmd != "" {
		bcTmpl, err := template.New("buildCmd").Parse(c.BuildCmd)
		if err != nil {
			return err
		}

		var bc bytes.Buffer
		err = bcTmpl.Execute(&bc, map[string]interface{}{
			"file": path,
			"dir":  path[0 : len(path)-len(fileName)-1],
		})
		if err != nil {
			return err
		}

		c.BuildCmd = bc.String()
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
