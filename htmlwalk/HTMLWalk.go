package htmlwalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/AnimusPEXUS/utils/cache01"
	"github.com/AnimusPEXUS/utils/filetools"
	"github.com/AnimusPEXUS/utils/logger"
	"github.com/AnimusPEXUS/utils/set"
	"github.com/antchfx/xquery/html"
)

var _ filetools.WalkerI = &HTMLWalk{}

type HTMLWalk struct {
	scheme string
	host   string

	cache *cache01.CacheDir

	log *logger.Logger

	exclude_paths []string
}

func NewHTMLWalk(
	scheme string,
	host string,
	cache *cache01.CacheDir,
	log *logger.Logger,
	exclude_paths []string,
) (*HTMLWalk, error) {
	self := new(HTMLWalk)
	self.scheme = scheme
	self.host = host

	self.cache = cache

	self.log = log

	self.exclude_paths = exclude_paths
	return self, nil
}

func (self *HTMLWalk) LogI(txt string) {
	if self.log != nil {
		self.log.Info(txt)
	}
}

func (self *HTMLWalk) LogE(txt string) {
	if self.log != nil {
		self.log.Error(txt)
	}
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

	if !strings.HasPrefix(strconv.Itoa(r.StatusCode), "2") {
		return []os.FileInfo{}, []os.FileInfo{},
			fmt.Errorf("http code %d", r.StatusCode)
	}

	b := new(bytes.Buffer)
	io.Copy(b, r.Body)
	r.Body.Close()

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

				// {
				// 	u, err := url.Parse(j.Val)
				// 	if err == nil {
				// 		fmt.Println("is abs?:", j.Val, u.IsAbs())
				// 	}
				// }

				if u, err := url.Parse(j.Val); err == nil && (u.Host != "" ||
					u.Scheme != "" ||
					u.RawQuery != "") {
					continue searching
				}

				for _, i := range []string{"/", "#"} {
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

				// ue, err := url.PathUnescape(i.Data)
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

func (self *HTMLWalk) ListDir(pth string) (
	[]os.FileInfo,
	[]os.FileInfo,
	error,
) {

	type Container struct {
		Dirs, Files []*FileInfoForMarshal
	}

	cf := func() ([]byte, error) {

		self.LogI(
			fmt.Sprintf(
				"updating cache %s",
				(&url.URL{
					Scheme: self.scheme,
					Host:   self.host,
					Path:   pth,
				}).String(),
			),
		)

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
	}

	var res []byte

	if self.cache != nil {
		c, err := self.cache.Cache(pth, cf)
		if err != nil {
			return []os.FileInfo{}, []os.FileInfo{}, err
		}

		res, err = c.GetValue()
		if err != nil {
			return []os.FileInfo{}, []os.FileInfo{}, err
		}
	} else {
		var err error
		res, err = cf()
		if err != nil {
			return []os.FileInfo{}, []os.FileInfo{}, err
		}
	}

	co := &Container{}
	err := json.Unmarshal(res, co)
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
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
		found := false
		for _, k := range self.exclude_paths {
			m, err := regexp.MatchString(k, j)
			if err != nil {
				return err
			}
			if m {
				found = true
				break
			}
		}
		if !found {
			err = self.Walk(j, target)
			if err != nil {
				return err
			}
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
