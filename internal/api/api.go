package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"slv.sh/slv/internal/cli/commands/utils"
	"slv.sh/slv/internal/core/session"
)

const (
	defaultPort = 47474
	slvLocalUrl = "http://local.slv.sh"
)

type apiResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Serve(jwtSecret []byte, session *session.Session, port uint16) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	apiGroup := router.Group("/api")
	apiGroup.Use(authMiddleware(jwtSecret))

	// Session API
	apiGroup.GET("/session", func(context *gin.Context) {
		getSession(context, session)
	})

	// Vaults API
	vaultAPI := apiGroup.Group("/vaults")
	vaultAPI.POST("", newVault)
	vaultAPI.GET("", listDirForVaults)
	vaultAPI.PUT("/:vaultFile", putItem)
	vaultAPI.GET("/:vaultFile", func(context *gin.Context) {
		getVault(context, session.SecretKey())
	})

	// Environments API
	envAPI := apiGroup.Group("/envs")
	envAPI.GET("/self", getSelf)
	envAPI.PUT("/self", setSelf)
	envAPI.GET("", getEnvs)
	envAPI.POST("", newEnv)
	envAPI.GET("/providers", getEnvProviders)

	// Profiles API
	profileAPI := apiGroup.Group("/profiles")
	profileAPI.GET("", getProfiles)
	profileAPI.POST("", newProfile)
	profileAPI.PUT("/:profileName", setActiveProfile)
	profileAPI.GET("/remotes", getProfileRemotes)

	return router.Run(fmt.Sprintf(":%d", port))
}

func Run(port uint16) error {
	session, err := session.GetSession()
	if err != nil {
		utils.ExitOnError(err)
	}
	jwtSecret := make([]byte, 32)
	if _, err = rand.Read(jwtSecret); err != nil {
		return err
	}
	jwtSecretStr := base64.RawURLEncoding.EncodeToString(jwtSecret)
	if port == 0 {
		port = defaultPort
		for {
			if conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
				conn.Close()
				break
			}
			port++
		}
	}
	url := fmt.Sprintf("%s:%d/#%s", slvLocalUrl, port, jwtSecretStr)
	fmt.Printf("Open this URL in your browser: %s\n", url)
	return Serve(jwtSecret, session, port)
}
