package githubUtil

import (
	"os"
	"os/exec"
)

func GitCommitAndPush(dir, username, repoName string) error {
	cmds := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "Initial commit with Dockerfile and CI setup"},
		{"git", "push", "origin", "main"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
