package checksums

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const SumsLineRe = `(.*?) \*(.*)`

var SumsLineReC = regexp.MustCompile(SumsLineRe)

type SumsLine struct {
	sum   []byte
	value string
}

func NewSumsLine(
	sum []byte,
	value string,
) *SumsLine {
	self := new(SumsLine)
	self.sum = sum
	self.value = value
	return self
}

func NewSumsLineFromBytes(data []byte) (*SumsLine, error) {
	return NewSumsLineFromString(string(data))
}

// NOTE: trailing newline trimmed from value
func NewSumsLineFromString(data string) (*SumsLine, error) {
	data = strings.TrimSuffix(data, "\n")
	if !SumsLineReC.MatchString(data) {
		return nil, errors.New("input data not matched as sums file line")
	}
	self := new(SumsLine)
	for k, v := range SumsLineReC.FindStringSubmatch(data) {
		switch k {
		case 0:
		case 1:
			vv, err := hex.DecodeString(v)
			if err != nil {
				return nil, err
			}
			self.sum = vv

		case 2:

			self.value = v

		default:
			return nil, errors.New("input data not matched as sums file line")
		}
	}
	return self, nil
}

// NOTE: result does not ends with newline
func (self *SumsLine) String() string {
	ret := fmt.Sprintf(
		"%s *%s",
		hex.EncodeToString(self.sum),
		self.value,
	)
	return ret
}

func (self *SumsLine) IsEqual(other *SumsLine) bool {
	return self.value == other.value && bytes.Compare(self.sum, other.sum) == 0
}
