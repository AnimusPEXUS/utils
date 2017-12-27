package tarballname

import (
	"fmt"
)

type ParsedTarballName struct {
	Basename  string
	Name      string
	Version   ParsedVersion
	Status    ParsedStatus
	Extension string
}

func (self *ParsedTarballName) InfoText() string {
	ret := fmt.Sprintf(` Basename:  "%s"
 Name:      "%s"
 Extension: "%s"
 Version:
%s
 Status:
%s
`,
		self.Basename,
		self.Name,
		self.Extension,
		self.Version.InfoText(),
		self.Status.InfoText(),
	)

	return ret
}
