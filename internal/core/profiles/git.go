package profiles

import (
	"fmt"
	"os"
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
	"golang.org/x/crypto/ssh"
	"slv.sh/slv/internal/core/commons"
)

const (
	configGitRepoKey              = "repo"
	configGitBranchKey            = "branch"
	configGitCommitAuthorEmailKey = "committer-email"
	configGitCommitAuthorNameKey  = "committer-name"
	configGitHTTPUserKey          = "username"
	configGitHTTPTokenKey         = "token"
	configGitSSHKey               = "ssh-key"

	gitUrlRegexPattern = `(?i)^(?:(https?|git|ssh):\/\/[\w.@\-~:/]+\.git|git@[\w.\-]+:[\w./~-]+\.git)$`
)

var gitArgs = []arg{
	{
		name:        configGitRepoKey,
		required:    true,
		description: "Git repository URL of the remote profile",
	},
	{
		name:        configGitBranchKey,
		description: "Git branch to be used for the remote profile",
	},
	{
		name:        configGitCommitAuthorEmailKey,
		description: "Email address to be used as the author of the commit",
	},
	{
		name:        configGitCommitAuthorNameKey,
		description: "Name to be used as the author of the commit",
	},
	{
		name:        configGitHTTPUserKey,
		sensitive:   true,
		description: "Username to authenticate with the git repository over HTTP",
	},
	{
		name:        configGitHTTPTokenKey,
		sensitive:   true,
		description: "Token to authenticate with the git repository over HTTP",
	},
	{
		name:        configGitSSHKey,
		sensitive:   true,
		description: "Path to the SSH private key file to authenticate with the git repository over SSH",
	},
}

func getGitAuth(config map[string]string) (auth transport.AuthMethod, err error) {
	gitUrl := config[configGitRepoKey]
	if !regexp.MustCompile(gitUrlRegexPattern).MatchString(gitUrl) {
		return nil, fmt.Errorf("invalid git URL: %s", gitUrl)
	}
	if strings.HasPrefix(gitUrl, "https://") || strings.HasPrefix(gitUrl, "http://") {
		username := config[configGitHTTPUserKey]
		token := config[configGitHTTPTokenKey]
		if username != "" && token != "" {
			auth = &http.BasicAuth{
				Username: username,
				Password: token,
			}
		} else if username != "" || token != "" {
			err = fmt.Errorf("both username and token must be provided for HTTP authentication")
		}
	} else {
		if sshKey := config[configGitSSHKey]; sshKey != "" {
			var keyBytes []byte
			if commons.FileExists(sshKey) {
				if keyBytes, err = os.ReadFile(sshKey); err != nil {
					return nil, fmt.Errorf("failed to read SSH key file %s: %w", sshKey, err)
				}
				config[configGitSSHKey] = string(keyBytes)
			} else {
				keyBytes = []byte(sshKey)
			}
			if _, err = ssh.ParsePrivateKey(keyBytes); err != nil {
				return nil, fmt.Errorf("failed to parse SSH key: %w", err)
			}
			if auth, err = gitssh.NewPublicKeys("git", keyBytes, ""); err != nil {
				return nil, fmt.Errorf("failed to create SSH auth from file %s: %w", sshKey, err)
			}
		} else if auth, err = gitssh.NewSSHAgentAuth("git"); err != nil {
			return nil, fmt.Errorf("failed to create SSH agent auth: %w", err)
		}
	}
	return
}

func gitCommit(repo *git.Repository, msg string, config map[string]string) error {
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
	cfg, _ := gitconfig.LoadConfig(gitconfig.GlobalScope)
	if signature.Email = config[configGitCommitAuthorEmailKey]; signature.Email == "" && cfg != nil {
		signature.Email = cfg.User.Email
	}
	if signature.Name = config[configGitCommitAuthorNameKey]; signature.Name == "" && cfg != nil {
		signature.Name = cfg.User.Name
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
	if cloneOptions.Auth, err = getGitAuth(config); err != nil {
		return err
	}
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
	var auth transport.AuthMethod
	if auth, err = getGitAuth(config); err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		Auth: auth,
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
	if err = gitCommit(repo, note, config); err != nil {
		return err
	}
	var auth transport.AuthMethod
	if auth, err = getGitAuth(config); err != nil {
		return err
	}
	return repo.Push(&git.PushOptions{
		Auth: auth,
	})
}
