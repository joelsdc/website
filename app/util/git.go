package util

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"aahframework.org/aah.v0"
	"aahframework.org/essentials.v0"
	"aahframework.org/log.v0"
)

var (
	// ErrRepoAleadyExists returned if repo/dir is already exists
	ErrRepoAleadyExists = errors.New("repo already exists")

	gitcmdPath string
)

// GitClone method to clone the Git Repository.
func GitClone(destDir, repoURL, branchName string) error {
	if ess.IsFileExists(destDir) {
		return ErrRepoAleadyExists
	}

	cloneArgs := []string{"clone"}

	if !ess.IsStrEmpty(branchName) {
		cloneArgs = append(cloneArgs, "-b", branchName)
	}

	cloneArgs = append(cloneArgs, repoURL, destDir)

	_, err := GitCmd(cloneArgs, true)
	return err
}

// GitPull method to take git update from Git Repository.
func GitPull(destDir string) error {
	if err := os.Chdir(destDir); err != nil {
		return err
	}
	gitDir := filepath.Join(destDir, ".git")
	pullArgs := []string{fmt.Sprintf("--git-dir=%s", gitDir), "pull"}

	_, err := GitCmd(pullArgs, true)
	return err
}

// GitCloneAndPull method does both clone and pull for git Repository.
func GitCloneAndPull(destDir, branchName string) error {
	repoURL := aah.AppConfig().StringDefault("docs.repo", "")
	err := GitClone(destDir, repoURL, branchName)
	if err == ErrRepoAleadyExists {
		if err = GitPull(destDir); err != nil {
			return err
		}
	}
	return err
}

// GitCmd method execute given git commands.
func GitCmd(args []string, stdout bool) (string, error) {
	return ExecCmd(gitcmdPath, args, stdout)
}

func init() {
	if ess.LookExecutable("git") {
		gitcmdPath, _ = exec.LookPath("git")
	} else {
		log.Warn("git program is not installed, site may produce unexpected behaviour")
	}
}
