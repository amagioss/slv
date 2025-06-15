package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/vaults"
	"slv.sh/slv/internal/helpers"
)

func listDirForVaults(context *gin.Context) {
	dir := context.Query("dir")
	var err error
	if dir == "" {
		dir, err = os.Getwd()
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
			return
		}
	}
	vaultFiles, err := helpers.ListVaultFiles(dir, false)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	vaultFilesMap := make(map[string]any)
	for _, vaultFile := range vaultFiles {
		vaultFilesMap[vaultFile] = nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	result := make(map[string]string)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			result[name] = "dir"
		} else if _, ok := vaultFilesMap[name]; ok {
			result[name] = "vault"
		} else {
			result[name] = "file"
		}
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: result})
}

func getVault(context *gin.Context, secretKey *crypto.SecretKey) {
	vaultFile := context.Param("vaultFile")
	dir := context.Query("dir")
	if dir != "" {
		vaultFile = filepath.Join(dir, vaultFile)
	}
	vault, err := vaults.Get(vaultFile)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	if context.Query("unlocked") == "true" || context.Query("unlock") == "true" {
		if err = vault.Unlock(secretKey); err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
			return
		}
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: helpers.GetVaultInfo(vault, true, false)})
}

type newVaultRequest struct {
	VaultFile    string   `json:"vaultFile,omitempty"`
	Name         string   `json:"name,omitempty"`
	K8sNamespace string   `json:"k8sNamespace,omitempty"`
	Hash         bool     `json:"hash,omitempty"`
	QuantumSafe  bool     `json:"quantumSafe,omitempty"`
	PublicKeys   []string `json:"publicKeys,omitempty"`
}

func newVault(context *gin.Context) {
	var request newVaultRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}
	dir := context.Query("dir")
	if dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: "directory does not exist"})
			return
		}
		request.VaultFile = filepath.Join(dir, request.VaultFile)
	}
	vault, err := helpers.NewVault(request.VaultFile, request.Name, request.K8sNamespace, request.Hash, request.QuantumSafe, request.PublicKeys)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: helpers.GetVaultInfo(vault, true, false)})
}

func putItem(context *gin.Context) {
	vaultFile := context.Param("vaultFile")
	dir := context.Query("dir")
	if dir != "" {
		vaultFile = filepath.Join(dir, vaultFile)
	}
	var request map[string]*struct {
		Value     string `json:"value,omitempty" binding:"required"`
		PlainText bool   `json:"plainText,omitempty"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}
	vault, err := vaults.Get(vaultFile)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	for key, item := range request {
		if err = vault.Put(key, []byte(item.Value), !item.PlainText); err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
			return
		}
	}
	context.JSON(http.StatusOK, apiResponse{Success: true})
}
