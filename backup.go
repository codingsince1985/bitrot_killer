package main

import (
	"fmt"
	"github.com/codingsince1985/checksum/md5"
	"github.com/codingsince1985/util"
	"os"
)

func main() {
	args := os.Args
	if len(args) >= 4 {
		switch args[1] {
		case "--create":
			createChecksumFile(args[2], args[3])
		case "--check":
			checkChecksumFile(args[2], args[3])
		}
	}
}

func createChecksumFile(root, file string) {
	folder, err := getChecksum(root)
	if err != nil {
		panic(err)
	}

	folder.Write(file)
}

func getChecksum(root string) (util.Folder, error) {
	fileList, err := util.GetFiles(root)
	if err != nil {
		panic(err)
	}

	files := []util.File{}
	for _, file := range fileList {
		md5sum, err := md5.MD5sum(file)
		if err != nil {
			panic(err)
		}
		files = append(files, util.File{file[len(root):], md5sum})
	}
	folder := util.Folder{root, "md5", files}
	return folder, nil
}

func checkChecksumFile(root, file string) {
	folderAfter, err := getChecksum(root)
	if err != nil {
		panic(err)
	}

	folderBefore, err := util.Read(file)
	if err != nil {
		panic(err)
	}

	changed, created, removed := util.Compare(folderBefore.Files, folderAfter.Files)

	if len(changed) > 0 {
		fmt.Println("Changed")
		printFiles(changed)
	}

	if len(created) > 0 {
		fmt.Println("Created")
		printFiles(created)
	}

	if len(removed) > 0 {
		fmt.Println("Removed")
		printFiles(removed)
	}
}

func printFiles(files []util.File) {
	for _, file := range files {
		fmt.Println(file.Name)
	}
}
