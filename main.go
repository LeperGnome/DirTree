package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	curVisual  string = "├───"
	prevVisual string = "│\t"
)

func getFilePrefix(path string) string {
	var prefix string
	n := strings.Count(path, "/")
	if n > 0 {
		prefix = strings.Repeat(prevVisual, n-1) + curVisual
	}
	return prefix

}

func step(path string, info os.FileInfo, err error) error {
	fname := info.Name()
	prefix := getFilePrefix(path)
	fmt.Println(prefix, fname)
	return nil
}

func dirTree(out io.Writer, path string, withFiles bool) error {
	var res string
	err := filepath.Walk(path, step)
	if err != nil {
		return err
	}
	fmt.Fprint(out, res)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
