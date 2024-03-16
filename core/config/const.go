package config

const (
	envar_SLV_ENV_SECRET_KEY      = "SLV_ENV_SECRET_KEY"
	envar_SLV_ENV_SECRET_BINDING  = "SLV_ENV_SECRET_BINDING"
	envar_SLV_ENV_SECRET_PASSWORD = "SLV_ENV_SECRET_PASSWORD"
	envar_SLV_ADMIN_MODE_ENABLED  = "SLV_ADMIN_MODE_ENABLED"
	envar_SLV_APP_DATA_DIR        = "SLV_APP_DATA_DIR"
	envar_SLV_GIT_HTTP_USER       = "SLV_GIT_HTTP_USER"
	envar_SLV_GIT_HTTP_PASS       = "SLV_GIT_HTTP_PASS"
	envar_SLV_GIT_HTTP_TOKEN      = "SLV_GIT_HTTP_TOKEN"

	AppNameLowerCase = "slv"
	AppNameUpperCase = "SLV"
	Description      = AppNameUpperCase + " (Secure Local Vault) : " + "Securely store, share, and access secrets alongside the codebase."
	Website          = "https://savesecrets.org/slv"
	Art              = `
	 ____  _ __     __
	/ ___|| |\ \   / /
	\___ \| | \ \ / / 
	 ___) | |__\ V /  
	|____/|_____\_/   
	`

	K8SLVGroup      = "slv.savesecrets.org"
	K8SLVVersion    = "v1"
	K8SLVKind       = AppNameUpperCase
	K8SLVVaultField = "spec"
)

var (
	adminMode         *bool
	envSecretKey      *string
	envSecretBinding  *string
	envSecretPassword *string
	appDataDir        *string
	gitHTTPUser       *string
	gitHTTPToken      *string
)
