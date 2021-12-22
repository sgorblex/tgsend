package utils

import (
	"errors"
	"os"
	"runtime"
)

func UserDataDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}

	case "darwin", "ios":
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir += "/Library/Application Support"

	default: // Unix
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_DATA_HOME nor $HOME are defined")
			}
			dir += "/.local/share"
		}
	}

	return dir, nil
}
