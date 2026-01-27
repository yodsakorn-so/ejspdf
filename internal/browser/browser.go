package browser

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const (
	// Using a known-good revision of Chromium.
	// Source: https://github.com/puppeteer/puppeteer/blob/main/packages/puppeteer-core/src/revisions.ts
	chromiumRevision = "1056772" 
)

// FindOrDownload attempts to find a locally installed Chrome/Chromium executable.
// If not found, it will download a suitable version into a local cache.
func FindOrDownload() (string, error) {
	// 1. First, try to find an installed version
	localPaths := findChromePaths()
	for _, path := range localPaths {
		if p, err := exec.LookPath(path); err == nil {
			log.Println("Found existing browser at:", p)
			return p, nil
		}
	}

	// 2. If not found, try to find it in our local cache
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", err
	}
	
	executablePath := getExecutablePath(cacheDir)
	if _, err := os.Stat(executablePath); err == nil {
		log.Println("Found cached browser at:", executablePath)
		return executablePath, nil
	}

	// 3. If still not found, download it
	log.Printf("Browser not found. Downloading Chromium revision %s to %s\n", chromiumRevision, cacheDir)
	return downloadAndUnzip(cacheDir)
}

// downloadAndUnzip downloads and unzips the appropriate version of Chromium.
func downloadAndUnzip(cacheDir string) (string, error) {
	url, err := getDownloadURL()
	if err != nil {
		return "", err
	}

	// Create a temporary file for the download
	tmpFile, err := os.CreateTemp("", "chromium-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Download the file
	log.Println("Downloading from:", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Setup progress bar
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading chromium",
	)

	_, err = io.Copy(io.MultiWriter(tmpFile, bar), resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	tmpFile.Close() // Close it so we can open it for unzipping

	// Unzip to cache directory
	if err := unzip(tmpFile.Name(), cacheDir); err != nil {
		return "", fmt.Errorf("failed to unzip: %w", err)
	}
	log.Println("Download and extraction complete.")

	// Return the path to the executable
	return getExecutablePath(cacheDir), nil
}

// unzip extracts a zip archive to a destination directory.
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// The top-level directory in the archive is something like "chrome-win",
	// we want to strip that and place the contents directly in cacheDir.
	var firstDir string

	for _, f := range r.File {
		// Determine the base directory
		if firstDir == "" {
			if idx := strings.Index(f.Name, "/"); idx != -1 {
				firstDir = f.Name[:idx]
			}
		}
		
		fpath := filepath.Join(dest, strings.TrimPrefix(f.Name, firstDir+"/"))
		
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}


func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(home, ".cache", "ejspdf", "browser")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}
	return cacheDir, nil
}

func getDownloadURL() (string, error) {
	base := "https://storage.googleapis.com/chromium-browser-snapshots"
	
	switch runtime.GOOS {
	case "windows":
		// For Windows, the structure is slightly different
		return fmt.Sprintf("%s/Win_x64/%s/chrome-win.zip", base, chromiumRevision), nil
	case "darwin":
		if runtime.GOARCH == "arm64" { // Apple Silicon
			return fmt.Sprintf("%s/Mac_Arm/%s/chrome-mac.zip", base, chromiumRevision), nil
		}
		return fmt.Sprintf("%s/Mac/%s/chrome-mac.zip", base, chromiumRevision), nil // Intel
	case "linux":
		return fmt.Sprintf("%s/Linux_x64/%s/chrome-linux.zip", base, chromiumRevision), nil
	}
	return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
}

func getExecutablePath(basePath string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(basePath, "chrome.exe")
	case "darwin":
		return filepath.Join(basePath, "Chromium.app", "Contents", "MacOS", "Chromium")
	case "linux":
		return filepath.Join(basePath, "chrome")
	}
	return ""
}

func findChromePaths() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			"chrome",
			"msedge",
			"chromium",
		}
	case "darwin":
		return []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
	case "linux":
		return []string{
			"google-chrome",
			"microsoft-edge",
			"chromium-browser",
			"chromium",
		}
	default:
		return []string{}
	}
}
