package cmdvault

import "gopkg.in/yaml.v3"

type k8Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Data map[string]string `yaml:"data"`
	Type string            `yaml:"type"`
}

func k8sSecretFromData(data []byte) (*k8Secret, error) {
	seceret := &k8Secret{}
	if err := yaml.Unmarshal(data, seceret); err != nil {
		return nil, err
	}
	return seceret, nil
}
