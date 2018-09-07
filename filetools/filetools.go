package filetools

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/AnimusPEXUS/utils/logger"
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

// TODO: no_symlink_delve option or somethin similar. symleanks should be dealth with
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

	upper_i := len(all_dirs)

	for {

		if upper_i == 0 {
			break
		}

		for i := 0; i != upper_i; i++ {
			all_dirs_i := all_dirs[i]

			dirs, files, err := ListDir(all_dirs_i)
			if err != nil {
				return err
			}

			for _, j := range dirs {
				all_dirs = append(all_dirs, path.Join(all_dirs_i, j.Name()))
			}

			err = target(all_dirs_i, dirs, files)
			if err != nil {
				return err
			}

		}

		all_dirs = append(all_dirs[upper_i:])

		upper_i = len(all_dirs)

	}

	return nil
}

// Walks inside src_path and copying files to dst_path
func CopyTree(
	src_path string,
	dst_path string,
	dst_not_empty_is_error bool,
	clear_dst_dir bool,
	files_already_exists_is_error bool,
	overwrite_existing_files bool,
	log logger.LoggerI,
	verbose_log bool,
	pass_log bool,
	copy_file_cb func(src, dst string, log logger.LoggerI) error,
) error {

	logi := func(t string, for_verbose_mode bool) {
		if !verbose_log && !for_verbose_mode {
			return
		}
		if log != nil {
			log.Info(t)
		}
	}

	loge := func(t interface{}, for_verbose_mode bool) {
		if !verbose_log && !for_verbose_mode {
			return
		}
		if log != nil {
			log.Error(t)
		}
	}

	src_path_stat, err := os.Lstat(src_path)
	if err != nil {
		loge(err, false)
		return err
	}

	if Is(src_path_stat.Mode()).Symlink() || (!src_path_stat.IsDir()) {

		logi("copying file or symlink "+src_path+" as "+dst_path, false)

		dst_exists := false

		dst_path_stat, err := os.Stat(dst_path)
		if err != nil {
			if !os.IsNotExist(err) {
				loge(err, false)
				return err
			}
		} else {
			dst_exists = true
		}

		if dst_exists {

			if files_already_exists_is_error {
				err = errors.New("dst file already exists")
				loge(err, false)
				return err
			}

			if overwrite_existing_files {
				logi("removing destination before overwritting", true)
				err = os.RemoveAll(dst_path)
				if err != nil {
					loge(err, false)
					return err
				}
			} else {
				logi("existing destination should not be overwritten", true)
				if dst_path_stat.IsDir() {
					err = errors.New("dst is dir and src is not dir")
					loge(err, false)
					return err
				}
				return nil
			}
		}

		{
			var l logger.LoggerI
			if pass_log {
				l = log
			}
			os.MkdirAll(path.Dir(dst_path), 0700)
			err = copy_file_cb(src_path, dst_path, l)
			if err != nil {
				loge(err, false)
				return err
			}
		}

	} else {

		logi("copying directory "+src_path+" as "+dst_path, false)

		dst_path_stat, err := os.Stat(dst_path)
		if err != nil {
			if !os.IsNotExist(err) {
				// some other error: we can't (and shouldn't) handle it here
				loge(err, false)
				return err
			}
		} else {

			if !dst_path_stat.IsDir() {
				err = errors.New("destination exists but it's not directory")
				loge(err, false)
				return err
			}

			// if dst already exists, we should decide what to do if it's not empty

			dst_dir_lst, err := ioutil.ReadDir(dst_path)
			if err != nil {
				loge(err, false)
				return err
			}

			if len(dst_dir_lst) != 0 {
				if dst_not_empty_is_error {
					err = errors.New("destination directory is not empty")
					loge(err, false)
					return err
				}

				if clear_dst_dir {

					for _, i := range dst_dir_lst {
						i_joined := path.Join(dst_path, i.Name())
						if i.IsDir() {
							err = os.RemoveAll(i_joined)
							if err != nil {
								loge(err, false)
								return err
							}
						} else {
							err = os.Remove(i_joined)
							if err != nil {
								loge(err, false)
								return err
							}
						}
					}
				}
			}
		}

		logi("making "+dst_path, true)
		err = os.MkdirAll(dst_path, 0700)
		if err != nil {
			return err
		}

		err = Walk(
			src_path,
			func(
				dir string,
				dirs []os.FileInfo,
				files []os.FileInfo,
			) error {
				logi("working inside "+dir, true)

				dir_rel_part, err := filepath.Rel(src_path, dir)
				if err != nil {
					loge(err, false)
					return err
				}

				dst_path := path.Join(dst_path, dir_rel_part)

				err = os.MkdirAll(dst_path, 0700)
				if err != nil {
					loge(err, false)
					return err
				}

				for _, i := range files {

					src_file_path := path.Join(dir, i.Name())
					dst_file_path := path.Join(dst_path, i.Name())

					dst_file_path_exists := false
					dst_file_path_stat, err := os.Lstat(dst_file_path)
					if err != nil {
						if !os.IsNotExist(err) {
							loge(err, false)
							return err
						} else {
							dst_file_path_exists = false
						}
					} else {
						dst_file_path_exists = true
					}

					if dst_file_path_exists {
						if dst_file_path_stat.IsDir() {
							err = errors.New("destination file already exists and it is directory " + dst_file_path)
							loge(err, false)
							return err
						}
						if files_already_exists_is_error {
							err = errors.New("destination file already exists")
							loge(err, false)
							return err
						}
					}

					if !dst_file_path_exists || (dst_file_path_exists && overwrite_existing_files) {

						var l logger.LoggerI
						if pass_log {
							l = log
						}

						err = copy_file_cb(src_file_path, dst_file_path, l)
						if err != nil {
							loge(err, false)
							return err
						}
					}
				}
				return nil
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyWithOptions(
	src string,
	dst string,
	log logger.LoggerI,
	copyinfo bool,
	dereference_symlinks bool,
) error {

	var src_stat os.FileInfo
	var err error

	if !dereference_symlinks {
		src_stat, err = os.Lstat(src)
		if err != nil {
			return err
		}
	} else {
		src_stat, err = os.Stat(src)
		if err != nil {
			return err
		}
	}

	if Is(src_stat.Mode()).Symlink() {

		link_value, err := os.Readlink(src)
		if err != nil {
			return err
		}

		os.Remove(dst)

		err = os.Symlink(link_value, dst)
		if err != nil {
			return err
		}

	} else {

		sf, err := os.Open(src)
		if err != nil {
			return err
		}
		defer sf.Close()

		err = os.Remove(dst)
		if err != nil {
			return err
		}

		df, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer df.Close()

		_, err = io.Copy(df, sf)
		if err != nil {
			return err
		}

		if copyinfo {
			sfs, err := os.Stat(src)
			if err != nil {
				return err
			}

			err = os.Chmod(dst, sfs.Mode())
			if err != nil {
				return err
			}

			err = os.Chtimes(dst, sfs.ModTime(), sfs.ModTime())
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func CopyWithInfo(src, dst string, log logger.LoggerI) error {
	return CopyWithOptions(src, dst, log, true, false)
}

type Is os.FileMode

func (self Is) Symlink() bool {
	return (os.FileMode(self) & os.ModeSymlink) != 0
}

func (self Is) Regular() bool {
	return os.FileMode(self).IsRegular()
}
