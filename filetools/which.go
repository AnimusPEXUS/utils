package filetools

import (
	"errors"
	"os"
	"path"
	"strings"
)

func Which(executable_name string, under []string) (string, error) {
	// TODO: is this function should have switch to not return nonexecutable files
	// , return_non_executable_files bool
	var where_to_search []string

	if len(under) != 0 {
		under_plus := []string{}

		for _, i := range under {
			if !strings.HasSuffix(i, "/bin") && !strings.HasSuffix(i, "/sbin") {
				under_plus = append(under_plus, path.Join(i, "/bin"))
				under_plus = append(under_plus, path.Join(i, "/sbin"))
			}
		}

		where_to_search = make([]string, 0)
		where_to_search = append(where_to_search, under...)
		where_to_search = append(where_to_search, under_plus...)
	} else {
		splitted_path := []string{}
		env := os.Environ()
		PATH, ok := env["PATH"]
		if !ok {
			return "", errors.New("no PATH environment variable")
		}
		splitted_path = strings.Split(PATH, ":")
		where_to_search = splitted_path
	}

	{
		t := make([]string, 0)
		for _, i := range where_to_search {
			s, err := os.Stat(i)
			if err != nil {
				if !os.IsNotExist(err) {
					return "", err
				} else {
					continue
				}
			}

			if !s.IsDir() {
				continue
			}

			t = append(t, i)
		}

		where_to_search = t
	}

	for _, i := range where_to_search {
		p := path.Join(i, executable_name)

		s, err := os.Stat(p)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			} else {
				continue
			}
		}

		if s.IsDir() {
			continue
		}

		return p, nil
	}

	return "", errors.New("not found")
}
