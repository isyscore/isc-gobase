package isc

type IntRange struct {
	Start int
	End   int
}

func MakeIntRange(AStart int, AEnd int) IntRange {
	return IntRange{
		Start: AStart,
		End:   AEnd,
	}
}
