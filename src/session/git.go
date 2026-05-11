package session

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitWorktree struct {
	Path   string
	Branch string
	IsBare bool
}

func findGitWorktreesFromPath(path string) ([]GitWorktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = path
	outBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error when getting git worktree list for %s: %w", path, err)
	}
	out := string(outBytes)

	worktrees := make([]GitWorktree, 0)

	currentTree := GitWorktree{}

	lines := strings.SplitSeq(out, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "worktree "); ok {
			currentTree.Path = after
		} else if after, ok := strings.CutPrefix(line, "branch refs/heads/"); ok {
			currentTree.Branch = after
		} else if line == "bare" {
			currentTree.IsBare = true
		} else if currentTree != (GitWorktree{}) && line == "" {
			worktrees = append(worktrees, currentTree)
			currentTree = GitWorktree{}
		}
	}

	return worktrees, nil
}
