package helpers

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"slv.sh/slv/internal/core/config"
)

const (
	vaultFileNameExt = config.AppNameLowerCase
)

func ListVaultFiles(dir string, recursive bool) (vaultFiles []string, err error) {
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	if recursive {
		err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(d.Name(), "."+vaultFileNameExt) ||
				strings.HasSuffix(d.Name(), vaultFileNameExt+".yaml") ||
				strings.HasSuffix(d.Name(), vaultFileNameExt+".yml") {
				if relPath, err := filepath.Rel(dir, path); err == nil {
					vaultFiles = append(vaultFiles, relPath)
				} else {
					vaultFiles = append(vaultFiles, path)
				}
			}
			return nil
		})
	} else {
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if strings.HasSuffix(file.Name(), "."+vaultFileNameExt) ||
				strings.HasSuffix(file.Name(), vaultFileNameExt+".yaml") ||
				strings.HasSuffix(file.Name(), vaultFileNameExt+".yml") {
				vaultFiles = append(vaultFiles, file.Name())
			}
		}
	}
	return vaultFiles, err
}
