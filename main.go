package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	SPACE_OFFSET  = "	"
	STICK_OFFSET  = "│	"
	DIR_DELIMITER = "/"
	CUR_DIR       = "." + DIR_DELIMITER
)

type DirInfo struct {
	fileCount     int
	endfOfReading bool
}

func directoryCount(root string) ([]os.FileInfo, error) {
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var dirs []os.FileInfo
	for i := 0; i < len(fileInfo); i++ {
		if fileInfo[i].IsDir() {
			dirs = append(dirs, fileInfo[i])
		}
	}
	return dirs, nil
}

func directoryPath(path, filename string) string {
	dirPath := strings.Replace(path, filename, "", 1)
	if dirPath == "" {
		dirPath = CUR_DIR
	} else {
		dirPath = CUR_DIR + dirPath
	}
	return dirPath
}

func pathDirectories(dirPath string) []string {
	return strings.Split(dirPath, DIR_DELIMITER)
}

func directoryFileCount(dirPath string, printFilesFlag bool) (int, error) {
	dirInfo, err := ioutil.ReadDir(dirPath)
	if !printFilesFlag {
		dirInfo, err = directoryCount(dirPath)
	}
	if err != nil {
		return 0, err
	}
	return len(dirInfo), nil
}

func printOffset(treeMap map[string]DirInfo, pathDirs []string, rootPath string, out io.Writer) {
	rootDir := ""
	i := 0
	for pathDirs[i] != rootPath {
		rootDir += pathDirs[i] + DIR_DELIMITER
		i++
	}
	for ; i < len(pathDirs)-2; i++ {
		rootDir += pathDirs[i] + DIR_DELIMITER

		if !treeMap[rootDir].endfOfReading {
			fmt.Fprintf(out, STICK_OFFSET)
		} else {
			fmt.Fprint(out, SPACE_OFFSET)
		}
	}
}

func printTreeNodes(treeMap map[string]DirInfo, dirPath string, info os.FileInfo, out io.Writer) {
	if !treeMap[dirPath].endfOfReading {
		if info.IsDir() {
			fmt.Fprintf(out, fmt.Sprintf("├───%s\n", info.Name()))
		} else {
			if info.Size() == 0 {
				fmt.Fprintf(out, fmt.Sprintf("├───%s (empty)\n", info.Name()))
			} else {
				fmt.Fprintf(out, fmt.Sprintf("├───%s (%db)\n", info.Name(), info.Size()))
			}
		}
	} else {
		if info.IsDir() {
			fmt.Fprintf(out, fmt.Sprintf("└───%s\n", info.Name()))
		} else {
			if info.Size() == 0 {
				fmt.Fprintf(out, fmt.Sprintf("└───%s (empty)\n", info.Name()))
			} else {
				fmt.Fprintf(out, fmt.Sprintf("└───%s (%db)\n", info.Name(), info.Size()))
			}
		}
	}
}

func dirTree(out io.Writer, rootPath string, printFiles bool) error {

	treeMap := make(map[string]DirInfo)

	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if path == rootPath || (!printFiles && !info.IsDir()) {
			return nil
		}

		dirPath := directoryPath(path, info.Name())
		pathDirs := pathDirectories(dirPath)
		dirSize, _ := directoryFileCount(dirPath, printFiles)

		if _, found := treeMap[dirPath]; !found {
			treeMap[dirPath] = DirInfo{1, false}
		} else {
			treeMap[dirPath] = DirInfo{treeMap[dirPath].fileCount + 1, false}
		}

		if treeMap[dirPath].fileCount == dirSize {
			treeMap[dirPath] = DirInfo{treeMap[dirPath].fileCount, true}
		}

		if info.Name() == ".DS_Store" {
			return nil
		}

		printOffset(treeMap, pathDirs, rootPath, out)
		printTreeNodes(treeMap, dirPath, info, out)

		return nil
	})

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
