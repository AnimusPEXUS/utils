package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	filename := os.Args[1]

	imports := os.Args[2:]

	for i := range imports {
		imports[i] = fmt.Sprintf("\"%s\"", imports[i])
	}

	log.Println("going to open file", filename, "and insert following imports into it:")
	for _, i := range imports {
		log.Println("  ", i)
	}

	t, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln("error", err)
	}

	ts := string(t)

	tss := strings.Split(ts, "\n")

	imports_line := -1
	for ind, i := range tss {
		if strings.HasPrefix(i, "package") {
			imports_line = ind
			log.Println("found line", i, "which index is", ind)
			log.Println("  inserting imports")
			break
		}
	}

	if imports_line == -1 {
		log.Fatalln("error", "can't find import's line")
	}

	imports_line++

	tss = append(
		tss[0:imports_line],
		append(
			append(
				append(
					[]string{"import ("},
					imports...,
				),
				[]string{")"}...,
			),
			tss[imports_line:]...,
		)...,
	)

	ts = strings.Join(tss, "\n")

	log.Println("formatting...")

	t, err = format.Source([]byte(ts))
	if err != nil {
		log.Fatalln("error", err)
	}

	err = ioutil.WriteFile(filename, t, 0700)
	if err != nil {
		log.Fatalln("error", err)
	}

	log.Println("completed")

}
