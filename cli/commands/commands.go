package commands

import "github.com/spf13/cobra"

// SLV Command
var slvCmd *cobra.Command

// Version Command
var versionCmd *cobra.Command

// System Commands
var systemCmd *cobra.Command
var systemResetCmd *cobra.Command

// Profile Commands
var profileCmd *cobra.Command
var profileNewCmd *cobra.Command
var profileListCmd *cobra.Command
var profileSetCmd *cobra.Command

// Environment Commands
var envCmd *cobra.Command
var envNewCmd *cobra.Command
var envProviderCmd *cobra.Command
var envAddCmd *cobra.Command
var envListCmd *cobra.Command

// Vault Commands
var vaultCmd *cobra.Command
var vaultNewCmd *cobra.Command
var vaultShareCmd *cobra.Command

// Secret Commands
var secretCmd *cobra.Command
var secretPutCmd *cobra.Command
var secretGetCmd *cobra.Command
var secretRefCmd *cobra.Command
var secretDerefCmd *cobra.Command
