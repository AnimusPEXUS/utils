package types

import "strconv"

type Status uint

type TarballNameParserI interface {
	ParseName(value string) (*ParseResult, error)
}

const (
	Development Status = iota
	PreAlpha
	Alpha
	Beta
	RC
	RTM
	GA
	Gold
)

func (self Status) String() string {
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

type ParseResult struct {
	Name        string
	HaveVersion bool
	Version     []uint
	HaveStatus  bool
	Status      Status
	HaveBuildId bool
	BuildId     string
}

func (self *ParseResult) VersionString() string {
	ret := ""
	l := len(self.Version) - 1
	for ii, i := range self.Version {
		ret += strconv.Itoa(int(i))
		if ii != l {
			ret += "."
		}
	}
	return ret
}

func (self *ParseResult) VersionStrings() []string {
	ret := make([]string, 0)
	for _, i := range self.Version {
		ret = append(ret, strconv.Itoa(int(i)))
	}
	return ret
}

func (self *ParseResult) VersionInt() []int {
	ret := make([]int, 0)
	for _, i := range self.Version {
		ret = append(ret, int(i))
	}
	return ret
}
