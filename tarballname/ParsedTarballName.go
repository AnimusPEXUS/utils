package tarballname

import (
	"fmt"

	"github.com/AnimusPEXUS/utils/versionorstatus"
)

type ParsedTarballName struct {
	Basename           string
	Name               string
	Version            *versionorstatus.ParsedVersionOrStatus
	Status             *versionorstatus.ParsedVersionOrStatus
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

func (self *ParsedTarballName) Render(apply_extension bool) (string, error) {
	name := ""
	if self.Name != "" {
		name = self.Name + "-"
	}
	status := ""
	if self.Status.StrSliceString("") != "" {
		status = "-" + self.Status.StrSliceString("")
	}

	ext := ""
	if apply_extension {
		if self.Extension != "" {
			ext = self.Extension
		}
	}

	version_str, err := self.Version.IntSliceString(".")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%s%s", name, version_str, status, ext), nil
}
