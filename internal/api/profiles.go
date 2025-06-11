package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/core/profiles"
)

func getProfiles(context *gin.Context) {
	activeProfileName, err := profiles.GetActiveProfileName()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	profiles, err := profiles.List()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: map[string]any{
		"active":   activeProfileName,
		"profiles": profiles,
	}})
}

func setActiveProfile(context *gin.Context) {
	profileName := context.Param("profileName")
	if err := profiles.SetActiveProfile(profileName); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: nil})
}

func getProfileRemotes(context *gin.Context) {
	type argResponse struct {
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Required    bool   `json:"required,omitempty"`
		Sensitive   bool   `json:"sensitive,omitempty"`
	}
	remotes := make(map[string][]argResponse)
	remoteNames := profiles.ListRemoteNames()
	for _, remoteName := range remoteNames {
		args := profiles.GetRemoteTypeArgs(remoteName)
		argResponses := make([]argResponse, 0, len(args))
		for _, arg := range args {
			argResponses = append(argResponses, argResponse{
				Name:        arg.Name(),
				Description: arg.Description(),
				Required:    arg.Required(),
				Sensitive:   arg.Sensitive(),
			})
		}
		remotes[remoteName] = argResponses
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: remotes})
}

func newProfile(context *gin.Context) {
	var request struct {
		ProfileName    string            `json:"profileName,omitempty" binding:"required"`
		RemoteType     string            `json:"remoteType,omitempty" binding:"required"`
		ReadOnly       bool              `json:"readOnly,omitempty"`
		UpdateInterval time.Duration     `json:"updateInterval,omitempty"`
		RemoteConfig   map[string]string `json:"remoteConfig,omitempty" binding:"required"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}
	if err := profiles.New(request.ProfileName, request.RemoteType, request.ReadOnly, request.UpdateInterval, request.RemoteConfig); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
}
