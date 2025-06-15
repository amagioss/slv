package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
	"slv.sh/slv/internal/core/profiles"
)

func getSelf(context *gin.Context) {
	self := environments.GetSelf()
	if self == nil {
		context.AbortWithStatusJSON(http.StatusNotFound, apiResponse{Success: false, Error: "self environment not found"})
		return
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: self})
}

func setSelf(context *gin.Context) {
	var request struct {
		EnvDef        string `json:"envDef,omitempty" binding:"required"`
		SecretBinding string `json:"secretBinding,omitempty"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}
	env, err := environments.FromDefStr(request.EnvDef)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}
	if request.SecretBinding != "" {
		env.SecretBinding = request.SecretBinding
	}
	if env.SecretBinding == "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: "secret binding is required"})
		return
	}
	if err = env.SetAsSelf(); err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: env})
}

func getEnvs(context *gin.Context) {
	profile, err := profiles.GetActiveProfile()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	var envs []*environments.Environment
	queries := []string{}
	query := context.Query("query")
	if context.Query("query") != "" {
		queries = append(queries, query)
	}
	if query = context.Query("search"); query != "" {
		queries = append(queries, query)
	}
	if query = context.Query("q"); query != "" {
		queries = append(queries, query)
	}
	if query = context.Query("s"); query != "" {
		queries = append(queries, query)
	}
	if len(queries) > 0 {
		envs, err = profile.SearchEnvs(queries)
	} else {
		envs, err = profile.ListEnvs()
	}
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
		return
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: envs})
}

func getEnvProviders(context *gin.Context) {
	type providerArgResponse struct {
		Id          string `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Required    bool   `json:"required,omitempty"`
		Description string `json:"description,omitempty"`
	}
	type providerResponse struct {
		Id          string                `json:"id,omitempty"`
		Name        string                `json:"name,omitempty"`
		Description string                `json:"description,omitempty"`
		Args        []providerArgResponse `json:"args,omitempty"`
	}
	providerIds := envproviders.ListIds()
	providers := make(map[string]providerResponse)
	for _, providerId := range providerIds {
		if providerId != envproviders.PasswordProviderId {
			provider := providerResponse{
				Id:          providerId,
				Name:        envproviders.GetName(providerId),
				Description: envproviders.GetDesc(providerId),
			}
			args := envproviders.GetArgs(providerId)
			provider.Args = make([]providerArgResponse, 0, len(args))
			for _, arg := range args {
				provider.Args = append(provider.Args, providerArgResponse{
					Id:          arg.Id(),
					Name:        arg.Name(),
					Required:    arg.Required(),
					Description: arg.Description(),
				})
			}
			providers[providerId] = provider
		}
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: providers})
}

func newEnv(context *gin.Context) {
	var request struct {
		Name           string            `json:"name,omitempty" binding:"required"`
		Email          string            `json:"email,omitempty"`
		Tags           []string          `json:"tags,omitempty"`
		QuantumSafe    bool              `json:"quantumSafe,omitempty"`
		ProviderId     string            `json:"providerId,omitempty"`
		ProviderInputs map[string]string `json:"inputs,omitempty"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, apiResponse{Success: false, Error: err.Error()})
		return
	}
	envType := environments.SERVICE
	omitBindingInESB := false
	if request.ProviderId == "" {
		env, sk, err := environments.New(request.Name, envType, request.QuantumSafe)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
			return
		}
		env.SetEmail(request.Email)
		env.AddTags(request.Tags...)
		esb, err := env.ToDefStr(omitBindingInESB)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
			return
		}
		context.JSON(http.StatusOK, apiResponse{Success: true, Data: map[string]any{
			"env": env,
			"sk":  sk.String(),
			"esb": esb,
		}})
		return
	} else {
		if request.ProviderId == envproviders.PasswordProviderId {
			envType = environments.USER
			omitBindingInESB = true
		}
		env, err := envproviders.NewEnv(request.ProviderId, request.Name, envType, request.ProviderInputs, request.QuantumSafe)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
			return
		}
		esb, err := env.ToDefStr(omitBindingInESB)
		if err != nil {
			context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
			return
		}
		env.SetEmail(request.Email)
		env.AddTags(request.Tags...)
		context.JSON(http.StatusOK, apiResponse{Success: true, Data: map[string]any{
			"env": env,
			"esb": esb,
		}})
		return
	}
}
