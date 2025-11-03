package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// openBrowserURL opens the specified URL in the default browser
func openBrowserURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		// Try common Linux browser commands
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", url)
		} else if _, err := exec.LookPath("x-www-browser"); err == nil {
			cmd = exec.Command("x-www-browser", url)
		} else {
			return fmt.Errorf("no browser command found (tried xdg-open, x-www-browser)")
		}
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start() // Don't wait for browser to close
}

func hasFrontendAssets() bool {
	if frontendFiles != nil {
		if entries, err := fs.ReadDir(frontendFiles, "web/dist"); err == nil && len(entries) > 0 {
			return true
		}
	}

	if _, err := os.Stat(filepath.Join("web", "dist", "index.html")); err == nil {
		return true
	}

	return false
}

// startServerAndOpenBrowser starts the serve command and opens the browser
func startServerAndOpenBrowser(resultsPath, graphPath string) error {
	if !hasFrontendAssets() {
		fmt.Fprintf(os.Stderr, "⚠️  Frontend not built. Skipping browser open.\n")
		fmt.Fprintf(os.Stderr, "   Run 'make frontend-build' to build the frontend, then run:\n")
		fmt.Fprintf(os.Stderr, "   barracuda serve --results %s", resultsPath)
		if graphPath != "" {
			fmt.Fprintf(os.Stderr, " --graph %s", graphPath)
		}
		fmt.Fprint(os.Stderr, "\n")
		return nil
	}

	// Start server in background
	fmt.Fprintf(os.Stdout, "\n🚀 Starting web server...\n")
	fmt.Fprintf(os.Stdout, "📊 Opening dashboard in browser...\n")

	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine executable path: %w", err)
	}

	// Build serve command
	serveArgs := []string{"serve", "--results", resultsPath, "--port", "8080"}
	if graphPath != "" {
		serveArgs = append(serveArgs, "--graph", graphPath)
	}

	cmd := exec.Command(executable, serveArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start server process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Give server a moment to start
	time.Sleep(2 * time.Second)

	// Open browser
	browserURL := "http://localhost:8080"
	if err := openBrowserURL(browserURL); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Failed to open browser: %v\n", err)
		fmt.Fprintf(os.Stderr, "   You can manually open: %s\n", browserURL)
	} else {
		fmt.Fprintf(os.Stdout, "✓ Browser opened: %s\n", browserURL)
		fmt.Fprintf(os.Stdout, "   Press Ctrl+C to stop the server\n")
	}

	// Wait for server process
	return cmd.Wait()
}
