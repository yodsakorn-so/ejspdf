package browser

import (
	"os"
	"testing"
)

func TestFindOrDownload(t *testing.T) {
	// 1. Run the function
	path, err := FindOrDownload()
	
	// 2. Assert results
	if err != nil {
		t.Fatalf("FindOrDownload failed: %v", err)
	}

	if path == "" {
		t.Fatal("FindOrDownload returned empty path")
	}

	// 3. Verify the file actually exists
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Returned path does not exist: %s, error: %v", path, err)
	}

	if info.IsDir() {
		t.Fatalf("Returned path is a directory, expected file: %s", path)
	}

	t.Logf("Success! Found/Downloaded browser at: %s", path)
}
