package htmlwalk

import (
	"bytes"
	"encoding/json"
	"errors"
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
	maxdepth      int
}

func NewHTMLWalk(
	scheme string,
	host string,
	cache *cache01.CacheDir,
	log *logger.Logger,
	exclude_paths []string,
	maxdepth int,
) (*HTMLWalk, error) {
	self := new(HTMLWalk)
	self.scheme = scheme
	self.host = host

	self.cache = cache

	self.log = log

	self.exclude_paths = exclude_paths
	self.maxdepth = maxdepth
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

	//	pth = path.Clean(pth)

	u := &url.URL{
		Scheme: self.scheme,
		Host:   self.host,
		Path:   pth + "/",
	}

	// c := new(http.Client)
	//
	// req, err := http.NewRequest("GET", u.String(), nil)
	// if err != nil {
	// 	return nil, nil, err
	// }
	//
	// req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:59.0) Gecko/20100101 Firefox/59.0")
	// req.Header.Set("Host", self.host)
	//
	// req.Header.Write(os.Stdout)
	//
	// r, err := c.Do(req)
	// if err != nil {
	// 	return nil, nil, err
	// }

	//	self.LogI("url " + u.String())

	r, err := http.Get(u.String())
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	res := htmlquery.Find(doc, ".//a")

	s := set.NewSetString()

searching:
	for _, i := range res {
		for _, j := range i.Attr {
			if j.Key == "href" {

				j_val := j.Val

				if u, err := url.Parse(j_val); err == nil && (u.Host != "" ||
					u.Scheme != "" ||
					u.RawQuery != "") {
					continue searching
				}

				is_dir := strings.HasSuffix(j_val, "/")

				if strings.HasPrefix(j_val, "#") {
					continue searching
				}

				if strings.HasPrefix(j_val, "/") {

					if is_dir {
						j_val = path.Dir(j_val)
					}

					if path.Dir(j_val) != path.Clean(pth) {
						continue
					}

					j_val = path.Base(j_val)

					if is_dir {
						j_val += "/"
					}

				}

				{
					c := path.Clean(j_val)
					for _, i := range []string{".", "..", ""} {
						if c == i {
							continue searching
						}
					}
				}

				// ue, err := url.PathUnescape(i.Data)
				ue, err := url.PathUnescape(j_val)
				if err != nil {
					continue searching
				}
				s.Add(ue)

				continue searching
			}
		}
	}

	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)

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
			return nil, err
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
			return nil, err
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
	return self._Walk(pth, target, self.maxdepth)
}

func (self *HTMLWalk) _Walk(
	pth string,
	target func(
		dir string,
		dirs []os.FileInfo,
		files []os.FileInfo,
	) error,
	maxdepth int,
) error {

	maxdepth--

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

			if self.maxdepth < 0 || maxdepth > 0 {
				err = self._Walk(j, target, maxdepth)
				if err != nil {
					return err
				}
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

func (self *HTMLWalk) GetDownloadingURIForFile(
	name string,
	tree_pth string,
) (string, error) {
	name = path.Base(name)

	tree, err := self.Tree(tree_pth)
	if err != nil {
		return "", err
	}

	for k, _ := range tree {
		if path.Base(k) == name {
			u := &url.URL{
				Scheme: self.scheme,
				Host:   self.host,
				Path:   k,
			}
			return u.String(), nil
		}
	}

	return "", errors.New("not found")
}
