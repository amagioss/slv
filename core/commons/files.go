package commons

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func WriteToYAML(filePath, notice string, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if notice != "" {
		bytes = append([]byte(notice), bytes...)
	}
	bytes = append([]byte(slvYamlNotice), bytes...)
	if err == nil {
		if err = WriteToFile(filePath, bytes); err != nil {
			return err
		}
	}
	return err
}

func ReadChildFromYAML(filePath, nodePath string, out interface{}) error {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into a map
	var objMap map[string]interface{}
	err = yaml.Unmarshal(data, &objMap)
	if err != nil {
		return err
	}

	// Split the node path and traverse the map
	nodes := strings.Split(nodePath, ".")
	for _, node := range nodes {
		if index, err := strconv.Atoi(node); err == nil {
			// If the node is an integer, treat it as an array index
			objArray, _ := objMap[nodes[0]].([]interface{})
			objMap, _ = objArray[index].(map[string]interface{})
			nodes = nodes[1:]
		} else {
			objMap, _ = objMap[node].(map[string]interface{})
		}
	}
	bytes, err := yaml.Marshal(objMap)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bytes, out)
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
