package martinez

type edgeType int

const (
	normal              edgeType = 0
	nonContributing              = 1
	sameTransition               = 2
	differentTransition          = 3
)
