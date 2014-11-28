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
		if file != root {
			md5sum, err := md5.MD5sum(file)
			if err != nil {
				panic(err)
			}
			files = append(files, util.File{file[len(root):], md5sum})
		}
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

	changedFiles := util.ChangedFiles(folderBefore.Files, folderAfter.Files)
	print("\nChanged", changedFiles, nil)

	createdFiles, createdDirs := util.CreatedFiles(folderBefore.Files, folderAfter.Files)
	print("\nCreated", createdFiles, createdDirs)

	removedFiles, removedDirs := util.RemovedFiles(folderBefore.Files, folderAfter.Files)
	print("\nRemoved", removedFiles, removedDirs)

	if len(changedFiles) > 0 || len(createdFiles) > 0 || len(createdDirs) > 0 || len(removedFiles) > 0 || len(removedDirs) > 0 {
		if remoteRoot != "" {
			fmt.Print("\nSync changes to remote folder? ")
			b := make([]byte, 1)
			if _, err := os.Stdin.Read(b); err != nil {
				panic(err)
			}

			if strings.ToUpper(string(b[0])) == "Y" {
				applyChanged(root, remoteRoot, changedFiles)
				applyCreated(root, remoteRoot, createdFiles, createdDirs)
				applyRemoved(root, remoteRoot, removedFiles, removedDirs)
			}
		}

		folderAfter.Write(checksumFile)
	}
}

func applyChanged(root, remoteRoot string, changedFiles []util.File) {
	if len(changedFiles) > 0 {
		fmt.Println("\nChanged")

		for _, file := range changedFiles {
			fmt.Println(file.Name)
			util.CopyRemoteFile(root+file.Name, remoteRoot+file.Name)
		}
	}
}

func applyCreated(root, remoteRoot string, createdFiles, createdDirs []util.File) {
	if len(createdFiles) > 0 {
		fmt.Println("\nCreated")

		for _, file := range createdDirs {
			fmt.Println(file.Name)
			util.MakeRemoteDir(remoteRoot + file.Name)
		}

		for _, file := range createdFiles {
			fmt.Println(file.Name)
			util.CopyRemoteFile(root+file.Name, remoteRoot+file.Name)
		}
	}
}

func applyRemoved(root, remoteRoot string, removedFiles, removedDirs []util.File) {
	if len(removedFiles) > 0 {
		fmt.Println("\nRemoved")

		for _, file := range removedFiles {
			fmt.Println(file.Name)
			util.DeleteRemoteFile(remoteRoot + file.Name)
		}

		for _, file := range removedDirs {
			fmt.Println(file.Name)
			util.DeleteRemoteFile(remoteRoot + file.Name)
		}
	}
}

func print(title string, files, dirs []util.File) {
	if len(files) > 0 {
		fmt.Println(title + " files")
		printFiles(files)
	}
	if len(dirs) > 0 {
		fmt.Println(title + " dirs")
		printFiles(dirs)
	}
}

func printFiles(files []util.File) {
	for _, file := range files {
		fmt.Println(file.Name)
	}
}
