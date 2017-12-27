package cycleclassification

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
		return "unknown"
	case Development:
		return "development"
	case PreAlpha:
		return "prealpha"
	case Alpha:
		return "alpha"
	case Beta:
		return "beta"
	case RC:
		return "RC"
	case RTM:
		return "RTM"
	case GA:
		return "GA"
	case Gold:
		return "Gold"
	}
}
