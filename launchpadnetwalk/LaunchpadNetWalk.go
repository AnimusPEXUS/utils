package launchpadnetwalk

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
	"strings"

	"github.com/AnimusPEXUS/utils/cache01"
	"github.com/AnimusPEXUS/utils/filetools"
	"github.com/AnimusPEXUS/utils/logger"
)

var _ filetools.WalkerI = &LaunchpadNetWalk{}

type LPReleasesStruct struct {
	Entries []struct {
		Title              string `json:"title"`
		FileCollectionLink string `json:"files_collection_link"`
	} `json:"entries"`
}

type LPReleaseFilesStruct struct {
	Entries []struct {
		FileLink string `json:"file_link"`
	} `json:"entries"`
}

type LaunchpadNetWalk struct {
	project string
	cache   *cache01.CacheDir
	log     *logger.Logger
}

func NewLaunchpadNetWalk(
	project string,
	cache *cache01.CacheDir,
	log *logger.Logger,
) (*LaunchpadNetWalk, error) {
	self := new(LaunchpadNetWalk)
	self.project = project
	self.cache = cache
	self.log = log

	return self, nil
}

func (self *LaunchpadNetWalk) LogI(txt string) {
	if self.log != nil {
		self.log.Info(txt)
	}
}

func (self *LaunchpadNetWalk) LogE(txt string) {
	if self.log != nil {
		self.log.Error(txt)
	}
}

func (self *LaunchpadNetWalk) _GetReleases() (*LPReleasesStruct, error) {
	c, err := self.cache.Cache(
		"releases",
		func() ([]byte, error) {
			self.LogI("updating cache for releases")
			u := &url.URL{
				Scheme: "https",
				Host:   "api.launchpad.net",
				Path:   fmt.Sprintf("1.0/%s/releases", self.project),
			}

			http_res, err := http.Get(u.String())
			if err != nil {
				return []byte{}, err
			}

			b := new(bytes.Buffer)
			io.Copy(b, http_res.Body)

			return b.Bytes(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := c.GetValue()
	if err != nil {
		return nil, err
	}

	dst := new(LPReleasesStruct)

	err = json.Unmarshal(res, dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func (self *LaunchpadNetWalk) _GetReleaseFiles(title string) ([][2]string, error) {

	var uri string

	{
		rels, err := self._GetReleases()
		if err != nil {
			return nil, err
		}

		found := false
		for _, i := range rels.Entries {

			if i.Title == title {
				uri = i.FileCollectionLink
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New("title not found")
		}
	}

	c, err := self.cache.Cache(
		title,
		func() ([]byte, error) {

			self.LogI(fmt.Sprintf("updating cache for %s", title))
			self.LogI("  " + uri)

			http_res, err := http.Get(uri)
			if err != nil {
				return []byte{}, err
			}

			b := new(bytes.Buffer)
			io.Copy(b, http_res.Body)

			return b.Bytes(), nil
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := c.GetValue()
	if err != nil {
		return nil, err
	}

	dst := new(LPReleaseFilesStruct)

	err = json.Unmarshal(res, dst)
	if err != nil {
		return nil, err
	}

	ret := make([][2]string, 0)
	for _, i := range dst.Entries {
		l := i.FileLink
		pl, err := url.Parse(l)
		if err != nil {
			return nil, err
		}
		plps := strings.Split(pl.Path, "/")
		ret = append(ret, [2]string{plps[len(plps)-2], l})
	}

	return ret, nil
}

func (self *LaunchpadNetWalk) ListDir(pth string) (
	[]os.FileInfo,
	[]os.FileInfo,
	error,
) {
	// TODO: fix required: this is hack to only make it work. but logical problem isn't
	//       solved. "path science" policy shold be developed for all Walk functions
	//       of aipsetup and utils.
	pth = strings.Trim(pth, "/")

	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)

	if pth == "" {
		rels, err := self._GetReleases()
		if err != nil {
			return []os.FileInfo{}, []os.FileInfo{}, err
		}

		for _, i := range rels.Entries {
			dirs = append(dirs, &FileInfo{name: i.Title, isdir: true})
		}
	} else {
		ff, err := self._GetReleaseFiles(pth)
		if err != nil {
			return []os.FileInfo{}, []os.FileInfo{}, err
		}

		for _, i := range ff {
			files = append(files, &FileInfo{name: i[0], isdir: false})
		}

	}

	return dirs, files, nil

}

func (self *LaunchpadNetWalk) Walk(
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

func (self *LaunchpadNetWalk) Tree(pth string) (map[string]os.FileInfo, error) {

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

func (self *LaunchpadNetWalk) GetDownloadingURIForFile(name string) (string, error) {
	name = path.Base(name)

	releases := make([]string, 0)

	rels, err := self._GetReleases()
	if err != nil {
		return "", err
	}

	for _, i := range rels.Entries {
		releases = append(releases, i.Title)
	}

	for _, i := range releases {
		ff, err := self._GetReleaseFiles(i)
		if err != nil {
			return "", err
		}

		for _, j := range ff {
			if j[0] == name {
				return j[1], nil
			}
		}
	}

	return "", errors.New("not found")
}
