package commons

import (
	"os"

	"gopkg.in/yaml.v3"
)

func WriteToYAML(filePath string, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	bytes = append([]byte(yamlNotice), bytes...)
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

func WriteToFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}

func ReadFromFile(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	return content, err
}