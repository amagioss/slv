package config

const (
	envar_SLV_ENV_SECRET_KEY      = "SLV_ENV_SECRET_KEY"
	envar_SLV_ENV_SECRET_BINDING  = "SLV_ENV_SECRET_BINDING"
	envar_SLV_ENV_SECRET_PASSWORD = "SLV_ENV_SECRET_PASSWORD"
	envar_SLV_APP_DATA_DIR        = "SLV_APP_DATA_DIR"

	AppNameLowerCase = "slv"
	AppNameUpperCase = "SLV"
	description      = AppNameUpperCase + " (Secure Local Vault) : " + "Securely store, share, and access secrets alongside the codebase."
	website          = "https://slv.sh"
	art              = `
	 ____  _ __     __
	/ ___|| |\ \   / /
	\___ \| | \ \ / / 
	 ___) | |__\ V /  
	|____/|_____\_/   
	`

	K8SLVGroup                = "slv.sh"
	K8SLVVersion              = "v1"
	K8SLVKind                 = AppNameUpperCase
	K8SLVAnnotationVersionKey = K8SLVGroup + "/version"
	K8SLVVaultField           = "spec"
)

var (
	Version    = "v" + version
	version    = "777.77.77"
	fullCommit = ""
	commitDate = ""
	releaseURL = ""
	appInfo    *string

	envSecretKey      *string
	envSecretBinding  *string
	envSecretPassword *string
	appDataDir        *string
)
