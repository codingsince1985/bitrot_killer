package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codingsince1985/bitrot_killer/util"
	"github.com/codingsince1985/checksum"
)

func main() {
	args := os.Args[1:]
	if len(args) >= 2 {
		var err error

		switch {
		case args[0] == "--create" && len(args) == 3:
			createChecksumFile(args[1], args[2])
		case args[0] == "--check":
			if len(args) == 4 {
				err = checkChecksumFile(args[1], args[2], args[3])
			} else if len(args) == 3 {
				err = checkChecksumFile(args[1], args[2], "")
			}
		case args[0] == "--dedup":
			err = checkDuplicated(args[1])
		}

		if err != nil {
			fmt.Println(err)
		}
	}
}

func checkDuplicated(checksumFile string) error {
	folder, err := util.Read(checksumFile)
	if err != nil {
		return err
	}

	fmt.Println("\nDuplicated files")
	for _, files := range groupByChecksum(folder.Files) {
		if len(files) > 1 {
			for i := range files {
				fmt.Println(files[i])
			}
			fmt.Println()
		}
	}

	fmt.Println("\nEmpty folders")
	for _, emptyFolder := range emptyFolders(folder.Files) {
		fmt.Println(emptyFolder.Name)
	}
	return nil
}

func groupByChecksum(files []util.File) map[string][]string {
	checksums := make(map[string][]string)
	for i := range files {
		if !files[i].IsFolder() {
			checksums[files[i].Checksum] = append(checksums[files[i].Checksum], files[i].Name)
		}
	}
	return checksums
}

func emptyFolders(files []util.File) []util.File {
	var folders []util.File
	for i := range files {
		if files[i].IsFolder() && hasNoFile(files[i].Name, files) {
			folders = append(folders, files[i])
		}
	}
	return folders
}

func hasNoFile(folderName string, files []util.File) bool {
	for i := range files {
		if !files[i].IsFolder() && strings.HasPrefix(files[i].Name, folderName) {
			return false
		}
	}
	return true
}

func createChecksumFile(root, checksumFile string) error {
	root = appendSlash(root)
	folder, err := getChecksum(root)
	if err != nil {
		return err
	}

	folder.Write(checksumFile)
	return nil
}

func getChecksum(root string) (util.Folder, error) {
	fileList, err := util.GetFiles(root)
	if err != nil {
		return util.Folder{}, err
	}
	return util.Folder{Folder: root, Algorithm: "md5", Files: getFiles(len(root), fileList[1:])}, nil
}

func getFiles(prefix int, files []string) []util.File {
	num := len(files)
	tasks := make(chan string, num)
	results := make(chan util.File, num)

	for i := 0; i < 2; i++ {
		go worker(prefix, tasks, results)
	}
	for i := range files {
		tasks <- files[i]
	}
	close(tasks)
	f := make([]util.File, num)
	for i := 0; i < num; i++ {
		f[i] = <-results
	}
	return f
}

func worker(prefix int, tasks <-chan string, results chan<- util.File) {
	for task := range tasks {
		md5sum, err := checksum.MD5sum(task)
		if err != nil {
			results <- util.File{Name: task[prefix:], Checksum: ""}
		}
		results <- util.File{Name: task[prefix:], Checksum: md5sum}
	}
}

func checkChecksumFile(root, checksumFile, remoteRoot string) error {
	root = appendSlash(root)
	folderAfter, err := getChecksum(root)
	if err != nil {
		return err
	}

	folderBefore, err := util.Read(checksumFile)
	if err != nil {
		return err
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
				return err
			}

			if strings.ToUpper(string(b[0])) == "Y" {
				remoteRoot = appendSlash(remoteRoot)
				applyChanged(root, remoteRoot, changedFiles)
				applyCreated(root, remoteRoot, createdFiles, createdDirs)
				applyRemoved(root, remoteRoot, removedFiles, removedDirs)

				folderAfter.Write(checksumFile)
			}
		}
	}
	return nil
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
	if len(dirs) > 0 {
		fmt.Println(title + " dirs")
		printFiles(dirs)
	}
	if len(files) > 0 {
		fmt.Println(title + " files")
		printFiles(files)
	}
}

func printFiles(files []util.File) {
	for _, file := range files {
		fmt.Println(file.Name)
	}
}
