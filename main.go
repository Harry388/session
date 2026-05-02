package main

import (
	"fmt"
	"harry/session/src/config"
	"harry/session/src/session"
	"harry/session/src/ui"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error when getting user config dir: %v\n", err)
		os.Exit(1)
		return
	}

	configDir := filepath.Join(userConfigDir, "session")

	conf, err := config.ParseFromConfigDir(configDir)
	if err != nil {
		fmt.Printf("Error when parsing config: %v\n", err)
		os.Exit(1)
		return
	}

	sessions, err := session.FindSessions(conf)
	if err != nil {
		fmt.Printf("Error when finding sessions: %v\n", err)
		os.Exit(1)
		return
	}

	program := tea.NewProgram(ui.InitialModel(conf, sessions))
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error when running program: %v\n", err)
		os.Exit(1)
		return
	}
}
