package session

import (
	"errors"
	"harry/session/src/config"
	"io/fs"
	"os"
	"sort"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Session struct {
	Name        string
	Path        string
	ProjectPath string
	Branch      string
	IsActive    bool
}

type sessionKey struct {
	name string
	path string
}

func (session Session) key() sessionKey {
	return sessionKey{
		name: session.Name,
		path: session.Path,
	}
}

func FindSessions(conf config.Config) ([]Session, error) {
	inactiveSessions := make([]Session, 0)
	for _, path := range conf.SearchPaths {
		sessions, err := findSessionsFromPath(path)
		if err != nil {
			return nil, err
		}
		inactiveSessions = append(inactiveSessions, sessions...)
	}

	for _, path := range conf.IncludePaths {
		realPath, err := expandPathHomeDir(path)
		if err != nil {
			return nil, err
		}
		dir, err := os.Stat(realPath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		if !dir.IsDir() {
			continue
		}
		branch, _ := getBranchFromRealPath(realPath)
		inactiveSessions = append(inactiveSessions, Session{
			Name:        dir.Name(),
			Path:        realPath,
			ProjectPath: realPath,
			Branch:      branch,
			IsActive:    false,
		})
	}

	activeSessions, err := findSessionsFromTmux()
	if err != nil {
		return nil, err
	}

	sessions := mergeSessions(inactiveSessions, activeSessions)

	return sessions, nil
}

func mergeSessions(sessionSlices ...[]Session) []Session {
	sessionMap := make(map[sessionKey]Session)

	for _, sessionSlice := range sessionSlices {
		for _, session := range sessionSlice {
			sessionMap[session.key()] = session
		}
	}

	mergedSessions := make([]Session, 0)
	for _, session := range sessionMap {
		mergedSessions = append(mergedSessions, session)
	}
	return mergedSessions
}

func FuzzySearch(search string, sessions []Session) []Session {
	fuzzyStrings := make([]string, 0)
	for _, session := range sessions {
		fuzzyStrings = append(fuzzyStrings, session.Path+session.Name)
	}

	matches := fuzzy.RankFind(search, fuzzyStrings)
	sort.Sort(matches)

	fuzzySessions := make([]Session, 0)
	for _, match := range matches {
		fuzzySessions = append(fuzzySessions, sessions[match.OriginalIndex])
	}

	return fuzzySessions
}
