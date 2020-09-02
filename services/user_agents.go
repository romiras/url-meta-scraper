package services

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

var userAgents []string

func loadUserAgents() error {
	rand.Seed(time.Now().UnixNano())
	userAgents = make([]string, 0)
	file, err := os.Open("user_agents.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ua string
	for scanner.Scan() {
		ua = scanner.Text()
		if ua != "" {
			userAgents = append(userAgents, ua)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func getRandomUA() string {
	if userAgents == nil {
		return ""
	}

	idx := rand.Intn(len(userAgents))
	return userAgents[idx]
}
