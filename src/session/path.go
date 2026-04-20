package session

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

func findSessionsFromPaths(paths []string) ([]Session, error) {
	sessions := make([]Session, 0)
	for _, path := range paths {
		childSessions, err := findSessionsFromPath(path)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, childSessions...)
	}
	return sessions, nil
}

func findSessionsFromPath(path string) ([]Session, error) {
	pathParts := strings.Split(path, "/")

	currentPath := ""

	sessions := make([]Session, 0)

	for i, part := range pathParts {
		if part == "*" {
			partChildren, err := os.ReadDir(currentPath)
			if errors.Is(err, fs.ErrNotExist) {
				return sessions, nil
			}
			if err != nil {
				return nil, err
			}
			for _, child := range partChildren {
				if isVisibleDirectory(child) {
					newRoot := currentPath + "/" + child.Name()
					if i < len(pathParts)-1 {
						newRoot = newRoot + "/" + strings.Join(pathParts[i+1:], "/")
					}
					childSessions, err := findSessionsFromPath(newRoot)
					if err != nil {
						return nil, err
					}
					sessions = append(sessions, childSessions...)
				}
			}
			return sessions, nil
		}
		if i == 0 {
			currentPath = part
		} else {
			currentPath = currentPath + "/" + part
		}
	}

	childSessions, err := sessionChildrenOfPath(currentPath)
	if err != nil {
		return nil, err
	}
	sessions = append(sessions, childSessions...)

	return sessions, nil
}

func sessionChildrenOfPath(path string) ([]Session, error) {
	children, err := os.ReadDir(path)
	if errors.Is(err, fs.ErrNotExist) {
		return make([]Session, 0), nil
	}
	if err != nil {
		return nil, err
	}

	sessions := make([]Session, 0)

	for _, child := range children {
		if isVisibleDirectory(child) {
			sessions = append(sessions, Session{
				Name:     child.Name(),
				Path:     path + "/" + child.Name(),
				IsActive: false,
			})
		}
	}

	return sessions, nil
}

func isVisibleDirectory(d fs.DirEntry) bool {
	return d.IsDir() && d.Name()[0] != '.'
}
