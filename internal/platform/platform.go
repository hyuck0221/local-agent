package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type OS string

const (
	Darwin  OS = "darwin"
	Linux   OS = "linux"
	Windows OS = "windows"
)

func Current() OS {
	return OS(runtime.GOOS)
}

func Arch() string {
	return runtime.GOARCH
}

// OpencodeConfigPath returns the absolute path to opencode's global config file.
// opencode reads from ~/.config/opencode/opencode.json on Unix and
// %APPDATA%\opencode\opencode.json on Windows.
func OpencodeConfigPath() (string, error) {
	switch Current() {
	case Windows:
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
		return filepath.Join(appdata, "opencode", "opencode.json"), nil
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".config", "opencode", "opencode.json"), nil
	}
}
