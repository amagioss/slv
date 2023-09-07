package commands

import "github.com/spf13/cobra"

// SLV Command
var slvCmd *cobra.Command

// Config Commands
var configCmd *cobra.Command
var configNewCmd *cobra.Command
var configListCmd *cobra.Command
var configSetCmd *cobra.Command

// Environment Commands
var envCmd *cobra.Command
var envNewCmd *cobra.Command
var envAddCmd *cobra.Command
var envListCmd *cobra.Command
var envRootInitCmd *cobra.Command

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
