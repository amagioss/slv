package profiles

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"savesecrets.org/slv/core/commons"
	"savesecrets.org/slv/core/config"
	"savesecrets.org/slv/core/crypto"
	"savesecrets.org/slv/core/environments"
	"savesecrets.org/slv/core/settings"
)

type Profile struct {
	name        *string
	dir         *string
	settings    *settings.Settings
	envManifest *environments.EnvManifest
	repo        *git.Repository
}

func (profile *Profile) Name() string {
	return *profile.name
}

func (profile *Profile) commit(msg string) error {
	if msg != "" {
		return profile.gitCommit(msg)
	}
	return nil
}

func newProfile(dir, gitURI, gitBranch string) (profile *Profile, err error) {
	if commons.DirExists(dir) {
		return nil, errProfilePathExistsAlready
	}
	var repo *git.Repository
	if gitURI != "" {
		repo, err = gitClone(dir, gitURI, gitBranch)
		if err != nil {
			return nil, err
		}
	} else if err = os.MkdirAll(dir, 0755); err != nil {
		return nil, errCreatingProfileDir
	}
	profile = &Profile{
		dir:  &dir,
		repo: repo,
	}
	if err = profile.commit(""); err != nil {
		return nil, err
	}
	return profile, nil
}

func getProfileForPath(dir string) (*Profile, error) {
	if !commons.DirExists(dir) {
		return nil, errProfilePathDoesNotExist
	}
	profile := &Profile{
		dir: &dir,
	}
	profile.gitLoadRepo()
	return profile, nil
}

func (profile *Profile) isWriteDenied() bool {
	return profile.repo != nil && !config.IsAdminModeEnabled()
}

func (profile *Profile) GetSettings() (*settings.Settings, error) {
	if profile.settings == nil {
		settingsManifest, err := settings.GetManifest(filepath.Join(*profile.dir, profileSettingsFileName))
		if err != nil {
			settingsManifest, err = settings.NewManifest(filepath.Join(*profile.dir, profileSettingsFileName))
			if err != nil {
				return nil, err
			}
		}
		profile.settings = settingsManifest
	}
	return profile.settings, nil
}

func (profile *Profile) getEnvManifest() (*environments.EnvManifest, error) {
	if profile.envManifest == nil {
		envManifest, err := environments.GetManifest(filepath.Join(*profile.dir, profileEnvironmentsFileName))
		if err != nil {
			envManifest, err = environments.NewManifest(filepath.Join(*profile.dir, profileEnvironmentsFileName))
			if err != nil {
				return nil, err
			}
		}
		profile.envManifest = envManifest
	}
	return profile.envManifest, nil
}

func (profile *Profile) PutEnv(env *environments.Environment) error {
	if profile.isWriteDenied() {
		return errChangesNotAllowedInGitProfile
	}
	if profile.repo != nil {
		if err := profile.Pull(); err != nil {
			return err
		}
	}
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return err
	}
	if err = envManifest.PutEnv(env); err != nil {
		return err
	}
	return profile.commit("Adding environment: " + env.Id() + " [" + env.Name + "]")
}

func (profile *Profile) RootPublicKey() (*crypto.PublicKey, error) {
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.RootPublicKey()
}

func (profile *Profile) SetRoot(env *environments.Environment) error {
	if profile.isWriteDenied() {
		return errChangesNotAllowedInGitProfile
	}
	if profile.repo != nil {
		if err := profile.Pull(); err != nil {
			return err
		}
	}
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return err
	}
	if err = envManifest.SetRoot(env); err != nil {
		return err
	}
	return profile.commit("Setting root environment: " + env.Id() + " [" + env.Name + "]")
}

func (profile *Profile) SearchEnvs(query string) ([]*environments.Environment, error) {
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.SearchEnvs(query), nil
}

func (profile *Profile) SearchEnvsForQueries(queries []string) ([]*environments.Environment, error) {
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.SearchEnvsForQueries(queries), nil
}

func (profile *Profile) DeleteEnv(id string) error {
	if profile.isWriteDenied() {
		return errChangesNotAllowedInGitProfile
	}
	if profile.repo != nil {
		if err := profile.Pull(); err != nil {
			return err
		}
	}
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return err
	}
	env, err := envManifest.DeleteEnv(id)
	if err != nil {
		return err
	}
	return profile.commit("Deleting environment: " + env.Id() + " [" + env.Name + "]")
}

func (profile *Profile) ListEnvs() ([]*environments.Environment, error) {
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.ListEnvs(), nil
}

func (profile *Profile) GetEnv(id string) (*environments.Environment, error) {
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.GetEnv(id), nil
}

func (profile *Profile) Pull() error {
	return profile.gitPull()
}

func (profile *Profile) Push() error {
	return profile.gitPush()
}
