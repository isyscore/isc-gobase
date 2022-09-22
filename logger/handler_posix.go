//go:build (linux && amd64) || darwin

package logger

import (
	"os"
	"syscall"
)

type Strategy struct {
}

func (s Strategy) Dup2(newfd *FileLevelWriter, oldfd *os.File) (err error) {
	return syscall.Dup2(int(oldfd.Fd()), int(newfd.Fd()))
}
