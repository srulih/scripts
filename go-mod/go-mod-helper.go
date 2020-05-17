package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"golang.org/x/mod/modfile"
)

func parseModFile(file string, bytes []byte) (*modfile.File, error) {
	f, err := modfile.Parse(file, bytes, nil)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func removeReplaces(mf *modfile.File) error {
	var err error
	for i := 0; i < len(mf.Replace); i++ {
		err = mf.DropReplace(mf.Replace[i].Old.Path, mf.Replace[i].Old.Version)
		if err != nil {
			return err
		}
	}
	return nil
}

func truncateAndWrite(f *os.File, b []byte) error {
	err := f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func appendCategory(a []string, b []string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}

func main() {
	replace := flag.Bool("replace", true, "whether to add a remove statement or remove a remove statement")
	root := flag.String("root", "../", "where to look for directories")
	flag.Parse()

	repos := flag.Args()
	modFile, err := os.OpenFile("go.mod", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	size, err := modFile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileBytes := make([]byte, size.Size())
	modFile.Read(fileBytes)
	f, err := parseModFile("go.mod", fileBytes)
	if err != nil {
		log.Fatal(err)
	}

	if *replace {
		for i := 0; i < len(f.Require); i++ {
			for j := 0; j < len(repos); j++ {
				requireRepo := f.Require[i].Syntax.Token[0]
				if strings.Contains(requireRepo, repos[j]) {
					err = f.AddReplace(requireRepo, "", *root+repos[j], "")
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	} else {
		err = removeReplaces(f)
		if err != nil {
			log.Fatal(err)
		}
	}
	modBytes, err := f.Format()
	err = truncateAndWrite(modFile, modBytes)
	if err != nil {
		log.Fatal(err)
	}
	modFile.Close()
}
