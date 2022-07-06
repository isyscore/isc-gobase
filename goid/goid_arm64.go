//go:build arm64 && gc

package goid

// Backdoor access to runtime·getg().
func getg() *g

func NativeGoid() int64 {
	return getg().goid
}
