package gnupghtmlftpwalk

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
	"sort"
	"strconv"
	"strings"

	"github.com/AnimusPEXUS/utils/cache01"
	"github.com/AnimusPEXUS/utils/filetools"
	"github.com/AnimusPEXUS/utils/logger"
	"github.com/antchfx/htmlquery"
)

var _ filetools.WalkerI = &Walk{}

type Walk struct {
	scheme string
	host   string
	cache  *cache01.CacheDir
	log    *logger.Logger
}

func NewWalk(
	scheme string,
	host string,
	cache *cache01.CacheDir,
	log *logger.Logger,
) (*Walk, error) {
	self := new(Walk)
	self.scheme = scheme
	self.host = host

	self.cache = cache
	self.log = log

	return self, nil
}

func (self *Walk) LogI(txt string) {
	if self.log != nil {
		self.log.Info(txt)
	}
}

func (self *Walk) LogE(txt string) {
	if self.log != nil {
		self.log.Error(txt)
	}
}

func (self *Walk) ListDirNotCached(pth string) (
	[]os.FileInfo,
	[]os.FileInfo,
	error,
) {
	u := &url.URL{
		Scheme: self.scheme,
		Host:   self.host,
		Path:   pth,
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, nil, err
	}

	if !strings.HasPrefix(strconv.Itoa(resp.StatusCode), "2") {
		return []os.FileInfo{}, []os.FileInfo{},
			fmt.Errorf("http code %d", resp.StatusCode)
	}

	b := new(bytes.Buffer)
	io.Copy(b, resp.Body)
	resp.Body.Close()

	doc, err := htmlquery.Parse(b)

	file_list_table_res := htmlquery.Find(doc, `.//table[@class="ftp"]`)
	// if len(file_list_table_res) != 1 {
	// 	return []os.FileInfo{}, []os.FileInfo{},
	// 		errors.New("found invalid number of ftp tables")
	// }

	dirs := make([]os.FileInfo, 0)
	files := make([]os.FileInfo, 0)

	for _, file_list_table := range file_list_table_res {

		file_list_table_tbody_res := htmlquery.Find(file_list_table, "tbody")
		if len(file_list_table_tbody_res) != 1 {
			return nil, nil, errors.New("found invalid number of tbody")
		}

		file_list_table_tbody := file_list_table_tbody_res[0]

		folder_trs := htmlquery.Find(file_list_table_tbody, "tr")

		for _, i := range folder_trs {

			row_tds := htmlquery.Find(i, "td")

			if len(row_tds) < 2 {
				continue
			}

			row_td_0 := row_tds[0]
			row_td_1 := row_tds[1]

			row_td_0_c := htmlquery.FindOne(row_td_0, "img")
			row_td_1_c := htmlquery.FindOne(row_td_1, "a")

			if row_td_0_c == nil || row_td_1_c == nil {
				continue
			}

			var (
				row_td_0_c_src_str            string
				row_td_1_c_href_str           string
				row_td_1_c_href_str_unescaped string
			)

			for _, v := range row_td_0_c.Attr {
				if v.Key == "src" {
					row_td_0_c_src_str = v.Val
					break
				}
			}

			if row_td_0_c_src_str == "" {
				continue
			}

			for _, v := range row_td_1_c.Attr {
				if v.Key == "href" {
					row_td_1_c_href_str = v.Val
					break
				}
			}

			if row_td_1_c_href_str == "" {
				continue
			}

			row_td_1_c_href_str_unescaped, err = url.PathUnescape(row_td_1_c_href_str)
			if err != nil {
				return nil, nil, err
			}

			switch row_td_0_c_src_str {
			case "/share/folder.png":
				dirs = append(
					dirs,
					&FileInfo{
						name:  row_td_1_c_href_str_unescaped,
						isdir: true,
					},
				)
			default:
				files = append(
					files,
					&FileInfo{
						name:  row_td_1_c_href_str_unescaped,
						isdir: false,
					},
				)
			}

		}
	}

	return dirs, files, nil
}

func (self *Walk) ListDir(pth string) (
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

func (self *Walk) Walk(
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

func (self *Walk) Tree(pth string) (map[string]os.FileInfo, error) {

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

func (self *Walk) GetDownloadingURIForFile(
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
