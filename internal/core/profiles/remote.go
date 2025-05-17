package profiles

var remotes = make(map[string]*remote)

func RegisterRemote(name string, setup setup, pull pull, push push, args []arg) {
	remotes[name] = &remote{
		name:  name,
		setup: setup,
		pull:  pull,
		push:  push,
		args:  args,
	}
}

func RegisterDefaultRemotes() {
	if !defaultRemoteRegistered {
		defaultRemoteRegistered = true
		RegisterRemote("git", gitSetup, gitPull, gitPush, gitArgs)
		RegisterRemote("http", httpSetup, httpPull, nil, httpArgs)
	}
}

func ListRemoteTypes() []string {
	RegisterDefaultRemotes()
	remoteNames := make([]string, 0, len(remotes))
	for name := range remotes {
		remoteNames = append(remoteNames, name)
	}
	return remoteNames
}

func GetRemoteTypeArgs(name string) []arg {
	if remote, ok := remotes[name]; ok {
		return remote.args
	}
	return nil
}

type setup func(dir string, config map[string]string) error
type pull func(dir string, config map[string]string) error
type push func(dir string, config map[string]string, note string) error

type remote struct {
	name  string
	setup setup
	pull  pull
	push  push
	args  []arg
}

type arg struct {
	name        string
	required    bool
	description string
}

func (a *arg) Name() string {
	return a.name
}

func (a *arg) Required() bool {
	return a.required
}

func (a *arg) Description() string {
	return a.description
}
