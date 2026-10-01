package projectfs

import (
	"os"

	"golang.org/x/sys/windows"
)

func lockFile(file *os.File) error {
	return windows.LockFileEx(windows.Handle(file.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, ^uint32(0), ^uint32(0), &windows.Overlapped{})
}
