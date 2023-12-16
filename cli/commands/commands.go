package commands

import "github.com/spf13/cobra"

var (

	// SLV Command
	slvCmd *cobra.Command

	// Version Command
	versionCmd *cobra.Command

	// System Commands
	systemCmd      *cobra.Command
	systemResetCmd *cobra.Command

	// Profile Commands
	profileCmd     *cobra.Command
	profileNewCmd  *cobra.Command
	profileListCmd *cobra.Command
	profileSetCmd  *cobra.Command

	// Environment Commands
	envCmd     *cobra.Command
	envNewCmd  *cobra.Command
	envAddCmd  *cobra.Command
	envListCmd *cobra.Command

	// Vault Commands
	vaultCmd      *cobra.Command
	vaultNewCmd   *cobra.Command
	vaultShareCmd *cobra.Command

	// Secret Commands
	secretCmd      *cobra.Command
	secretPutCmd   *cobra.Command
	secretGetCmd   *cobra.Command
	secretRefCmd   *cobra.Command
	secretDerefCmd *cobra.Command
)
