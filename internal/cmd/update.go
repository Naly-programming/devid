package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update devid to the latest version",
	RunE:  runUpdate,
}

type ghRelease struct {
	TagName string `json:"tag_name"`
}

func runUpdate(cmd *cobra.Command, args []string) error {
	current := rootCmd.Version

	// Get latest version from GitHub
	resp, err := http.Get("https://api.github.com/repos/Naly-programming/devid/releases/latest")
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	var release ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("failed to parse release info: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	if latest == current {
		fmt.Printf("Already on the latest version (v%s)\n", current)
		return nil
	}

	fmt.Printf("Current: v%s\n", current)
	fmt.Printf("Latest:  v%s\n", latest)
	fmt.Println("Updating...")

	// Determine platform
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	ext := "tar.gz"
	if goos == "windows" {
		ext = "zip"
	}

	url := fmt.Sprintf("https://github.com/Naly-programming/devid/releases/download/v%s/devid_%s_%s_%s.%s",
		latest, latest, goos, goarch, ext)

	// Download
	dlResp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer dlResp.Body.Close()

	if dlResp.StatusCode != 200 {
		return fmt.Errorf("download failed: HTTP %d (no release for %s/%s?)", dlResp.StatusCode, goos, goarch)
	}

	tmpDir, err := os.MkdirTemp("", "devid-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, "devid."+ext)
	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, dlResp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()

	// Extract
	if ext == "zip" {
		cmd := exec.Command("powershell", "-Command",
			fmt.Sprintf("Expand-Archive -Force '%s' '%s'", archivePath, tmpDir))
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("extract failed: %s", string(out))
		}
	} else {
		cmd := exec.Command("tar", "-xzf", archivePath, "-C", tmpDir)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("extract failed: %s", string(out))
		}
	}

	// Find current binary location
	currentBin, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find current binary: %w", err)
	}
	currentBin, _ = filepath.EvalSymlinks(currentBin)

	// Find new binary
	newBin := filepath.Join(tmpDir, "devid")
	if goos == "windows" {
		newBin += ".exe"
	}

	if _, err := os.Stat(newBin); err != nil {
		return fmt.Errorf("new binary not found in download")
	}

	// Replace
	// On Windows, can't replace a running binary directly - rename first
	if goos == "windows" {
		oldBin := currentBin + ".old"
		os.Remove(oldBin)
		if err := os.Rename(currentBin, oldBin); err != nil {
			return fmt.Errorf("failed to replace binary (try running as admin): %w", err)
		}
		if err := copyFile(newBin, currentBin); err != nil {
			// Rollback
			os.Rename(oldBin, currentBin)
			return fmt.Errorf("failed to install new binary: %w", err)
		}
		os.Remove(oldBin)
	} else {
		if err := copyFile(newBin, currentBin); err != nil {
			return fmt.Errorf("failed to replace binary (try: sudo devid update): %w", err)
		}
		os.Chmod(currentBin, 0o755)
	}

	fmt.Printf("Updated to v%s\n", latest)
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
