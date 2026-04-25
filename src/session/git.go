package session

import (
	"os/exec"
	"strings"
)

func findWorktreesFromRealPath(path string) ([]string, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	cmd.Dir = path
	outBytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	out := string(outBytes)

	worktrees := make([]string, 0)

	currentTree := ""

	lines := strings.SplitSeq(out, "\n")
	for line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "worktree") {
			currentTree = line[len("worktree "):]
		} else if line == "bare" {
			currentTree = ""
		} else if currentTree != "" && line == "" {
			worktrees = append(worktrees, currentTree)
			currentTree = ""
		}
	}

	return worktrees, nil
}
