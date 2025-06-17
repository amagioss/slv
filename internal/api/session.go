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
	data := map[string]any{
		"publicKeys": map[string]string{
			"publicKeyEC": session.PublicKeyEC(),
			"publicKeyPQ": session.PublicKeyPQ(),
		},
	}
	if env != nil {
		data["env"] = env
	}
	context.JSON(http.StatusOK, apiResponse{Success: true, Data: data})
}
