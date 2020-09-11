package cache03

import (
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AnimusPEXUS/utils/worker"
)

// cache03 adds optional (by default disabled) UNIX Socket to cache02, so when
// new document pushed into cache - socket is triggered.
// also, the us_NextFile (and NextFile throught it) transparently caches list
// of files, to minimise calls to ioutil.ListDit() and speedup in case if
// many files is stored in the cache

// TODO: add trigger, using inotify

type EmptyStruct struct{}

type CacheDirOptions struct {
	DirPath       string
	WorkExtension string

	HashMaker     func() hash.Hash
	HashExtension string

	// UnixSocketServerSide bool
	UnixSocketEnabled  bool
	UnixSocketFileName string
}

type CacheDir struct {
	options *CacheDirOptions
	lenwe   int

	lockedFiles []string

	cached_working_files []os.FileInfo

	_RWMutex *sync.RWMutex

	// unix_conn_enabled     bool
	unix_conn_addr        *net.UnixAddr
	unix_conn             *net.UnixConn
	unix_conn_worker      *worker.Worker
	unix_conn_server_chan chan EmptyStruct
}

func NewCacheDir(options *CacheDirOptions) *CacheDir {

	self := &CacheDir{options: options}
	self.lenwe = len(options.WorkExtension)
	self.lockedFiles = make([]string, 0)

	self._RWMutex = &sync.RWMutex{}
	self.cached_working_files = make([]os.FileInfo, 0)

	// don't allow user to change this value during runtime
	// self.unix_conn_enabled = self.options.UnixSocketEnabled

	self.unix_conn_worker = worker.New(self.unixSocketSideWorkerThread)

	if self.options.UnixSocketEnabled {
		self.unix_conn_addr = &net.UnixAddr{
			Name: self.JoinFileName(self.options.UnixSocketFileName),
			Net:  "unixgram",
		}

		self.unix_conn_server_chan = make(chan EmptyStruct)
	}

	return self
}

const TXT_SERVER_SIDE_UNIX_SOCKET = "Server Side UNIX Socket"

func (self *CacheDir) unixSocketSideWorkerThread(
	set_starting func(),
	set_working func(),
	set_stopping func(),
	set_stopped func(),

	is_stop_flag func() bool,
) {
	set_starting()
	defer set_stopped()
	// if !self.options.UnixSocketServerSide {
	// 	log.Println(
	// 		TXT_SERVER_SIDE_UNIX_SOCKET,
	// 		"err:", "trying start server worker in client mode",
	// 	)
	// 	return
	// }

	os.Remove(self.unix_conn_addr.Name)
	conn, err := net.ListenUnixgram("unixgram", self.unix_conn_addr)
	if err != nil {
		log.Println(
			TXT_SERVER_SIDE_UNIX_SOCKET,
			"err:", "can't open listening unix socket",
			self.unix_conn_addr.String(),
		)
		return
	}
	defer func() {
		conn.Close()
		os.Remove(self.unix_conn_addr.Name)
	}()

	self.unix_conn = conn

	go func() {
		defer func() {
			conn.Close()
		}()
		for {
			if is_stop_flag() {
				break
			}

			time.Sleep(time.Second)
		}
	}()

	set_working()
	b := make([]byte, 1)
	for {
		if is_stop_flag() {
			break
		}
		_, err := conn.Read(b)
		if err != nil {
			log.Println(TXT_SERVER_SIDE_UNIX_SOCKET, "err:", err)
			break
		}

		if len(self.unix_conn_server_chan) < 1 {
			self.unix_conn_server_chan <- EmptyStruct{}
		}

	}
}

func (self *CacheDir) GetServerThreadWorker() worker.WorkerI {
	return self.unix_conn_worker
}

func (self *CacheDir) GetServerNewItemSignalChan() <-chan EmptyStruct {
	return self.unix_conn_server_chan
}

func (self *CacheDir) lockFile(name string) {

	if !self.isLocked(name) {
		self.lockedFiles = append(self.lockedFiles, name)
	}
	return
}

func (self *CacheDir) unlockFile(name string) {

	for i := len(self.lockedFiles) - 1; i != -1; i = i - 1 {
		if self.lockedFiles[i] == name {
			self.lockedFiles = append(self.lockedFiles[:i], self.lockedFiles[i+1:]...)
			break
		}
	}
	return
}

func (self *CacheDir) isLocked(name string) bool {

	for _, i := range self.lockedFiles {
		if i == name {
			return true
		}
	}
	return false
}

// Ensure directory acceptable to be used for file storage
func (self *CacheDir) EnsureDirectory(
	try_create_dir bool,
	try_write_file bool,
) error {
	dirstat, err := os.Stat(self.options.DirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if try_create_dir {
				err = os.MkdirAll(self.options.DirPath, 0o700)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
	} else {
		if !dirstat.IsDir() {
			err = errors.New(self.options.DirPath + " not a dirrectory")
			return err
		}
	}

	if try_write_file {
		tmpfile := self.JoinFileName("dirwritetestfilename" + time.Now().UTC().Format(time.RFC3339Nano))
		f, err := os.Create(tmpfile)
		if err != nil {
			return err
		}
		f.Close()
		err = os.Remove(tmpfile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *CacheDir) HaveCache() (bool, error) {
	res, err := self.WorkingFiles()
	if err != nil {
		return false, err
	}
	return len(res) != 0, nil
}

// preforms filepath joining of directory path and provided filename
func (self *CacheDir) JoinFileName(filename string) string {
	filename = filepath.Base(filename)
	return filepath.Join(self.options.DirPath, filename)
}

func (self *CacheDir) WorkingFiles() (ret []os.FileInfo, err error) {
	self._RWMutex.RLock()
	defer self._RWMutex.RUnlock()
	return self.us_WorkingFiles()
}

// return file info for files, which names matching working cache names criteria
// no checksum checking performed except echecksum file existance
// doesn't return err, if no files found
func (self *CacheDir) us_WorkingFiles() (ret []os.FileInfo, err error) {

	ret = make([]os.FileInfo, 0)

	files, err := ioutil.ReadDir(self.options.DirPath)
	if err != nil {
		return
	}

	if len(files) == 0 {
		return nil, nil
	}

	for i := len(files) - 1; i != -1; i = i - 1 {

		name := files[i].Name()

		_, err = self.ParseWorkFileName(name)
		if err != nil {
			goto remove
		}

		if self.isLocked(name) {
			goto remove
		}

		{
			_, _, name_sum, _ := self.GenNames(name)
			_, err = os.Stat(name_sum)
			if err != nil {
				goto remove
			}
		}

		continue
	remove:
		files = append(files[:i], files[i+1:]...)
	}

	err = nil
	ret = files

	return
}

// err = os.ErrNotExist if no any file found.
// acceptable filenames should be already checked by self.WorkingFiles()
// so NextFile() only find's the oldest one and treats any errors as not acceptable
func (self *CacheDir) NextFile() (name string, err error) {
	return self.us_NextFile()
}

func (self *CacheDir) us_NextFile() (name string, err error) {

	if len(self.cached_working_files) == 0 {
		var files []os.FileInfo
		files, err = self.us_WorkingFiles()
		if err != nil {
			return
		}

		len_files := len(files)

		if len_files < 2 {
			self.cached_working_files = files
			goto ifexit
		}

		for i := 0; i != len_files-1; i++ {
			for j := i + 1; j != len_files; j++ {
				var comp_res int
				comp_res, err = self.ComparisonFunction(files[i], files[j])
				if err != nil {
					return
				}

				if comp_res < 0 {
					z := files[i]
					files[i] = files[j]
					files[j] = z
				}

			}
		}

		self.cached_working_files = files
	}
ifexit:

	if len(self.cached_working_files) == 0 {
		err = os.ErrNotExist
		return
	}

	oldest := self.cached_working_files[0]
	self.cached_working_files = self.cached_working_files[1:]

	name = oldest.Name()

	return
}

func (self *CacheDir) ParseWorkFileName(n string) (t time.Time, err error) {
	if !strings.HasSuffix(n, self.options.WorkExtension) {
		err = errors.New("file name not acceptable")
		return
	}
	return time.Parse(time.RFC3339Nano, n[:len(n)-self.lenwe])
}

func (self *CacheDir) ParseWorkFileNameByFileInfo(fi os.FileInfo) (t time.Time, err error) {
	return self.ParseWorkFileName(fi.Name())
}

// NOTE: synonym not needed imo
// func (self *CacheDir) TimeFromName(name string) (t time.Time, err error) {
// 	return self.ParseWorkFileName(name)
// }

func (self *CacheDir) ComparisonFunction(f1, f2 os.FileInfo) (int, error) {

	fn1, err := self.ParseWorkFileNameByFileInfo(f1)
	if err != nil {
		return 0, err
	}

	fn2, err := self.ParseWorkFileNameByFileInfo(f2)
	if err != nil {
		return 0, err
	}

	if fn1.Equal(fn2) {
		return 0, nil
	}

	if fn1.Before(fn2) {
		return -1, nil
	}

	return 1, nil
}

// to many functionality. this detalisation if overhead, imo. use GenNames()
// func (self *CacheDir) GenChecksumName(name string) string {
// 	name = filepath.Base(name)
// 	return name + self.options.HashExtension
// }

// TODO: probaby, this function requires optimization
// oname must be suffixed with WorkExtension. it will be suffized automatically, if it's not
func (self *CacheDir) GenNames(oname string) (name, name_disabled, name_sum, name_sum_disabled string) {
	if !strings.HasSuffix(oname, self.options.WorkExtension) {
		oname = oname + self.options.WorkExtension
	}

	name = oname
	name = filepath.Base(name)
	name_disabled = name + ".disabled"

	name_sum = name + self.options.HashExtension
	name_sum_disabled = name_sum + ".disabled"

	name = self.JoinFileName(name)
	name_disabled = self.JoinFileName(name_disabled)
	name_sum = self.JoinFileName(name_sum)
	name_sum_disabled = self.JoinFileName(name_sum_disabled)
	return
}

func (self *CacheDir) UnlockFile(name string) {
	self._RWMutex.Lock()
	defer self._RWMutex.Unlock()

	self.us_UnlockFile(name)

	return
}

func (self *CacheDir) us_UnlockFile(name string) {

	self.unlockFile(name)

	return
}

func (self *CacheDir) Disable(name string) {
	self._RWMutex.Lock()
	defer self._RWMutex.Unlock()

	self.us_Disable(name)

	return
}

func (self *CacheDir) us_Disable(name string) {

	nname, name_disabled, name_sum, name_sum_disabled := self.GenNames(name)

	os.Rename(nname, name_disabled)
	os.Rename(name_sum, name_sum_disabled)

	self.unlockFile(name)

	return
}

func (self *CacheDir) Delete(name string) {
	self._RWMutex.Lock()
	defer self._RWMutex.Unlock()

	self.us_Delete(name)

	return
}

func (self *CacheDir) us_Delete(name string) {

	nname, name_disabled, name_sum, name_sum_disabled := self.GenNames(name)

	for _, i := range []string{nname, name_disabled, name_sum, name_sum_disabled} {
		os.Remove(i)
	}

	self.unlockFile(name)

	return
}

func (self *CacheDir) Put(data io.Reader) (err error) {

	name := time.Now().UTC().Format(time.RFC3339Nano) + self.options.WorkExtension

	name, _, name_sum, _ := self.GenNames(name)

	_, err = os.Stat(name)

	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		err = os.ErrExist
		return
	}

	f, err := os.Create(name)
	if err != nil {
		return
	}

	h := self.options.HashMaker()

	data_tee := io.TeeReader(data, h)

	_, err = io.Copy(f, data_tee)
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}

	sum := h.Sum([]byte{})
	f, err = os.Create(name_sum)
	if err != nil {
		return
	}

	_, err = f.WriteString(hex.EncodeToString(sum))
	if err != nil {
		return
	}

	err = f.Close()
	if err != nil {
		return
	}

	if self.options.UnixSocketEnabled {
		go func() {
			conn, err := net.DialUnix("unixgram", nil, self.unix_conn_addr)
			if err != nil {
				log.Println("err:", err)
				return
			}
			_, err = conn.Write([]byte{1})
			if err != nil {
				log.Println("err:", err)
				return
			}
		}()
	}

	return
}

func (self *CacheDir) CheckFileIntegrity(name string) (ok bool, fullpath string, err error) {
	_, err = self.ParseWorkFileName(name)
	if err != nil {
		return
	}

	name, _, name_sum, _ := self.GenNames(name)

	fullpath = name

	var saved_sum string
	var fresh_sum string
	{
		var data []byte
		data, err = ioutil.ReadFile(name_sum)
		if err != nil {
			return
		}
		saved_sum = strings.Trim(string(data), "\n\r\t\x00 ")
	}

	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	{
		h := self.options.HashMaker()

		_, err = io.Copy(h, f)
		if err != nil {
			return
		}

		fresh_sum = hex.EncodeToString(h.Sum([]byte{}))
	}

	ok = saved_sum == fresh_sum

	return
}

// Get() locks file with returned name, so Get() can be called asyncronously
// and Get() will not return locked files.
// files can be unlocked with Unlock(), Disable() or Delete() functions
func (self *CacheDir) Get() (name string, data io.ReadCloser, err error) {
	self._RWMutex.Lock()
	defer self._RWMutex.Unlock()

start:
	name, err = self.us_NextFile()
	if err != nil {
		return
	}

	ok, fullpath, err := self.CheckFileIntegrity(name)
	if err != nil {
		goto disable_and_restart
	}

	if !ok {
		goto disable_and_restart
	}

	{
		var f io.ReadCloser
		f, err = os.Open(fullpath)
		if err != nil {
			goto disable_and_restart
		}
		name = filepath.Base(name)
		self.lockFile(name)
		data = f
		return
	}

disable_and_restart:
	self.us_Disable(name)
	goto start

	return
}
