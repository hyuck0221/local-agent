package ollama

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hyuck0221/local-agent/internal/platform"
)

const (
	DefaultHost = "http://localhost:11434"
	BaseURL     = DefaultHost + "/v1"
)

// Installed reports whether the `ollama` binary is on PATH.
func Installed() bool {
	_, err := exec.LookPath("ollama")
	return err == nil
}

// Install runs the best-effort OS-appropriate installer for Ollama.
// On macOS it prefers Homebrew; on Linux it uses the official shell installer;
// on Windows it points the user at winget.
func Install() error {
	switch platform.Current() {
	case platform.Darwin:
		if _, err := exec.LookPath("brew"); err == nil {
			return run("brew", "install", "ollama")
		}
		return errors.New("homebrew not found; install it from https://brew.sh or download Ollama from https://ollama.com/download")
	case platform.Linux:
		if _, err := exec.LookPath("curl"); err != nil {
			return errors.New("curl is required to install Ollama on Linux")
		}
		return runShell("curl -fsSL https://ollama.com/install.sh | sh")
	case platform.Windows:
		if _, err := exec.LookPath("winget"); err == nil {
			return run("winget", "install", "--id", "Ollama.Ollama", "-e", "--source", "winget")
		}
		return errors.New("winget not found; download Ollama from https://ollama.com/download/windows")
	default:
		return fmt.Errorf("unsupported OS: %s", platform.Current())
	}
}

// Serving reports whether an Ollama server is already accepting TCP on the default port.
func Serving() bool {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:11434", 500*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// StartServer launches `ollama serve` detached from the current process and
// waits up to 5s for the port to accept connections.
func StartServer() error {
	if Serving() {
		return nil
	}
	cmd := exec.Command("ollama", "serve")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start ollama serve: %w", err)
	}
	// Release the child so it outlives us.
	_ = cmd.Process.Release()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if Serving() {
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	return errors.New("ollama serve did not become ready within 5s")
}

// Pull downloads a model if it is not already available locally.
func Pull(model string) error {
	have, err := Has(model)
	if err != nil {
		return err
	}
	if have {
		return nil
	}
	fmt.Printf("Pulling %s (this can take a while on first run)...\n", model)
	cmd := exec.Command("ollama", "pull", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Has reports whether a given model tag is already present in `ollama list`.
func Has(model string) (bool, error) {
	out, err := exec.Command("ollama", "list").Output()
	if err != nil {
		return false, fmt.Errorf("ollama list: %w", err)
	}
	want := strings.TrimSpace(model)
	for _, line := range strings.Split(string(out), "\n")[1:] {
		name := strings.Fields(line)
		if len(name) > 0 && (name[0] == want || strings.HasPrefix(name[0], want+":")) {
			return true, nil
		}
	}
	return false, nil
}

// List returns locally installed models (first column of `ollama list`).
func List() ([]string, error) {
	out, err := exec.Command("ollama", "list").Output()
	if err != nil {
		return nil, err
	}
	var models []string
	lines := strings.Split(string(out), "\n")
	if len(lines) <= 1 {
		return models, nil
	}
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}
	return models, nil
}

// Stop attempts to stop a running Ollama server. On macOS/Linux we simply
// signal any `ollama serve` process; on Windows we rely on `taskkill`.
func Stop() error {
	switch platform.Current() {
	case platform.Windows:
		return run("taskkill", "/IM", "ollama.exe", "/F")
	default:
		// `pkill -f "ollama serve"` is the least-fragile option across distros.
		return run("pkill", "-f", "ollama serve")
	}
}

// Reachable performs a real HTTP GET against /api/tags to distinguish a
// listening port from a working server.
func Reachable() bool {
	client := http.Client{Timeout: time.Second}
	resp, err := client.Get(DefaultHost + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runShell(script string) error {
	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	var buf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
	return cmd.Run()
}
