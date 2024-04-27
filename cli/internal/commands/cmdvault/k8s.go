package cmdvault

import (
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"oss.amagi.com/slv/cli/internal/commands/utils"
	"oss.amagi.com/slv/core/commons"
	"oss.amagi.com/slv/core/crypto"
	"oss.amagi.com/slv/core/input"
	"oss.amagi.com/slv/core/vaults"
)

type k8Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
	StringData map[string]string `yaml:"stringData"`
	Type       string            `yaml:"type"`
}

func parseK8sSecret(data []byte) (*k8Secret, error) {
	seceret := &k8Secret{}
	if err := yaml.Unmarshal(data, seceret); err != nil {
		return nil, err
	}
	return seceret, nil
}

func toK8slvVaultFile(vault *vaults.Vault, vaultFilePath, k8slvName, k8sSecretType string) error {
	obj := make(map[string]interface{})
	obj[k8sVaultField] = vault
	obj["apiVersion"] = k8sApiVersion
	obj["kind"] = k8sKind
	obj["metadata"] = map[string]interface{}{
		"name": k8slvName,
	}
	if k8sSecretType != "" {
		obj["type"] = k8sSecretType
	}
	return commons.WriteToYAML(vaultFilePath, "", obj)
}

func newK8sVault(filePath, k8sValue string, hashLength uint8, pq bool, rootPublicKey *crypto.PublicKey, publicKeys ...*crypto.PublicKey) (*vaults.Vault, error) {
	k8slvName := k8sValue
	var secretDataMap map[string][]byte
	var k8sSecretType string
	if strings.HasSuffix(k8sValue, ".yaml") || strings.HasSuffix(k8sValue, ".yml") || strings.HasSuffix(k8sValue, ".json") || k8sValue == "-" {
		var data []byte
		var err error
		if k8sValue == "-" {
			data, err = input.ReadBufferFromStdin("Input the k8s secret object as yaml/json: ")
		} else {
			data, err = os.ReadFile(k8sValue)
		}
		if err != nil {
			return nil, err
		}
		secret, err := parseK8sSecret(data)
		if err != nil {
			return nil, err
		}
		k8slvName = secret.Metadata.Name
		secretDataMap = make(map[string][]byte)
		if secret.Data != nil {
			for key, value := range secret.Data {
				decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(value))
				secretValue, err := io.ReadAll(decoder)
				if err != nil {
					return nil, err
				}
				secretDataMap[key] = secretValue
			}
		}
		if secret.StringData != nil {
			for key, value := range secret.StringData {
				secretDataMap[key] = []byte(value)
			}
		}
		k8sSecretType = secret.Type
	}
	vault, err := vaults.New(filePath, k8sVaultField, hashLength, pq, rootPublicKey, publicKeys...)
	if err != nil {
		return nil, err
	}
	if len(secretDataMap) > 0 {
		for key, value := range secretDataMap {
			if err = vault.PutSecret(key, value); err != nil {
				return nil, err
			}
		}
	}
	return vault, toK8slvVaultFile(vault, filePath, k8slvName, k8sSecretType)
}

func vaultToK8sCommand() *cobra.Command {
	if vaultToK8sCmd != nil {
		return vaultToK8sCmd
	}
	vaultToK8sCmd = &cobra.Command{
		Use:     "tok8s",
		Aliases: []string{"k8s", "tok8slv"},
		Short:   "Transform an existing SLV vault file to a K8s compatible one",
		Run: func(cmd *cobra.Command, args []string) {
			vaultFilePath := cmd.Flag(vaultFileFlag.Name).Value.String()
			k8sResourceName := cmd.Flag(vaultK8sNameFlag.Name).Value.String()
			vault, err := getVault(vaultFilePath)
			if err != nil {
				utils.ExitOnError(err)
			}
			if err = toK8slvVaultFile(vault, vaultFilePath, k8sResourceName, ""); err != nil {
				utils.ExitOnError(err)
			}
		},
	}
	vaultToK8sCmd.Flags().StringP(vaultK8sNameFlag.Name, vaultK8sNameFlag.Shorthand, "", vaultK8sNameFlag.Usage)
	vaultToK8sCmd.MarkFlagRequired(vaultK8sNameFlag.Name)
	return vaultToK8sCmd
}
