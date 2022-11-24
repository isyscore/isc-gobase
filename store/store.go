package store

import (
	"github.com/isyscore/isc-gobase/goid"
)

var RequestStorage goid.LocalStorage
var MdcStorage goid.LocalStorage

func init() {
	RequestStorage = goid.NewLocalStorage()
	MdcStorage = goid.NewLocalStorage()
}
