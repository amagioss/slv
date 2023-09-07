package projects

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/shibme/slv/core/crypto"
	"github.com/shibme/slv/core/vaults"
)

type Project struct {
	projectRoot string
	vaults      map[string]*vaults.Vault
}

func isProjectConfigDir(dir string) bool {
	dir = filepath.Clean(dir)
	for {
		if f, err := os.Stat(dir); err == nil && f.IsDir() && f.Name() == slvDirName {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
}

func findProjectRoot(dir string) (string, error) {
	dir = filepath.Clean(dir)
	for {
		if f, err := os.Stat(filepath.Join(dir, slvDirName)); err == nil && f.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrProjectNotFound
		}
		dir = parent
	}
}

func NewProjectFor(dir string) (*Project, error) {
	dir = filepath.Clean(dir)
	if isProjectConfigDir(dir) {
		return nil, ErrProhibitedProjectDir
	}
	existingProjectRoot, err := findProjectRoot(dir)
	if err == nil && existingProjectRoot == dir {
		return nil, ErrProjectExists
	}
	projectConfigDirPath := filepath.Join(dir, slvDirName)
	err = os.MkdirAll(projectConfigDirPath, os.FileMode(0755))
	if err != nil {
		return nil, ErrProjectCreation
	}
	return &Project{
		projectRoot: dir,
		vaults:      make(map[string]*vaults.Vault),
	}, nil
}

func NewProject() (*Project, error) {
	currentDir, _ := os.Getwd()
	return NewProjectFor(currentDir)
}

func GetProjectFor(dir string) (project *Project, err error) {
	var projectRoot string
	projectRoot, err = findProjectRoot(dir)
	if err == nil {
		project = &Project{
			projectRoot: projectRoot,
			vaults:      make(map[string]*vaults.Vault),
		}
	}
	return
}

func GetCurrentProject() (*Project, error) {
	currentDir, _ := os.Getwd()
	return GetProjectFor(currentDir)
}

func (project *Project) GetProjectRoot() string {
	return project.projectRoot
}

func validateProjectVaultName(vaultName string) (err error) {
	if match, _ := regexp.MatchString(projectVaultNamePattern, vaultName); !match {
		err = ErrInvalidProjectVaultName
	}
	return
}

func (project *Project) NewVault(vaultName string, hashLength uint32, publicKeys ...crypto.PublicKey) (v *vaults.Vault, err error) {
	if err = validateProjectVaultName(vaultName); err != nil {
		return
	}
	simplifiedVaultName := strings.ToLower(vaultName)
	vaultPath := filepath.Join(project.projectRoot, slvDirName, vaultsDirName, simplifiedVaultName+vaultFileExtension)
	v, err = vaults.New(vaultPath, hashLength, publicKeys...)
	if err == nil {
		project.vaults[vaultName] = v
		project.vaults[simplifiedVaultName] = v
	}
	return
}

func (project *Project) GetVault(vaultName string) (v *vaults.Vault, err error) {
	v = project.vaults[vaultName]
	if v == nil {
		simplifiedVaultName := strings.ToLower(vaultName)
		v = project.vaults[simplifiedVaultName]
		if v == nil {
			if err = validateProjectVaultName(vaultName); err != nil {
				return
			}
			vaultPath := filepath.Join(project.projectRoot, slvDirName, vaultsDirName, simplifiedVaultName+vaultFileExtension)
			v, err = vaults.Get(vaultPath)
			if err == nil {
				project.vaults[vaultName] = v
				project.vaults[simplifiedVaultName] = v
			}
		}
	}
	return
}

func (project *Project) ListVaults() []*vaults.Vault {
	//TODO: Add logic
	return []*vaults.Vault{}
}
