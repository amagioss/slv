package environments_new

import (
	"fmt"

	"slv.sh/slv/internal/core/crypto"
	"slv.sh/slv/internal/core/environments"
	"slv.sh/slv/internal/core/environments/envproviders"
	"slv.sh/slv/internal/core/profiles"
)

// createEnvironment creates the environment based on the collected information
func (nep *NewEnvironmentPage) createEnvironment() {
	var env *environments.Environment
	var secretKey *crypto.SecretKey
	var err error

	// Handle direct service creation
	if nep.selectedProvider == "direct" {
		env, secretKey, err = environments.New(nep.envName, nep.selectedType, nep.quantumSafe)
		if err != nil {
			nep.ShowError(fmt.Sprintf("Failed to create environment: %v", err))
			return
		}

		// Store secret key string if available
		if secretKey != nil {
			nep.secretKey = secretKey.String()
		}
	} else {
		// Create environment using provider
		env, err = envproviders.NewEnv(
			nep.selectedProvider,
			nep.envName,
			nep.selectedType,
			nep.providerInputs,
			nep.quantumSafe,
		)
		if err != nil {
			nep.ShowError(fmt.Sprintf("Failed to create environment: %v", err))
			return
		}
	}

	// Set metadata
	if nep.envEmail != "" {
		env.SetEmail(nep.envEmail)
	}
	if len(nep.envTags) > 0 {
		env.AddTags(nep.envTags...)
	}

	// Handle self environment
	if nep.selectedType == environments.USER {
		if err = env.SetAsSelf(); err != nil {
			nep.ShowError(fmt.Sprintf("Failed to set as self environment: %v", err))
			return
		}
	}

	// Add to profile if requested
	if nep.addToProfile {
		profile, err := profiles.GetActiveProfile()
		if err != nil {
			nep.ShowError(fmt.Sprintf("Failed to get active profile: %v", err))
			return
		}

		if !profile.IsPushSupported() {
			nep.ShowError(fmt.Sprintf("Profile (%s) does not support adding environments", profile.Name()))
			return
		}

		if err = profile.PutEnv(env); err != nil {
			nep.ShowError(fmt.Sprintf("Failed to add environment to profile: %v", err))
			return
		}
	}

	// Store created environment
	nep.createdEnv = env

	// Show result
	nep.showResult()
}
