//go:build (amd64 || amd64p32) && gc

#include "go_asm.h"
#include "textflag.h"

// func NativeGoid() int64
TEXT ·NativeGoid(SB),NOSPLIT,$0-8
MOVQ (TLS), R14
MOVQ g_goid(R14), R13
MOVQ R13, ret+0(FP)
RET
