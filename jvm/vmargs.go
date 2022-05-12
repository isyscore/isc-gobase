package jvm

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

/*
#include <stdlib.h>
#include <jni.h>
*/
import "C"

func ArgumentsCheck(args ...any) ([]string, error) {
	var types []string
	var err error
	for _, a := range args {
		t := reflect.TypeOf(a).String()
		if t == "bool" {
			types = append(types, "Z")
		} else if t == "uint8" {
			types = append(types, "B")
		} else if t == "int16" {
			types = append(types, "S")
		} else if t == "int64" {
			types = append(types, "J")
		} else if t == "float32" {
			types = append(types, "F")
		} else if t == "float64" {
			types = append(types, "D")
		} else if t == "string" {
			types = append(types, "Ljava/lang/String;")
		} else if strings.Contains(t, "int") {
			types = append(types, "I")
		} else {
			err = fmt.Errorf("unsupported type: %s", t)
		}
	}
	return types, err
}

func ParseArguments(types []string, retType string, args ...any) (string, **C.char, *unsafe.Pointer) {
	size := C.size_t(unsafe.Sizeof((*C.char)(nil)))
	clen := C.size_t(len(args))
	typesArg := (**C.char)(C.malloc(size * clen))
	typesArgView := (*[1 << 30]*C.char)(unsafe.Pointer(typesArg))[0:len(args):len(args)]
	sizev := C.size_t(unsafe.Sizeof((*C.void)(nil)))
	argArg := (*unsafe.Pointer)(C.malloc(sizev * clen))
	argArgView := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(argArg))[0:len(args):len(args)]

	sigStr := "("
	for i, t := range types {
		sigStr += t
		typesArgView[i] = C.CString(t)
		if t == "Ljava/lang/String;" {
			argArgView[i] = unsafe.Pointer(C.CString(args[i].(string)))
		} else if t == "I" {
			ci := C.int(args[i].(int))
			argArgView[i] = unsafe.Pointer(&ci)
		} else if t == "J" {
			li := C.long(args[i].(int64))
			argArgView[i] = unsafe.Pointer(&li)
		} else if t == "Z" {
			var bi C.int
			if args[i].(bool) {
				bi = 1
			} else {
				bi = 0
			}
			argArgView[i] = unsafe.Pointer(&bi)
		} else if t == "S" {
			si := C.short(args[i].(int16))
			argArgView[i] = unsafe.Pointer(&si)
		} else if t == "F" {
			fi := C.float(args[i].(float32))
			argArgView[i] = unsafe.Pointer(&fi)
		} else if t == "D" {
			di := C.double(args[i].(float64))
			argArgView[i] = unsafe.Pointer(&di)
		} else if t == "B" {
			bi := C.uchar(args[i].(uint8))
			argArgView[i] = unsafe.Pointer(&bi)
		}
	}
	sigStr += ")" + retType
	return sigStr, typesArg, argArg
}

func FreeArgs(size int, types []string, typesArg **C.char, valArgs *unsafe.Pointer) {
	typView := (*[1 << 30]*C.char)(unsafe.Pointer(typesArg))[0:size:size]
	valView := (*[1 << 30]unsafe.Pointer)(unsafe.Pointer(valArgs))[0:size:size]
	for i := 0; i < size; i++ {
		C.free(unsafe.Pointer(typView[i]))
		if types[i] == "Ljava/lang/String;" {
			C.free(valView[i])
		}
	}
	C.free(unsafe.Pointer(typesArg))
	C.free(unsafe.Pointer(valArgs))
}

func ParseNameSig(name string, sig string) (*C.char, *C.char) {
	nameC := C.CString(name)
	sigC := C.CString(sig)
	return nameC, sigC
}

func FreeNameSig(nameC *C.char, sigC *C.char) {
	C.free(unsafe.Pointer(nameC))
	C.free(unsafe.Pointer(sigC))
}
