package types

// https://en.wikipedia.org/wiki/Software_release_life_cycle

type StabilityClassification uint

const (
	Development StabilityClassification = iota
	PreAlpha
	Alpha
	Beta
	RC
	Release
)

func (self StabilityClassification) String() string {
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
		return "rc"
	case Release:
		return "release"
	}
}

func IsStable(value StabilityClassification) bool {
	return value >= Release
}
