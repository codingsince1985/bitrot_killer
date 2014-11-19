package main

import (
	"github.com/codingsince1985/checksum/md5"
	"github.com/codingsince1985/util"
	"strings"
)

func main() {
	root := "/home/jerry/Downloads/Go/"
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
		files = append(files, util.File{file[strings.Index(file, root)+len(root):], md5sum})
	}
	folder := util.Folder{root, "md5", files}

	folder.Write("/home/jerry/Downloads/Go.json")
}
