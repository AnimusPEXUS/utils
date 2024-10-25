package sfnetwalk

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
	"github.com/antchfx/htmlquery"
)

var _ filetools.WalkerI = &SFNetWalk{}

type SFNetWalk struct {
	project string
	cache   *cache01.CacheDir
	log     *logger.Logger

	exclude_paths []string
	maxdepth      int
}

func NewSFNetWalk(
	project string,
	cache *cache01.CacheDir,
	log *logger.Logger,
	exclude_paths []string,
	maxdepth int,
) (*SFNetWalk, error) {
	self := new(SFNetWalk)
	self.project = project
	self.cache = cache
	self.log = log

	self.exclude_paths = exclude_paths
	self.maxdepth = maxdepth

	return self, nil
}

func (self *SFNetWalk) LogI(txt string) {
	if self.log != nil {
		self.log.Info(txt)
	}
}

func (self *SFNetWalk) LogE(txt string) {
	if self.log != nil {
		self.log.Error(txt)
	}
}

func (self *SFNetWalk) ListDirNotCached(pth string) (
	[]os.FileInfo,
	[]os.FileInfo,
	error,
) {
	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)

	url_path := []string{"projects", self.project, "files"}
	url_path = append(url_path, strings.Split(pth, "/")...)

	u := &url.URL{
		Scheme: "https",
		Host:   "sourceforge.net",
		Path:   path.Clean(path.Join(url_path...)),
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return []os.FileInfo{}, []os.FileInfo{}, err
	}

	if !strings.HasPrefix(strconv.Itoa(resp.StatusCode), "2") {
		return []os.FileInfo{}, []os.FileInfo{},
			fmt.Errorf("http code %d", resp.StatusCode)
	}

	b := new(bytes.Buffer)
	io.Copy(b, resp.Body)
	resp.Body.Close()

	doc, err := htmlquery.Parse(b)

	file_list_table_res := htmlquery.Find(doc, `.//table[@id="files_list"]`)
	if len(file_list_table_res) != 1 {
		return []os.FileInfo{}, []os.FileInfo{},
			errors.New("found invalid number of files_list tables")
	}

	file_list_table := file_list_table_res[0]

	file_list_table_tbody_res := htmlquery.Find(file_list_table, "tbody")
	if len(file_list_table_tbody_res) != 1 {
		return []os.FileInfo{}, []os.FileInfo{},
			errors.New("found invalid number of tbody")
	}

	file_list_table_tbody := file_list_table_tbody_res[0]

	folder_trs := htmlquery.Find(file_list_table_tbody, "tr")

	for _, i := range folder_trs {
		cls := ""

		for _, i := range i.Attr {
			if i.Key == "class" {
				cls = i.Val
				break
			}
		}

		name := ""
		for _, i := range i.Attr {
			if i.Key == "title" {
				name = i.Val
				break
			}
		}
		name, err = url.PathUnescape(name)
		if err != nil {
			return []os.FileInfo{}, []os.FileInfo{}, err
		}

		if strings.Contains(cls, "folder") {
			dirs = append(dirs, &FileInfo{name: name, isdir: true})
		} else if strings.Contains(cls, "file") {
			files = append(files, &FileInfo{name: name, isdir: false})
		} else {
			// ignore
		}

	}

	return dirs, files, nil
}

func (self *SFNetWalk) ListDir(pth string) (
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
					Scheme: "https",
					Host:   "sourceforge.net",
					Path:   path.Join("projects", self.project, "files", pth),
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

func (self *SFNetWalk) Walk(
	pth string,
	target func(
		dir string,
		dirs []os.FileInfo,
		files []os.FileInfo,
	) error,
) error {
	return self._Walk(pth, target, self.maxdepth)
}

func (self *SFNetWalk) _Walk(
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

func (self *SFNetWalk) Tree(pth string) (map[string]os.FileInfo, error) {

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

func (self *SFNetWalk) GetDownloadingURIForFile(name string) (string, error) {
	name = path.Base(name)

	tree, err := self.Tree("/")
	if err != nil {
		return "", err
	}

	for k, _ := range tree {
		if path.Base(k) == name {
			u := &url.URL{
				Scheme: "https",
				Host:   "sourceforge.net",
				Path:   path.Join("projects", self.project, "files", k, "download"),
			}
			return u.String(), nil
		}
	}

	return "", errors.New("not found")
}
