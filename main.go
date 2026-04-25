package main

import (
	"bufio"
	"fmt"
	"harry/session/src/config"
	"harry/session/src/session"
	"os"
	"strconv"
)

func main() {
	sessions, err := session.FindSessions(config.Config{
		SearchPaths: []string{
			"~/dev/*",
			"~/work/*",
		},
		IncludePaths: []string{
			"~/env",
		},
	})
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
		err = session.AttachToSession(selection)
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
