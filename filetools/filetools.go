package filetools

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
)

type WalkerI interface {
	ListDir(pth string) ([]os.FileInfo, []os.FileInfo, error)
	Walk(
		pth string,
		target func(
			dir string,
			dirs []os.FileInfo,
			files []os.FileInfo,
		) error,
	) error
	Tree(pth string) (map[string]os.FileInfo, error)
}

func ListDir(pth string) ([]os.FileInfo, []os.FileInfo, error) {
	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)

	data, err := ioutil.ReadDir(pth)
	if err != nil {
		return dirs, files, nil
	}

	for _, i := range data {
		if i.Mode()&os.ModeSymlink != 0 {
			files = append(files, i)
		} else {
			if i.IsDir() {
				dirs = append(dirs, i)
			} else {
				files = append(files, i)
			}
		}
	}

	{
		names := make([]string, 0)
		for _, i := range dirs {
			names = append(names, i.Name())
		}
		sort.Strings(names)

		new_structs := make([]os.FileInfo, 0)

		for _, i := range names {
			for _, j := range dirs {
				if j.Name() == i {
					new_structs = append(new_structs, j)
				}
			}
		}

		dirs = new_structs
	}

	{
		names := make([]string, 0)
		for _, i := range files {
			names = append(names, i.Name())
		}
		sort.Strings(names)

		new_structs := make([]os.FileInfo, 0)

		for _, i := range names {
			for _, j := range files {
				if j.Name() == i {
					new_structs = append(new_structs, j)
				}
			}
		}

		files = new_structs
	}

	return dirs, files, nil
}

func Walk(
	pth string,
	target func(
		dir string,
		dirs []os.FileInfo,
		files []os.FileInfo,
	) error,
) error {

	pth, err := filepath.Abs(pth)
	if err != nil {
		return err
	}

	{
		pths, err := os.Stat(pth)
		if err != nil {
			return err
		}
		if !pths.IsDir() {
			return errors.New("pth must be dir")
		}
	}

	all_dirs := make([]string, 0)

	all_dirs = append(all_dirs, pth)

	// lower_i := 0
	upper_i := len(all_dirs)

	for {
		// fmt.Println("all_dirs", all_dirs)

		if upper_i == 0 {
			// fmt.Println("break")
			break
		}

		for i := 0; i != upper_i; i++ {
			all_dirs_i := all_dirs[i]

			dirs, files, err := ListDir(all_dirs_i)
			if err != nil {
				return err
			}
			// fmt.Println("dirs, files, err", dirs, files, err)

			for _, j := range dirs {
				all_dirs = append(all_dirs, path.Join(all_dirs_i, j.Name()))
			}

			err = target(all_dirs_i, dirs, files)
			if err != nil {
				return err
			}

		}

		all_dirs = append(all_dirs[upper_i:])

		// lower_i = 0
		upper_i = len(all_dirs)

	}

	return nil

}
