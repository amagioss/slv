package commons

import (
	"os"

	"gopkg.in/yaml.v3"
)

func WriteToYAML(filePath string, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err == nil {
		if err = os.WriteFile(filePath, bytes, 0644); err != nil {
			return err
		}
	}
	return err
}

func ReadFromYAML(filePath string, out interface{}) error {
	bytes, err := os.ReadFile(filePath)
	if err == nil {
		if err = yaml.Unmarshal(bytes, out); err != nil {
			return err
		}
	}
	return err
}
