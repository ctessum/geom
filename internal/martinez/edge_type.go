package martinez

type edgeType int

const (
	NORMAL               edgeType = 0
	NON_CONTRIBUTING              = 1
	SAME_TRANSITION               = 2
	DIFFERENT_TRANSITION          = 3
)
