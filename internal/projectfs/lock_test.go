package projectfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLockRejectsForeignAndLinkedFiles(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	files, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer files.Close()
	if err := os.WriteFile(filepath.Join(outside, "foreign"), []byte("must survive"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "foreign")); err != nil {
		t.Skipf("native symlinks unavailable: %v", err)
	}
	if lock, err := files.Lock("foreign/new-lock"); err == nil {
		lock.Close()
		t.Error("external lock creation accepted")
	}
	if _, err := os.Stat(filepath.Join(outside, "new-lock")); !os.IsNotExist(err) {
		t.Error("external lock created")
	}
	if err := files.WriteAtomic("regular", []byte("ordinary")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("regular", filepath.Join(root, "linked-lock")); err != nil {
		t.Fatal(err)
	}
	if lock, err := files.Lock("linked-lock"); err == nil {
		lock.Close()
		t.Error("linked lock file accepted")
	}
}

func TestLockWorksWithoutWriteAccessToExistingFile(t *testing.T) {
	root := t.TempDir()
	files, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer files.Close()
	if err := os.WriteFile(filepath.Join(root, "lock"), []byte("unchanged"), 0444); err != nil {
		t.Fatal(err)
	}
	lock, err := files.Lock("lock")
	if err != nil {
		t.Fatal(err)
	}
	if err := lock.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := files.ReadFile("lock")
	if err != nil || string(data) != "unchanged" {
		t.Error("locking altered existing file data")
	}
}
