//go:build windows

package logger

import (
	"os"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

func setStdHandle(stdHandle int32, handle syscall.Handle) error {
	r0, _, err := syscall.Syscall(procSetStdHandle.Addr(), 2, uintptr(stdHandle), uintptr(handle), 0)
	if r0 == 0 {
		if err != 0 {
			return error(err)
		}
		return syscall.EINVAL
	}
	return nil
}

type Strategy struct {
}

func (s Strategy) Dup2(newFile *FileLevelWriter, oldfd *os.File) (err error) {
	if err := setStdHandle(syscall.STD_ERROR_HANDLE, syscall.Handle(newFile.File.Fd())); err != nil {
		return err
	}
	// SetStdHandle does not affect prior references to stde
	//os.Stderr = newFile.File
	//os.Stdout = newFile.File
	return nil
}
