package rl

import (
	"context"
	"runtime/pprof"
	"strconv"
	"sync/atomic"
	"unsafe"
)

const labelTag = "github.com/isyscore/isc-gobase"

var lastID uintptr

type labelMapAndContext struct {
	lm  labelMap
	ctx context.Context
}

func Set(ctx context.Context) {
	var gctx labelMapAndContext
	pprof.SetGoroutineLabels(ctx)
	lm := (*labelMap)(runtimeGetProfLabel())
	if lm != nil && len(*lm) > 0 {
		gctx.lm = *lm
	} else {
		gctx.lm = make(labelMap)
	}
	id := atomic.AddUintptr(&lastID, 1)
	gctx.lm[labelTag] = strconv.FormatUint(uint64(id), 16)
	gctx.ctx = ctx
	runtimeSetProfLabel(unsafe.Pointer(&gctx))
}

func Get() context.Context {
	lm := (*labelMap)(runtimeGetProfLabel())
	if lm != nil {
		if _, ok := (*lm)[labelTag]; ok {
			gctx := (*labelMapAndContext)(unsafe.Pointer(lm))
			return gctx.ctx
		}
	}
	return nil
}
