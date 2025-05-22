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
)

const (
	configGitRepoKey      = "repo"
	configGitBranchKey    = "branch"
	configGitHTTPUserKey  = "username"
	configGitHTTPTokenKey = "token"
	configGitHTTPSSHKey   = "ssh-key"
)

var gitArgs = []arg{
	{
		name:        configGitRepoKey,
		required:    true,
		description: "The Git repository URL of the remote profile",
	},
	{
		name:        configGitBranchKey,
		description: "The Git branch to be used for the remote profile",
	},
	{
		name:        configGitHTTPUserKey,
		sensitive:   true,
		description: "The username to authenticate with the git repository over HTTP",
	},
	{
		name:        configGitHTTPTokenKey,
		sensitive:   true,
		description: "The token to authenticate with the git repository over HTTP",
	},
	{
		name:        configGitHTTPSSHKey,
		sensitive:   true,
		description: "The path to the SSH private key file to authenticate with the git repository over SSH",
	},
}

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

func getGitAuth(config map[string]string) transport.AuthMethod {
	gitUrl := config[configGitRepoKey]
	if strings.HasPrefix(gitUrl, "https://") || strings.HasPrefix(gitUrl, "http://") {
		username := config[configGitHTTPUserKey]
		token := config[configGitHTTPTokenKey]
		if username != "" && token != "" {
			return &http.BasicAuth{
				Username: username,
				Password: token,
			}
		}
	}
	if sshAgentAuth, err := gitssh.NewSSHAgentAuth("git"); err == nil {
		return sshAgentAuth
	}
	if sshKeyFile := config[configGitHTTPSSHKey]; sshKeyFile != "" {
		if keyBytes, err := os.ReadFile(sshKeyFile); err == nil {
			_, err = ssh.ParsePrivateKey(keyBytes)
			if err == nil {
				auth, err := gitssh.NewPublicKeysFromFile("git", sshKeyFile, "")
				if err == nil {
					return auth
				}
			}
		}
	}
	if sshKeyFiles := getSSHKeyFiles(gitUrl); len(sshKeyFiles) > 0 {
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

func gitCommit(repo *git.Repository, msg string) error {
	worktree, err := repo.Worktree()
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

func gitSetup(dir string, config map[string]string) (err error) {
	gitUrl := config[configGitRepoKey]
	cloneOptions := &git.CloneOptions{
		URL: gitUrl,
	}
	cloneOptions.Auth = getGitAuth(config)
	branch := config[configGitBranchKey]
	if branch != "" {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(branch)
		cloneOptions.SingleBranch = true
	}
	_, err = git.PlainClone(dir, false, cloneOptions)
	return
}

func gitPull(dir string, config map[string]string) (err error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		Progress: os.Stderr,
		Auth:     getGitAuth(config),
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func gitPush(dir string, config map[string]string, note string) (err error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}
	if err = gitCommit(repo, note); err != nil {
		return err
	}
	return repo.Push(&git.PushOptions{
		Progress: os.Stderr,
		Auth:     getGitAuth(config),
	})
}
