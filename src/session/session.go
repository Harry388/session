package session

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Session struct {
	Name           string
	WorkingPath    string
	RepositoryPath string
	Branch         string
	IsActive       bool
}

func NewSessionFromWorkingPath(path string, isActive bool) Session {
	session := Session{
		Name:           strings.ReplaceAll(filepath.Base(path), ".", "_"),
		WorkingPath:    path,
		RepositoryPath: path,
		Branch:         "",
		IsActive:       isActive,
	}

	worktrees, err := findGitWorktreesFromPath(path)

	// There are no worktrees, or there is only one worktree in the same directory
	if err != nil ||
		len(worktrees) == 0 ||
		(len(worktrees) == 1 &&
			(worktrees[0].Path == path)) {
		if len(worktrees) == 1 {
			session.Branch = worktrees[0].Branch
		}
		return session
	} else {
		for _, worktree := range worktrees {
			if worktree.Path == path {
				session.Branch = worktree.Branch
			} else if worktree.IsBare {
				session.RepositoryPath = worktree.Path
			}
		}
		return session
	}
}

type SessionFinder interface {
	FindSessions() ([]Session, error)
	MergeSessions(currentSessions []Session, newSessions []Session) []Session
}

func defaultMergeSessions(currentSessions []Session, newSessions []Session) []Session {
	return append(currentSessions, newSessions...)
}

func FindSessions(sources []SessionFinder) ([]Session, error) {
	var sessions []Session
	for _, source := range sources {
		sourceSessions, err := source.FindSessions()
		if err != nil {
			return nil, err
		}
		sessions = source.MergeSessions(sessions, sourceSessions)
	}
	return sessions, nil
}

func FuzzySearch(sessions []Session, search string) []Session {
	fuzzyStrings := make([]string, 0, len(sessions))
	for _, session := range sessions {
		fuzzyStrings = append(fuzzyStrings, session.WorkingPath+session.Name)
	}

	matches := fuzzy.RankFind(search, fuzzyStrings)
	sort.Sort(matches)

	fuzzySessions := make([]Session, 0, len(matches))
	for _, match := range matches {
		fuzzySessions = append(fuzzySessions, sessions[match.OriginalIndex])
	}

	return fuzzySessions
}
