package profiles

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
	"oss.amagi.com/slv/internal/core/config"
)

func expandTilde(path string) string {
	if len(path) > 0 && path[0] == '~' {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[1:])
		}
	}
	return path
}

func getSSHKeyFiles(uri string) []string {
	pattern := regexp.MustCompile(`(?:[^@]+@)?([^:/]+)`)
	matches := pattern.FindStringSubmatch(uri)
	if len(matches) < 2 {
		return nil
	}
	hostname := matches[1]
	if hostname != "" {
		allKeyPaths := ssh_config.GetAll(hostname, "IdentityFile")
		var keyPaths []string
		keyPathMap := make(map[string]struct{})
		for _, keyPath := range allKeyPaths {
			keyPath = expandTilde(keyPath)
			if _, found := keyPathMap[keyPath]; !found {
				keyPaths = append(keyPaths, keyPath)
				keyPathMap[keyPath] = struct{}{}
			}
		}
		return keyPaths
	}
	return nil
}

func getGitAuth(gitURI string) transport.AuthMethod {
	if strings.HasPrefix(gitURI, "https://") {
		if !gitHttpAuthProcessed {
			gitHttpUsername := config.GetGitHTTPUsername()
			if gitHttpUsername != "" {
				gitHttpToken := config.GetGitHTTPToken()
				if gitHttpToken != "" {
					gitHttpAuth = &http.BasicAuth{
						Username: gitHttpUsername,
						Password: gitHttpToken,
					}
				}
			}
			gitHttpAuthProcessed = true
		}
		return gitHttpAuth
	}
	if sshKeyFiles := getSSHKeyFiles(gitURI); len(sshKeyFiles) > 0 {
		keyPath := sshKeyFiles[0]
		keyBytes, err := os.ReadFile(keyPath)
		if err == nil {
			_, err = ssh.ParsePrivateKey(keyBytes)
			if err == nil {
				auth, err := gitssh.NewPublicKeysFromFile("git", keyPath, "")
				if err == nil {
					return auth
				}
			}
		}
	}
	return nil
}

func (profile *Profile) getGitAuth() (transport.AuthMethod, error) {
	remotes, err := profile.repo.Remotes()
	if err != nil {
		return nil, err
	}
	if len(remotes) == 0 {
		return nil, errProfileNotGitRepository
	} else {
		return getGitAuth(remotes[0].Config().URLs[0]), nil
	}
}

func (profile *Profile) gitLoadRepo() {
	profile.repo, _ = git.PlainOpen(*profile.dir)
	lastUpdated := profile.gitLastPulledAt()
	if time.Since(lastUpdated) > profileGitSyncInterval {
		profile.gitPull()
	}
}

func (profile *Profile) gitLastPulledAt() time.Time {
	if profile.repo == nil {
		return time.Now()
	}
	slvPullMarkFile := filepath.Join(*profile.dir, ".git", ".slv-pull")
	slvPullMarkStat, err := os.Stat(slvPullMarkFile)
	if err != nil {
		return time.Time{}
	}
	return slvPullMarkStat.ModTime()
}

func (profile *Profile) gitCommit(msg string) error {
	if profile.repo == nil {
		return nil
	}
	worktree, err := profile.repo.Worktree()
	if err != nil {
		return err
	}
	if _, err = worktree.Add("."); err != nil {
		return err
	}
	signature := &object.Signature{
		When: time.Now(),
	}
	cfg, err := gitconfig.LoadConfig(gitconfig.GlobalScope)
	if err == nil {
		if userEmail := cfg.User.Email; userEmail != "" {
			signature.Email = userEmail
		}
		if userName := cfg.User.Name; userName != "" {
			signature.Name = userName
		}
	}
	_, err = worktree.Commit(msg, &git.CommitOptions{
		Author: signature,
	})
	return err
}

func gitClone(dir, uri, branch string) (*git.Repository, error) {
	cloneOptions := &git.CloneOptions{
		URL: uri,
	}
	if auth := getGitAuth(uri); auth != nil {
		cloneOptions.Auth = auth
	}
	if branch != "" {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(branch)
		cloneOptions.SingleBranch = true
	}
	return git.PlainClone(dir, false, cloneOptions)
}

func (profile *Profile) gitMarkPull() error {
	slvPullMarkFile := filepath.Join(*profile.dir, ".git", ".slv-pull")
	if _, err := os.Create(slvPullMarkFile); err != nil {
		return errProfileGitPullMarking
	}
	return nil
}

func (profile *Profile) gitPull() error {
	if profile.repo == nil {
		return errProfileNotGitRepository
	}
	if err := profile.gitMarkPull(); err != nil {
		return err
	}
	worktree, err := profile.repo.Worktree()
	if err != nil {
		return err
	}
	auth, err := profile.getGitAuth()
	if err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		Progress: os.Stderr,
		Auth:     auth,
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func (profile *Profile) gitPush() error {
	if profile.repo == nil {
		return errProfileNotGitRepository
	}
	auth, err := profile.getGitAuth()
	if err != nil {
		return err
	}
	return profile.repo.Push(&git.PushOptions{
		Progress: os.Stderr,
		Auth:     auth,
	})
}
