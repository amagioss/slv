package vaults

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/shibme/slv/core/commons"
	"github.com/shibme/slv/core/crypto"
	"gopkg.in/yaml.v3"
)

type autoRef struct {
	keyId   string
	fileId  string
	objPath string
}

func (ref autoRef) String() string {
	return getRefStr(ref.keyId, ref.fileId, ref.objPath)
}

func getRefStr(keyId, fileId, objPath string) string {
	return fmt.Sprintf("%s_%s_%s_%s_%s", commons.SLV, autoRefAbbrev,
		keyId, fileId, commons.Encode([]byte(objPath)))
}

func getAutoRefFromStr(refStr string) (*autoRef, error) {
	sliced := strings.Split(refStr, "_")
	if len(sliced) != 5 || sliced[0] != commons.SLV || sliced[1] != autoRefAbbrev {
		return nil, ErrInvalidAutoRefString
	}
	return &autoRef{
		keyId:   sliced[2],
		fileId:  sliced[3],
		objPath: string(commons.Decode(sliced[4])),
	}, nil
}

func randomStr(bytecount uint8) (string, error) {
	randBytes := make([]byte, bytecount)
	if _, err := io.ReadFull(rand.Reader, randBytes); err != nil {
		return "", err
	}
	return commons.Encode(randBytes), nil
}

func getCurrentFileId(refFile string) (string, error) {
	file, err := os.Open(refFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	autoRefPattern := regexp.MustCompile(commons.SLV + "_" + autoRefAbbrev + "_" +
		"[A-Za-z0-9]+_[A-Za-z0-9]+_[A-Za-z0-9]+")
	for scanner.Scan() {
		line := scanner.Text()
		if match := autoRefPattern.FindString(line); match != "" {
			autoRef, err := getAutoRefFromStr(match)
			if err != nil {
				return "", err
			}
			return autoRef.fileId, nil
		}
	}
	return "", nil
}

func getFileId(refFile string) (fileId string, err error) {
	fileId, err = getCurrentFileId(refFile)
	if err == nil && fileId == "" {
		fileId, err = randomStr(autoReferenceLength)
	}
	return fileId, err
}

func (vlt *Vault) getAutoReferenceName(fileId, objPath string) string {
	return getRefStr(vlt.Config.PublicKey.IdStr(), fileId, objPath)
}

func isReferencedSecret(str string) bool {
	autoRefPattern := regexp.MustCompile(commons.SLV + "_" + autoRefAbbrev + "_" +
		"[A-Za-z0-9]+_[A-Za-z0-9]+_[A-Za-z0-9]+")
	directRefPattern := regexp.MustCompile(commons.SLV + "_" + directRefAbbrev + "_" +
		"[A-Za-z0-9]+_[A-Za-z0-9]+_\\[[A-Za-z0-9]+\\]")
	return autoRefPattern.MatchString(str) || directRefPattern.MatchString(str)
}

func (vlt *Vault) autoReferenceSecret(fileId, objPath, secret string) (secretRef string, err error) {
	var sealedSecret *crypto.SealedSecret
	sealedSecret, err = vlt.Config.PublicKey.EncryptSecretString(secret, vlt.Config.HashLength)
	if err == nil {
		if vlt.vault.Secrets.Referenced == nil {
			vlt.vault.Secrets.Referenced = make(map[string]*crypto.SealedSecret)
		}
		secretRef = vlt.getAutoReferenceName(fileId, objPath)
		vlt.vault.Secrets.Referenced[secretRef] = sealedSecret
		err = vlt.commit()
	}
	return
}

func (vlt *Vault) yamlTraverseAndUpdateRefSecrets(data *map[string]interface{},
	fileId string, path []string, previewOnly bool) (err error) {
	for key, value := range *data {
		switch v := value.(type) {
		case map[string]interface{}:
			if err = vlt.yamlTraverseAndUpdateRefSecrets(&v, fileId, append(path, key),
				previewOnly); err != nil {
				return err
			}
		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err = vlt.yamlTraverseAndUpdateRefSecrets(&itemMap, fileId,
						append(path, key+"["+strconv.Itoa(i)+"]"), previewOnly); err != nil {
						return err
					}
				}
			}
		case string:
			if !isReferencedSecret(v) {
				objPath := strings.Join(append(path, key), ".")
				if previewOnly {
					v = vlt.getAutoReferenceName(fileId, objPath)
				} else {
					if v, err = vlt.autoReferenceSecret(fileId, objPath, v); err != nil {
						return err
					}
				}
				(*data)[key] = v
			}
		}
	}
	return nil
}

func (vlt *Vault) RefSecrets(file string, previewOnly bool) (string, error) {
	if !strings.HasSuffix(file, ".yaml") && !strings.HasSuffix(file, ".yml") {
		return "", ErrInvalidReferenceFileFormat
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	var yMap map[string]interface{}
	err = yaml.Unmarshal(data, &yMap)
	if err != nil {
		return "", err
	}
	fileId, err := getFileId(file)
	if err != nil {
		return "", err
	}
	err = vlt.yamlTraverseAndUpdateRefSecrets(&yMap, fileId, []string{}, previewOnly)
	if err == nil {
		updatedYaml, err := yaml.Marshal(yMap)
		if err == nil {
			if !previewOnly {
				err = os.WriteFile(file, updatedYaml, 0644)
				vlt.commit()
			}
			return string(updatedYaml), err
		}
	}
	return "", err
}
