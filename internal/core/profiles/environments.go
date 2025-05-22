package profiles

import (
	"os"

	"slv.sh/slv/internal/core/environments"
)

func (profile *Profile) getEnvManifestToWrite() (*environments.EnvManifest, error) {
	if !profile.IsPushSupported() {
		return nil, errRemotePushNotSupported
	}
	if err := profile.Pull(); err != nil {
		return nil, err
	}
	return profile.getEnvManifest()
}

func (profile *Profile) pushAndUndoOnError(note string) error {
	if err := profile.Push(note); err != nil {
		if e := os.RemoveAll(profile.dataDir); e != nil {
			return e
		}
		return err
	}
	return nil
}

func (profile *Profile) PutEnv(env *environments.Environment) error {
	envManifest, err := profile.getEnvManifestToWrite()
	if err != nil {
		return err
	}
	if err = envManifest.PutEnv(env); err != nil {
		return err
	}
	return profile.pushAndUndoOnError("Adding environment: " + env.PublicKey + " [" + env.Name + "]")
}

func (profile *Profile) SetRoot(env *environments.Environment) error {
	envManifest, err := profile.getEnvManifestToWrite()
	if err != nil {
		return err
	}
	if err = envManifest.SetRoot(env); err != nil {
		return err
	}
	return profile.pushAndUndoOnError("Setting root environment: " + env.PublicKey + " [" + env.Name + "]")
}

func (profile *Profile) GetRoot() (*environments.Environment, error) {
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.GetRoot()
}

func (profile *Profile) SearchEnvs(queries []string) ([]*environments.Environment, error) {
	envManifest, err := profile.getEnvManifest()
	if err != nil {
		return nil, err
	}
	return envManifest.SearchEnvs(queries), nil
}

func (profile *Profile) DeleteEnv(id string) error {
	envManifest, err := profile.getEnvManifestToWrite()
	if err != nil {
		return err
	}
	var env *environments.Environment
	if env, err = envManifest.DeleteEnv(id); err != nil {
		return err
	}
	return profile.pushAndUndoOnError("Removing environment: " + env.PublicKey + " [" + env.Name + "]")
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
