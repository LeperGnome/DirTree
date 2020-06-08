package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

const (
	curVisual  string = "├───"
	lastVisual string = "└───"
	prevVisual string = "│\t"
)

func getFilePrefix(path string, isLast bool, depth int) string {
	var prefix string
	vis := curVisual
	if isLast {
		vis = lastVisual
	}
	prefix = strings.Repeat(prevVisual, depth-1) + vis
	return prefix

}

func getOnlyDirs(objects []os.FileInfo) (dirs []os.FileInfo) {
	for _, info := range objects {
		if info.IsDir() {
			dirs = append(dirs, info)
		}
	}
	return
}

func sizeOrEmpty(size int64) (sizeStr string) {
	if size == 0 {
		sizeStr = "(empty)"
	} else {
		sizeStr = fmt.Sprintf("(%vb)", size)
	}
	return
}

func describeDir(out io.Writer, path string, withFiles bool, curDepth int) error {
	curObject, err := os.Open(path)
	if err != nil {
		return err
	}
	dirContent, err := curObject.Readdir(0)
	if err != nil {
		return err
	}
	sort.Slice(dirContent, func(i, j int) bool { return dirContent[i].Name() < dirContent[j].Name() })
	if !withFiles {
		dirContent = getOnlyDirs(dirContent)
	}
	for i, info := range dirContent {
		isLast := i == len(dirContent)-1
		if info.IsDir() {
			fmt.Fprintf(out, "%v%v\n", getFilePrefix(path, isLast, curDepth), info.Name())
			err := describeDir(out, path+"/"+info.Name(), withFiles, curDepth+1)
			if err != nil {
				return err
			}
		} else if withFiles {
			fmt.Fprintf(out, "%v%v %v\n", getFilePrefix(path, isLast, curDepth), info.Name(), sizeOrEmpty(info.Size()))
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, withFiles bool) error {
	err := describeDir(out, path, withFiles, 1)
	if err != nil {
		return err
	}
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
