package session

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func findSessionsFromTmux() ([]Session, error) {
	names := listTmuxSessionsF("#{session_name}")
	paths := listTmuxSessionsF("#{session_path}")

	if len(names) != len(paths) {
		return nil, errors.New("Unable to read tmux sessions")
	}

	sessions := make([]Session, 0)

	for i, name := range names {
		sessions = append(sessions, Session{
			Name:     name,
			Path:     paths[i],
			IsActive: true,
		})
	}

	return sessions, nil
}

func listTmuxSessionsF(format string) []string {
	cmd := exec.Command("tmux", "list-sessions", "-F", format)
	out, err := cmd.Output()
	if err != nil {
		return make([]string, 0)
	}
	str := string(out)
	trimmed := strings.TrimSpace(str)
	lines := strings.Split(trimmed, "\n")
	return lines
}

func AttachToSession(session Session) error {
	if !session.IsActive {
		err := startNewTmuxSession(session)
		if err != nil {
			return err
		}
	}
	err := attachTmuxToSession(session)
	if err != nil {
		return err
	}
	return nil
}

func attachTmuxToSession(session Session) error {
	inSession := os.Getenv("TMUX") != ""
	if inSession {
		cmd := exec.Command("tmux", "switch", "-t", session.Name)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
		return nil
	} else {
		cmd := exec.Command("tmux", "attach", "-t", session.Name)
		cmd.Stdin = os.Stdin
		_, err := cmd.Output()
		if err != nil {
			return err
		}
		return nil
	}
}

func startNewTmuxSession(session Session) error {
	cmd := exec.Command("tmux", "new-session", "-c", session.Path, "-s", session.Name, "-d")
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
