package commons

import (
	"os"

	"gopkg.in/yaml.v3"
)

const slvYamlNotice = "# This file is managed by SLV. Please avoid editing this file manually.\n"

func WriteToYAML(filePath string, data any) error {
	bytes, err := yaml.Marshal(data)
	if err == nil {
		bytes = append([]byte(slvYamlNotice), bytes...)
		err = os.WriteFile(filePath, bytes, 0644)
	}
	return err
}

func ReadFromYAML(filePath string, out any) error {
	bytes, err := os.ReadFile(filePath)
	if err == nil {
		err = yaml.Unmarshal(bytes, out)
	}
	return err
}
