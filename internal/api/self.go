package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/core/environments"
)

func getSelf(context *gin.Context) {
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: environments.GetSelf()})
}
