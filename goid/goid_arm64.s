//go:build arm64 && gc

#include "textflag.h"

// func getg() *g
TEXT ·getg(SB),NOSPLIT,$0-8
	MOVD g, ret+0(FP)
	RET
