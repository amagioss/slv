package commands

import "github.com/spf13/cobra"

// SLV Command
var slvCmd *cobra.Command

// Profile Commands
var profileCmd *cobra.Command
var profileNewCmd *cobra.Command
var profileListCmd *cobra.Command
var profileSetCmd *cobra.Command
var profileInitRootCmd *cobra.Command

// Environment Commands
var envCmd *cobra.Command
var envNewCmd *cobra.Command
var envAddCmd *cobra.Command
var envListCmd *cobra.Command
var envUserRegisterCmd *cobra.Command

// Vault Commands
var vaultCmd *cobra.Command
var vaultNewCmd *cobra.Command
var vaultShareCmd *cobra.Command

// Secret Commands
var secretCmd *cobra.Command
var secretAddCmd *cobra.Command
var secretGetCmd *cobra.Command
var secretRefCmd *cobra.Command
var secretDerefCmd *cobra.Command
