package libo5

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type State struct {
	Domain string `json:"domain"`
}

var reUnsafe = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func readState(domain string) (*State, error) {
	filename, err := stateFilename(domain)
	if err != nil {
		return nil, err
	}
	fileOut, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Domain: domain}, nil
		}
		return nil, fmt.Errorf("read state file %s: %w", filename, err)
	}

	state := &State{}
	err = json.Unmarshal(fileOut, &state)
	if err != nil {
		return nil, fmt.Errorf("unmarshal state file %s: %w", filename, err)
	}

	return state, nil
}

func writeState(state *State) error {
	filename, err := stateFilename(state.Domain)
	if err != nil {
		return err
	}

	fileOut, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state file %s: %w", filename, err)
	}
	return os.WriteFile(filename, fileOut, 0644)
}

func stateFilename(domain string) (string, error) {
	stateDir, err := getStateDir()
	if err != nil {
		return "", fmt.Errorf("getStateDir: %w", err)
	}
	logDir := filepath.Join(stateDir, "o5")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("mkdir %s: %w", logDir, err)
	}

	filename := fmt.Sprintf("%s.json", reUnsafe.ReplaceAllString(domain, "_"))
	return filepath.Join(logDir, filename), nil
}

func getStateDir() (string, error) {
	dir := os.Getenv("XDG_STATE_HOME")
	if dir == "" {
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("neither $XDG_STATE_HOME nor $HOME are defined")
		}
		dir += "/.local/state"
	} else if !filepath.IsAbs(dir) {
		return "", errors.New("path in $XDG_STATE_HOME is relative")
	}

	return dir, nil
}
