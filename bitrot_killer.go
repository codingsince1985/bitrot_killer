package main

import (
	"fmt"
	"os"
	"sort"
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
		case args[0] == "--check" && len(args) == 3:
			err = checkChecksumFile(args[1], args[2])
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
			for _, file := range files {
				fmt.Println(file)
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
	for _, file := range files {
		if !file.IsFolder() {
			checksums[file.Checksum] = append(checksums[file.Checksum], file.Name)
		}
	}
	return checksums
}

func emptyFolders(files []util.File) []util.File {
	var folders []util.File
	for _, file := range files {
		if file.IsFolder() && hasNoFile(file.Name, files) {
			folders = append(folders, file)
		}
	}
	return folders
}

func hasNoFile(folderName string, files []util.File) bool {
	for _, file := range files {
		if !file.IsFolder() && strings.HasPrefix(file.Name, folderName) {
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
	for _, file := range files {
		tasks <- file
	}
	close(tasks)
	f := make([]util.File, num)
	for i := 0; i < num; i++ {
		f[i] = <-results
	}
	sort.Slice(f, func(i, j int) bool { return f[i].Name < f[j].Name })
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

func checkChecksumFile(root, checksumFile string) error {
	folderBefore, err := util.Read(checksumFile)
	if err != nil {
		return err
	}

	root = appendSlash(root)
	folderAfter, err := getChecksum(root)
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
		fmt.Print("\nUpdate? ")
		b := make([]byte, 1)
		if _, err := os.Stdin.Read(b); err != nil {
			return err
		}
		if strings.ToUpper(string(b[0])) == "Y" {
			folderAfter.Write(checksumFile)
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
