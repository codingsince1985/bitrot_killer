package main

import (
	"github.com/codingsince1985/checksum/md5"
	"github.com/codingsince1985/util"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 4 {
		switch args[1] {
		case "--create":
			createChecksumFile(args[2], args[3])
		}
	}
}

func createChecksumFile(root, file string) {
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

	folder.Write(file)
}
