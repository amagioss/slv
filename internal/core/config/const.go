package config

import (
	"strings"

	"github.com/fatih/color"
)

const (
	envar_SLV_APP_DATA_DIR = "SLV_APP_DATA_DIR"

	AppNameLowerCase = "slv"
	AppNameUpperCase = "SLV"
	description      = AppNameUpperCase + " (Secure Local Vault) : " + "Securely store, share, and access secrets alongside the codebase."
	website          = "https://slv.sh"
	art              = ` ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
 ▓▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
   ▓▓▓▓▓▓▓▓▓▓▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒  
    ▓▓▓▓▓▓▓▓▓▓▓   ░░░░░▒▒▒▒▒   
      ▓▓▓▓▓▓▓▓▓▓ ░░░░░▒▒▒▒     
       ▓▓▓▓▓▓▓▓▓▓░░░▒▒▒▒       
         ▓▓▓▓▓▓▓▓▓▓▒▒▒▒        
           ▓▓▓▓▓▓▓▓▓▓          
         ░░░░░▓▓▓▓▓▓▓▓▓        
        ░░░░░▒▒▒▓▓▓▓▓▓▓▓       
      ░░░░░░▒▒▒  ▓▓▓▓▓▓▓▓▓     
    ░░░░░░▒▒▒▒    ▓▓▓▓▓▓▓▓▓▓   
   ░░░░░░▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  
 ░░░░░░▒▒▒▒▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
 ░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒   `

	K8SLVGroup                = "slv.sh"
	K8SLVVersion              = "v1"
	K8SLVKind                 = AppNameUpperCase
	K8SLVAnnotationVersionKey = K8SLVGroup + "/version"
	K8SLVVaultField           = "spec"
)

func coloredArt() string {
	rubyColored := strings.ReplaceAll(art, "▓", color.RGB(157, 58, 79).Sprint("▓"))
	grayColored := strings.ReplaceAll(rubyColored, "░", color.RGB(79, 85, 89).Sprint("░"))
	return strings.ReplaceAll(grayColored, "▒", color.RGB(79, 85, 89).Sprint("▒"))
}

var (
	Version    = "v" + version
	version    = "777.77.77"
	fullCommit = ""
	commitDate = ""
	releaseURL = ""
	appInfo    *string

	appDataDir *string

	Art = coloredArt()
)
