package main

import (
	"fmt"
	"github.com/codingsince1985/checksum/md5"
	"github.com/codingsince1985/util"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) >= 4 {
		root := args[2]
		if !strings.HasSuffix(args[2], "/") {
			root = args[2] + "/"
		}

		switch args[1] {
		case "--create":
			createChecksumFile(root, args[3])
		case "--check":
			if len(args) == 4 {
				checkChecksumFile(root, args[3], "")
			} else {
				remoteRoot := args[4]
				if !strings.HasSuffix(args[4], "/") {
					remoteRoot = args[4] + "/"
				}
				checkChecksumFile(root, args[3], remoteRoot)
			}
		}
	}
}

func createChecksumFile(root, checksumFile string) {
	folder, err := getChecksum(root)
	if err != nil {
		panic(err)
	}

	folder.Write(checksumFile)
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

func checkChecksumFile(root, checksumFile, remoteRoot string) {
	folderAfter, err := getChecksum(root)
	if err != nil {
		panic(err)
	}

	folderBefore, err := util.Read(checksumFile)
	if err != nil {
		panic(err)
	}

	changed, created, removed := util.Compare(folderBefore.Files, folderAfter.Files)
	printComparison(changed, created, removed)

	if len(changed) > 0 || len(created) > 0 || len(removed) > 0 {
		if remoteRoot != "" {
			fmt.Print("\nSync changes to remote folder? ")
			b := make([]byte, 2)
			if _, err := os.Stdin.Read(b); err != nil {
				panic(err)
			}

			if strings.ToUpper(string(b[0])) == "Y" {
				performUpdate(root, remoteRoot, changed, created, removed)
			}
		}

		folderAfter.Write(checksumFile)
	}
}

func performUpdate(root, remoteRoot string, changed, created, removed []util.File) {
	if len(changed) > 0 {
		fmt.Println("\nChanged")
		for _, file := range changed {
			fmt.Println(file.Name)
			util.CopyRemoteFile(root+file.Name, remoteRoot+file.Name)
		}
	}

	if len(created) > 0 {
		fmt.Println("\nCreated")
		for _, file := range created {
			fmt.Println(file.Name)
			util.CopyRemoteFile(root+file.Name, remoteRoot+file.Name)
		}
	}

	if len(removed) > 0 {
		fmt.Println("\nRemoved")
		for _, file := range removed {
			fmt.Println(file.Name)
			util.DeleteRemoteFile(remoteRoot + file.Name)
		}
	}
}

func printComparison(changed, created, removed []util.File) {
	if len(changed) > 0 {
		fmt.Println("\nChanged")
		printFiles(changed)
	}

	if len(created) > 0 {
		fmt.Println("\nCreated")
		printFiles(created)
	}

	if len(removed) > 0 {
		fmt.Println("\nRemoved")
		printFiles(removed)
	}
}

func printFiles(files []util.File) {
	for _, file := range files {
		fmt.Println(file.Name)
	}
}
