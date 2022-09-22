//go:build linux && arm64

package logger

import (
	"os"
	"syscall"
)

type Strategy struct {
}

func (s Strategy) Dup2(newfd *FileLevelWriter, oldfd *os.File) (err error) {
	return syscall.Dup3(int(oldfd.Fd()), int(newfd.Fd()), 0)
}
