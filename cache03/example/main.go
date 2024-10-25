package main

import (
	"bytes"
	"crypto/sha512"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/AnimusPEXUS/utils/cache02"
)

func main() {
	log.SetFlags(log.Llongfile)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	cdo := &cache02.CacheDirOptions{
		DirPath:       filepath.Join(cwd, "testcache"),
		WorkExtension: ".cached",
		HashMaker:     sha512.New,
		HashExtension: ".sha512",
	}

	cd := cache02.NewCacheDir(cdo)
	err = cd.EnsureDirectory(true, true)
	if err != nil {
		log.Fatalln(err)
	}

	err = cd.Put(bytes.NewBufferString("somedata"))
	if err != nil {
		log.Fatalln(err)
	}

	// name, err := cd.NextFile()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	for {

		have_cache, err := cd.HaveCache()
		if err != nil {
			log.Fatalln(err)
		}
		if !have_cache {
			break
		}

		name, data, err := cd.Get()
		if err != nil {
			log.Fatalln(err)
		}

		data_b := bytes.NewBuffer([]byte{})

		_, err = io.Copy(data_b, data)
		data.Close()

		log.Printf("name: %s, data: %s", name, string(data_b.Bytes()))
		cd.Delete(name)
	}

}
