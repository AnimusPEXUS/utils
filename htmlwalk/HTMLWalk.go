package htmlwalk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/AnimusPEXUS/utils/cache01"
	"github.com/AnimusPEXUS/utils/directory"
	"github.com/AnimusPEXUS/utils/filetools"
	"github.com/AnimusPEXUS/utils/set"
	"github.com/antchfx/xquery/html"
)

var _ filetools.WalkerI = &HTMLWalk{}

type HTMLWalk struct {
	scheme string
	host   string

	cache *cache01.CacheDir

	tree *directory.File
}

func NewHTMLWalk(
	scheme string,
	host string,
	cache *cache01.CacheDir,
) (*HTMLWalk, error) {
	ret := new(HTMLWalk)
	ret.scheme = scheme
	ret.host = host

	ret.cache = cache

	ret.tree = directory.NewFile(nil, "", true, nil)
	return ret, nil
}

func (self *HTMLWalk) ListDir(pth string) (
	[]os.FileInfo,
	[]os.FileInfo,
	error,
) {

	type Container struct {
		Dirs, Files []*FileInfoForMarshal
	}

	c, err := self.cache.Cache(
		pth,
		func() ([]byte, error) {

			d, f, err := self.ListDirNotCached(pth)
			if err != nil {
				return []byte{}, err
			}

			c := &Container{}

			for _, i := range d {
				c.Dirs = append(c.Dirs, NewFileInfoForMarshal(i))
			}
			for _, i := range f {
				c.Files = append(c.Files, NewFileInfoForMarshal(i))
			}

			ret, err := json.Marshal(c)
			if err != nil {
				return []byte{}, err
			}

			return ret, nil
		},
	)
	if err != nil {
		t := make([]os.FileInfo, 0)
		return t, t, err
	}

	res, err := c.GetValue()
	if err != nil {
		t := make([]os.FileInfo, 0)
		return t, t, err
	}

	co := &Container{}
	err = json.Unmarshal(res, co)
	if err != nil {
		t := make([]os.FileInfo, 0)
		return t, t, err
	}

	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)

	for _, i := range co.Dirs {
		dirs = append(dirs, i.GetFileInfo())
	}

	for _, i := range co.Files {
		files = append(files, i.GetFileInfo())
	}

	sort.Sort(OsFileInfoSort(dirs))
	sort.Sort(OsFileInfoSort(files))

	return dirs, files, nil

}

func (self *HTMLWalk) ListDirNotCached(pth string) (
	[]os.FileInfo,
	[]os.FileInfo,
	error,
) {

	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)

	pth = path.Clean(pth)

	u := &url.URL{
		Scheme: self.scheme,
		Host:   self.host,
		Path:   pth,
	}

	r, err := http.Get(u.String())
	if err != nil {
		return dirs, files, err
	}

	b := new(bytes.Buffer)

	r.Write(b)

	doc, err := htmlquery.Parse(b)

	if err != nil {
		return dirs, files, err
	}

	res := htmlquery.Find(doc, ".//a")

	s := set.NewSetString()

searching:
	for _, i := range res {
		for _, j := range i.Attr {
			if j.Key == "href" {

				if u, err := url.Parse(j.Val); err == nil && (u.Host != "" ||
					u.Scheme != "" ||
					u.RawQuery != "" ||
					u.IsAbs()) {
					continue searching
				}

				for _, i := range []string{"#"} {
					if strings.HasPrefix(j.Val, i) {
						continue searching
					}
				}

				{
					c := path.Clean(j.Val)
					for _, i := range []string{".", "..", ""} {
						if c == i {
							continue searching
						}
					}
				}

				ue, err := url.PathUnescape(j.Val)
				if err != nil {
					continue searching
				}
				s.Add(ue)

				continue searching
			}
		}
	}

	for _, i := range s.ListStrings() {
		if strings.HasSuffix(i, "/") {
			t := &FileInfo{name: i, isdir: true}
			dirs = append(dirs, t)
		} else {
			t := &FileInfo{name: i, isdir: false}
			files = append(files, t)
		}
	}

	sort.Sort(OsFileInfoSort(dirs))
	sort.Sort(OsFileInfoSort(files))

	return dirs, files, nil
}

func (self *HTMLWalk) Walk(
	pth string,
	target func(
		dir string,
		dirs []os.FileInfo,
		files []os.FileInfo,
	) error,
) error {

	dirs, files, err := self.ListDir(pth)
	if err != nil {
		return err
	}

	err = target(pth, dirs, files)
	if err != nil {
		return err
	}

	for _, i := range dirs {
		j := path.Join(pth, i.Name())
		err = self.Walk(j, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *HTMLWalk) Tree(pth string) (map[string]os.FileInfo, error) {

	ret := make(map[string]os.FileInfo)

	err := self.Walk(
		pth,
		func(
			dir string,
			dirs []os.FileInfo,
			files []os.FileInfo,
		) error {
			for _, i := range files {
				j := path.Join(dir, i.Name())
				ret[j] = i
			}
			for _, i := range dirs {
				j := path.Join(dir, i.Name())
				j += "/"
				ret[j] = i
			}
			return nil
		},
	)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
