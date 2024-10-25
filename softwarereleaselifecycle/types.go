package softwarereleaselifecycle

// https://en.wikipedia.org/wiki/Software_release_life_cycle

type DevelopmentStatus uint

const (
	Development DevelopmentStatus = iota
	PreAlpha
	Alpha
	Beta
	RC
	RTM
	GA
	Gold
)

func (self DevelopmentStatus) String() string {
	switch self {
	default:
		panic("invalid value")
	case Development:
		return "development"
	case PreAlpha:
		return "prealpha"
	case Alpha:
		return "alpha"
	case Beta:
		return "beta"
	case RC:
		return "RC" // Release candidate
	case RTM:
		return "RTM" // Release to manufacturing
	case GA:
		return "GA" // General availability
	case Gold:
		return "Gold"
	}
}
