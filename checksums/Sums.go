package checksums

import (
	"bufio"
	"io"
	"os"

	mystrings "github.com/AnimusPEXUS/utils/strings"
)

type Sums struct {
	Lines []*SumsLine
}

func NewSums() *Sums {
	ret := new(Sums)
	ret.Lines = make([]*SumsLine, 0)
	return ret
}

func (self *Sums) LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	file_eof := false

	for {
		if file_eof {
			break
		}
		prefix_line := make([]byte, 0)

	continue_readline:
		line, isprefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				file_eof = true
				goto line_read
			} else {
				return err
			}
		}
		if isprefix {
			prefix_line = append(prefix_line, line...)
			goto continue_readline
		}
	line_read:

		line = append(prefix_line, line...)
		prefix_line = nil

		line_str := string(line)
		line = nil

		if len(line_str) == 0 || mystrings.StringOfSpaces(line_str) {
			continue
		}

		sums_file_line, err := NewSumsLineFromString(line_str)
		if err != nil {
			return err
		}

		line_str = ""

		self.Lines = append(self.Lines, sums_file_line)
	}

	return nil
}

func (self *Sums) SaveToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	for _, i := range self.Lines {
		_, err := f.WriteString(i.String() + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// warning: this may be slow for large list. working with 1000 takes few seconds
func (self *Sums) SortLinesByValue() {
	len_self_lines := len(self.Lines)

	if len_self_lines < 2 {
		return
	}

	for i := 0; i != len_self_lines-1; i++ {

		for j := i + 1; j != len_self_lines; j++ {
			if self.Lines[i].value > self.Lines[j].value {
				z := self.Lines[i]
				self.Lines[i] = self.Lines[j]
				self.Lines[j] = z
			}
		}
	}
	return
}

// warning: this may be slow for large list. working with 1000 takes few seconds
func (self *Sums) RemoveDuplicates() {

	if len(self.Lines) < 2 {
		return
	}

	for i := 0; i != len(self.Lines); i++ {
		for j := len(self.Lines) - 1; j != i; j-- {
			if self.Lines[i].IsEqual(self.Lines[j]) {
				self.Lines = append(self.Lines[:j], self.Lines[j+1:]...)
			}
		}
	}

	return
}
