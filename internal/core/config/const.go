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

func ColorizedArt() string {
	if colorizedArt == nil {
		colorizedArt = new(string)
		*colorizedArt = strings.ReplaceAll(art, "▓", color.RGB(157, 58, 79).Sprint("▓"))
		*colorizedArt = strings.ReplaceAll(*colorizedArt, "░", color.RGB(79, 85, 89).Sprint("░"))
		*colorizedArt = strings.ReplaceAll(*colorizedArt, "▒", color.RGB(79, 85, 89).Sprint("▒"))
	}
	return *colorizedArt
}

// Art returns the plain ASCII art without color codes
func Art() string {
	return art
}

var (
	Version    = "v" + version
	version    = "777.77.77"
	fullCommit = ""
	commitDate = ""
	releaseURL = ""
	appInfo    *string

	appDataDir *string

	colorizedArt *string
)

func GetVersion() string {
	return version
}

func GetFullCommit() string {
	return fullCommit
}

func GetCommitDate() string {
	return commitDate
}

func GetReleaseURL() string {
	return releaseURL
}
