package insertimports

import (
	"errors"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
)

func InsertImports(filename string, imports []string, verbose bool) error {

	if verbose {
		log.Println("reading", filename)
	}

	t, err := ioutil.ReadFile(filename)
	if err != nil {
		if verbose {
			log.Println("error", err)
		}
		return err
	}

	ts := string(t)

	tss := strings.Split(ts, "\n")

	imports_line := -1
	for ind, i := range tss {
		if strings.HasPrefix(i, "package") {
			if verbose {
				log.Println("package found. inserting..")
			}
			imports_line = ind
			break
		}
	}

	if imports_line == -1 {
		err = errors.New("can't find import's line")
		if verbose {
			log.Println("error", err)
		}
		return err
	}

	imports_line++

	tss = append(
		tss[0:imports_line],
		append(
			append(
				append(
					[]string{"//inserted", "import ("},
					imports...,
				),
				[]string{")"}...,
			),
			tss[imports_line:]...,
		)...,
	)

	ts = strings.Join(tss, "\n")

	// t = []byte(ts)

	if verbose {
		log.Println("formatting..")
	}

	t, err = format.Source([]byte(ts))
	if err != nil {
		if verbose {
			log.Println("error", err)
		}
		return err
	}

	if verbose {
		log.Println("saving..")
	}

	err = ioutil.WriteFile(filename, t, 0700)
	if err != nil {
		if verbose {
			log.Println("error", err)
		}
		return err
	}

	if verbose {
		log.Println("complete")
	}

	return nil
}
