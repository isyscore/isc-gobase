//go:build !windows

package logger

import (
	"os"
	"syscall"
)

type Strategy struct {
}

func (s Strategy) Dup2(newfd *FileLevelWriter, oldfd *os.File) (err error) {
	return syscall.Dup2(int(newfd.Fd()), int(oldfd.Fd()))
}
