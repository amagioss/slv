package projects

import (
	"errors"
)

const (
	slvDirName              = ".slv"
	vaultsDirName           = "vaults"
	vaultFileExtension      = ".vault.slv"
	projectVaultNamePattern = "^[a-zA-Z][a-zA-Z0-9_]*[a-zA-Z0-9]$"
)

var ErrProjectNotFound = errors.New("not a project directory")
var ErrProhibitedProjectDir = errors.New("the given directory is inside the project config directory")
var ErrProjectCreation = errors.New("failed to create project in the given directory")
var ErrProjectExists = errors.New("a project already exists in the given directory")
var ErrInvalidProjectVaultName = errors.New("invalid vault name")
