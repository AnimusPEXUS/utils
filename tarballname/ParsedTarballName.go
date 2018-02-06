package tarballname

import (
	"fmt"
)

type ParsedTarballName struct {
	Basename           string
	Name               string
	Version            ParsedVersion
	Status             ParsedStatus
	Extension          string
	OriginalInputValue string
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

func (self *ParsedTarballName) Render(apply_extension bool) string {
	name := ""
	if self.Name != "" {
		name = self.Name + "-"
	}
	status := ""
	if self.Status.Str != "" {
		status = "-" + self.Status.Str
	}

	ext := ""
	if apply_extension {
		if self.Extension != "" {
			ext = self.Extension
		}
	}

	return fmt.Sprintf("%s%s%s%s", name, self.Version.Str, status, ext)
}
