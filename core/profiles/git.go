package profiles

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (profile *Profile) gitLoadRepo() {
	profile.repo, _ = git.PlainOpen(*profile.dir)
	lastUpdated := profile.gitLastUpdated()
	if time.Since(lastUpdated) > profileGitSyncInterval {
		profile.gitPull()
	}
}

func (profile *Profile) gitLastUpdated() time.Time {
	if profile.repo == nil {
		return time.Now()
	}
	fetchHeadPath := filepath.Join(*profile.dir, ".git", "index")
	fetchHeadStat, err := os.Stat(fetchHeadPath)
	if err != nil {
		return time.Time{}
	}
	return fetchHeadStat.ModTime()
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
	_, err = worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			When: time.Now(),
		},
	})
	return err
}

func gitClone(dir, uri, branch string) (*git.Repository, error) {
	if branch != "" {
		return git.PlainClone(dir, false, &git.CloneOptions{
			URL:           uri,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			SingleBranch:  true,
		})
	} else {
		return git.PlainClone(dir, false, &git.CloneOptions{
			URL: uri,
		})
	}
}

func (profile *Profile) gitPull() error {
	if profile.repo == nil {
		return errProfileNotGitRepository
	}
	worktree, err := profile.repo.Worktree()
	if err != nil {
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		Progress: os.Stderr,
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}
