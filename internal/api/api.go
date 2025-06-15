package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/session"
)

type apiResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Serve(jwtSecret []byte, session *session.Session, port uint16) {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/api/session", func(context *gin.Context) {
		getSession(context, session)
	})
	router.POST("/api/vaults", newVault)
	router.GET("/api/vaults", listDirForVaults)
	router.PUT("/api/vaults/:vaultFile", putItem)
	router.GET("/api/vaults/:vaultFile", func(context *gin.Context) {
		getVault(context, session.SecretKey())
	})
	router.GET("/api/envs", getEnvs)
	router.POST("/api/envs", newEnv)
	router.GET("/api/envs/self", getSelf)
	router.GET("/api/envs/providers", getEnvProviders)

	router.GET("/api/profiles", getProfiles)
	router.POST("/api/profiles", newProfile)
	router.PUT("/api/profiles/:profileName", setActiveProfile)
	router.GET("/api/profiles/remotes", getProfileRemotes)
	router.Run(fmt.Sprintf(":%d", port))
}

func Run() {
	session, err := session.GetSession()
	if err != nil {
		utils.ExitOnError(err)
	}
	jwtSecret := []byte(os.Getenv("SLV_JWT_SECRET"))
	port := uint16(8888)
	Serve(jwtSecret, session, port)
}
