package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/core/session"
)

func getSession(context *gin.Context, session *session.Session) {
	if session.SecretKey() == nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, apiResponse{Success: false, Error: "Unauthorized"})
		return
	}
	env, err := session.Env()
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, apiResponse{Success: false, Error: err.Error()})
	}
	if env == nil {
		context.JSON(http.StatusNotFound, apiResponse{Success: true, Data: map[string]string{
			"publicKeyEC": session.PublicKeyEC(),
			"publicKeyPQ": session.PublicKeyPQ(),
		}})
		return
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: env})
}
