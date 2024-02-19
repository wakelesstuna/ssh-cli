package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GitAddCommitPush(repoURL string, projectDir string) error {
	// Navigate to the project directory
	if err := os.Chdir(filepath.Join(projectDir)); err != nil {
		return fmt.Errorf("error changing directory: %v", err)
	}

	// Initialize a new Git repository
	if err := exec.Command("git", "init").Run(); err != nil {
		return fmt.Errorf("error initializing git repository: %v", err)
	}

	// Add remote origin
	if err := exec.Command("git", "remote", "add", "origin", repoURL).Run(); err != nil {
		return fmt.Errorf("error adding remote origin: %v", err)
	}

	// Add all files to the staging area
	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return fmt.Errorf("error adding files to git: %v", err)
	}

	// Commit the files
	if err := exec.Command("git", "commit", "-m", "first commit").Run(); err != nil {
		return fmt.Errorf("error committing files: %v", err)
	}

	// Push the commit to the remote repository
	if err := exec.Command("git", "push", "-u", "origin", "master").Run(); err != nil {
		return fmt.Errorf("error pushing to remote repository: %v", err)
	}

	// Cleanup: Remove the .git directory to reverse the git init operation
	if err := os.RemoveAll(filepath.Join(projectDir, ".git")); err != nil {
		return fmt.Errorf("error cleaning up git repository: %v", err)
	}

	return nil
}
