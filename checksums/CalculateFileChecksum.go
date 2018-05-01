package checksums

import (
	"hash"
	"io"
	"os"
)

func CalculateFileChecksum(
	filename string,
	new_hash hash.Hash,
	bs int,
) ([]byte, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buff := make([]byte, bs)
	for {
		s, err := f.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		new_hash.Write(buff[:s])
	}

	ret := new_hash.Sum(nil)

	return ret, nil
}

// func CalculateDirChecksums(
// 	dirname string,
// 	output_filename string,
// 	make_rel bool,
// 	rel_to string,
// 	conv_to_rooted bool,
// 	exclude []string,
// 	bs int,
// 	create_new_hash_cb func() (hash.Hash, error),
// 	log *logger.Logger, // can be nil, to disable
// ) error {
//
// 	if s, err := os.Stat(dirname); err != nil {
// 		return err
// 	} else {
// 		if !s.IsDir() {
// 			return errors.New("dirname not a deirectory")
// 		}
// 	}
//
// 	dirname, err := filepath.Abs(dirname)
// 	if err != nil {
// 		return err
// 	}
//
// 	output_fileobj, err := os.Create(output_filename)
// 	if err != nil {
// 		return err
// 	}
//
// 	defer output_fileobj.Close()
//
// 	if !make_rel {
// 		rel_to = dirname
// 	}
//
// 	rel_to, err = filepath.Abs(rel_to)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = filetools.Walk(
// 		dirname,
// 		func(
// 			root string,
// 			dirs []os.FileInfo,
// 			files []os.FileInfo,
// 		) error {
//
// 			for _, f := range files {
//
// 				root_f := path.Join(root, path.Base(f.Name()))
//
// 				for _, i := range exclude {
// 					if root_f == i {
// 						continue
// 					}
// 				}
//
// 				rel_path, err := filepath.Rel(rel_to, root_f)
// 				if err != nil {
// 					return err
// 				}
//
// 				s, err := os.Lstat(root_f)
// 				if err != nil {
// 					return err
// 				}
//
// 				// TODO: excluding symlincs is questionable.
// 				if s.IsDir() || (os.ModeSymlink&s.Mode() != 0) {
// 					continue
// 				}
//
// 				new_hash, err := create_new_hash_cb()
// 				if err != nil {
// 					return err
// 				}
//
// 				sum, err := CalculateFileChecksum(root_f, new_hash, bs)
// 				if err != nil {
// 					return err
// 				}
//
// 				wfn := rel_path
// 				if conv_to_rooted && !strings.HasPrefix(wfn, string(filepath.Separator)) {
// 					wfn = string(filepath.Separator) + wfn
// 				}
//
// 				_, err = output_fileobj.WriteString(
// 					fmt.Sprintf(
// 						"%s\n", checksums.NewSumsLine(sum, wfn).String(),
// 					),
// 				)
// 				if err != nil {
// 					return err
// 				}
//
// 			}
//
// 			return nil
// 		},
// 	)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
