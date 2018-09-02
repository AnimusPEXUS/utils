package archive

import "strings"

// TODO: I don't like this package. now it is somewhat bogus

func DetermineCompressorByFilename(file_name string) (bool, string) {

	if strings.HasSuffix(file_name, ".lzma") {
		return true, "xz"
	}

	if strings.HasSuffix(file_name, ".bz2") {
		return true, "bzip2"
	}

	if strings.HasSuffix(file_name, ".gz") {
		return true, "gzip"
	}

	if strings.HasSuffix(file_name, ".xz") {
		return true, "xz"
	}

	return false, ""
}

func DetermineExtensionByFilename(file_name string) (bool, string) {
	for _, i := range []string{
		".lzma",
		".bz2",
		".gz",
		".xz",
	} {
		if strings.HasSuffix(file_name, i) {
			return true, i
		}
	}
	return false, ""
}
