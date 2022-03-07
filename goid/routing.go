package goid

// LocalStorage provides goroutine-local variables.
type LocalStorage interface {

	// Get returns the value in the current goroutine's local storage, if it was set before.
	Get() (value any)

	// Set copy the value into the current goroutine's local storage, and return the old value.
	Set(value any) (oldValue any)

	// Del delete the value from the current goroutine's local storage, and return it.
	Del() (oldValue any)

	// Clear delete values from all goroutine's local storages.
	Clear()
}

// ImmutableContext represents all local storages of one goroutine.
type ImmutableContext struct {
	gid    int64
	values map[uintptr]any
}

// Go start an new goroutine, and copy all local storages from current goroutine.
func Go(f func()) {
	ic := BackupContext()
	go func() {
		InheritContext(ic)
		f()
	}()
}

// BackupContext copy all local storages into an ImmutableContext instance.
func BackupContext() *ImmutableContext {
	s := loadCurrentStore()
	data := make(map[uintptr]any, len(s.values))
	for k, v := range s.values {
		data[k] = v
	}
	return &ImmutableContext{gid: s.gid, values: data}
}

// InheritContext load the specified ImmutableContext instance into the local storage of current goroutine.
func InheritContext(ic *ImmutableContext) {
	if ic == nil || ic.values == nil {
		return
	}
	s := loadCurrentStore()
	for k, v := range ic.values {
		s.values[k] = v
	}
}

// NewLocalStorage create and return an new LocalStorage instance.
func NewLocalStorage() LocalStorage {
	t := new(storage)
	t.Clear()
	return t
}

// Goid return the current goroutine's unique id.
// It will try get gid by native cgo/asm for better performance,
// and could parse gid from stack for failover supporting.
func Goid() (id int64) {
	return getGoidByStack()
}

// AllGoids return all goroutine's goid in the current golang process.
// It will try load all goid from runtime natively for better performance,
// and fallover to runtime.Stack, which is realy inefficient.
func AllGoids() (ids []int64) {
	return getAllGoidByStack()
}
