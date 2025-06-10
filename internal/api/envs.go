package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/profiles"
)

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
