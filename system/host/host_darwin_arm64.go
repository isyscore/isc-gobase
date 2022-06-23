//go:build darwin && arm64

package host

type Utmpx struct {
	User [256]int8
	Id   [4]int8
	Line [32]int8
	Pid  int32
	Type int16
	Tv   Timeval
	Host [256]int8
	Pad  [16]uint32
}
type Timeval struct {
	Sec       int64
	Usec      int32
	Pad_cgo_0 [4]byte
}
