package main

import (
	"bufio"
	"fmt"
	"harry/session/src/config"
	"harry/session/src/session"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error when getting user config dir: %v\n", err)
		return
	}

	configDir := filepath.Join(userConfigDir, "session")

	conf, err := config.ParseFromConfigDir(configDir)
	if err != nil {
		fmt.Printf("Error when parsing config: %v\n", err)
		return
	}

	sessions, err := session.FindSessions(conf)
	if err != nil {
		fmt.Printf("Error when finding sessions: %v\n", err)
		return
	}

	for i, session := range sessions {
		fmt.Printf("%d: %v\n", i, session)
	}

	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Scan() {
		input := scanner.Text()
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Error when converting input to int: %v\n", err)
			return
		}
		if index < 0 || index >= len(sessions) {
			fmt.Printf("Index out of range: %d\n", index)
			return
		}
		selection := sessions[index]
		err = session.AttachToSession(conf, selection)
		if err != nil {
			fmt.Printf("Error when attaching to session: %v\n", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error when reading input: %v\n", err)
		return
	}
}
