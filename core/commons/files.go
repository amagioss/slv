package commons

import (
	"os"

	"gopkg.in/yaml.v3"
)

func WriteToYAML(filePath, notice string, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	bytes = append([]byte(notice), bytes...)
	bytes = append([]byte(slvYamlNotice), bytes...)
	if err == nil {
		if err = WriteToFile(filePath, bytes); err != nil {
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
