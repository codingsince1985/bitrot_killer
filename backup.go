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
	if len(args) >= 3 {
		switch {
		case args[1] == "--create" && len(args) == 4:
			createChecksumFile(args[2], args[3])
		case args[1] == "--check":
			if len(args) == 5 {
				checkChecksumFile(args[2], args[3], args[4])
			} else {
				checkChecksumFile(args[2], args[3], "")
			}
		case args[1] == "--dedup":
			checkDuplicated(args[2])
		}
	}
}

func checkDuplicated(checksumFile string) {
	folder, err := util.Read(checksumFile)
	if err != nil {
		panic(err)
	}

	dupFiles := [][]util.File{}
	for i, _ := range folder.Files {
		dupFiles = checkDuplicatedIn(folder.Files[i:], dupFiles)
	}

	fmt.Println("\nDuplicated files")
	for _, dupFile := range dupFiles {
		for _, file := range dupFile {
			fmt.Println(file.Name)
		}
		fmt.Println("")
	}
}

func checkDuplicatedIn(files []util.File, dupFiles [][]util.File) [][]util.File {
	if len(files) >= 2 && files[0].Checksum != "" {
		testedFile := files[0]

		for _, dupFile := range dupFiles {
			if testedFile.Checksum == dupFile[0].Checksum {
				return dupFiles
			}
		}

		dupFile := []util.File{testedFile}
		for _, file := range files[1:] {
			if testedFile.Checksum == file.Checksum {
				dupFile = append(dupFile, file)
			}
		}

		if len(dupFile) > 1 {
			dupFiles = append(dupFiles, dupFile)
		}
	}
	return dupFiles
}

func createChecksumFile(root, checksumFile string) {
	root = appendSlash(root)
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
	root = appendSlash(root)
	remoteRoot = appendSlash(remoteRoot)

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

func appendSlash(folder string) string {
	if !strings.HasSuffix(folder, "/") {
		return folder + "/"
	}
	return folder
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
