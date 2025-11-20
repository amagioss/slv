package config

const (
	envar_SLV_APP_DATA_DIR = "SLV_APP_DATA_DIR"

	AppNameLowerCase = "slv"
	AppNameUpperCase = "SLV"
	description      = AppNameUpperCase + " (Secure Local Vault) : " + "Securely store, share, and access secrets alongside the codebase."
	website          = "https://slv.sh"
	Art              = ` ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
 ▓▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
   ▓▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒  
    ▓▓▓▓▓▓▓▓▓▓▓   ▓▓▓▓▓▒▒▒▒▒   
      ▓▓▓▓▓▓▓▓▓▓ ▓▓▓▓▓▒▒▒▒     
       ▓▓▓▓▓▓▓▓▓▓▓▓▓▒▒▒▒       
         ▓▓▓▓▓▓▓▓▓▓▒▒▒▒        
           ▓▓▓▓▓▓▓▓▓▓          
         ▓▓▓▓▓▓▓▓▓▓▓▓▓▓        
        ▓▓▓▓▓▒▒▒▓▓▓▓▓▓▓▓       
      ▓▓▓▓▓▓▒▒▒  ▓▓▓▓▓▓▓▓▓     
    ▓▓▓▓▓▓▒▒▒▒    ▓▓▓▓▓▓▓▓▓▓   
   ▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  
 ▓▓▓▓▓▓▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
 ▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒   `

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

	appDataDir *string
)
