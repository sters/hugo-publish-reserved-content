package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
	"github.com/sters/hugo-publish-reserved-content/publish"
)

func abs(p string) string {
	f, _ := filepath.Abs(p)
	return f
}

func dirwalk(dir string) ([]string, error) {
	dir = abs(dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			child, err := dirwalk(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, failure.Wrap(err)
			}
			paths = append(paths, child...)
			continue
		}

		p, err := filepath.Abs(filepath.Join(dir, file.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%+v", err)
			continue
		}
		paths = append(paths, p)
	}

	return paths, nil
}

func main() {
	var (
		reservedKey string
		draftKey    string
		basePath    string
	)
	flag.StringVar(&reservedKey, "reservedKey", "reserved", "hugo content's reservation key (bool), default = reserved")
	flag.StringVar(&draftKey, "draftKey", "draft", "hugo content's draft key (bool), default = draft")
	flag.StringVar(&basePath, "basePath", "", "hugo content's root directory")

	if reservedKey == "" {
		fmt.Fprintf(os.Stderr, "reservedKey is required.")
	}
	if draftKey == "" {
		fmt.Fprintf(os.Stderr, "draftKey is required.")
	}
	if basePath == "" {
		fmt.Fprintf(os.Stderr, "basePath is required.")
	}

	dirs, err := dirwalk(basePath)
	if err != nil {
		log.Fatal(err)
	}

	p := publish.New("", "")
	for _, filepath := range dirs {
		c, ok := failure.CodeOf(p.CheckReservedAndPublish(filepath))
		if !ok {
			// = no error
			fmt.Fprintf(os.Stdout, "%s is published.", filepath)
			continue
		}

		switch c {
		case publish.ErrContentIsReservedButNotDraft:
			fmt.Fprintf(os.Stderr, "%s is reserved but not draft.", filepath)
		case publish.ErrFileContentMismatch:
			fmt.Fprintf(os.Stderr, "%s is maybe breaking content.", filepath)
		case publish.ErrContentIsNotTheTimeYet:
			fmt.Fprintf(os.Stderr, "%s is still waiting.", filepath)
		}
	}
}
