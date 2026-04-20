package session

type Session struct {
	Name     string
	Path     string
	IsActive bool
}

func FindSessions(paths []string) ([]Session, error) {
	return findSessionsFromPaths(paths)
}
