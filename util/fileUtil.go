package util

import (
	"os"
	"os/exec"
	"path/filepath"
)

var fileList = []string{}

func visit(path string, info os.FileInfo, err error) error {
	fileList = append(fileList, path)
	return nil
}

func GetFiles(path string) ([]string, error) {
	err := filepath.Walk(path, visit)
	if err != nil {
		return nil, err
	}

	return fileList, nil
}

func CopyRemoteFile(sourceFile, destFile string) error {
	cmd := exec.Command("gvfs-copy", sourceFile, destFile)
	err := cmd.Run()
	return err
}

func MakeRemoteDir(remoteDir string) error {
	cmd := exec.Command("gvfs-mkdir", "-p", remoteDir)
	err := cmd.Run()
	return err
}

func DeleteRemoteFile(filepath string) error {
	cmd := exec.Command("gvfs-rm", filepath)
	err := cmd.Run()
	return err
}
