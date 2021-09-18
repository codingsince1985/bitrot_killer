package util

import (
	"encoding/json"
	"os"
)

type (
	Folder struct {
		Folder    string `json:"folder"`
		Algorithm string `json:"algorithm"`
		Files     []File `json:"files"`
	}
	File struct {
		Name     string `json:"name"`
		Checksum string `json:"checksum"`
	}
	ChecksumBytes []byte
)

func (file File) IsFolder() bool {
	return file.Checksum == ""
}

func (folder Folder) Encode() ([]byte, error) {
	if bytes, err := json.MarshalIndent(folder, "", "  "); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func (jsonBytes ChecksumBytes) decode() (Folder, error) {
	var folder Folder

	if err := json.Unmarshal(jsonBytes, &folder); err != nil {
		return Folder{}, err
	}
	return folder, nil
}

func (folder Folder) Write(filepath string) error {
	if bytes, err := folder.Encode(); err != nil {
		return err
	} else {
		return os.WriteFile(filepath, bytes, 0644)
	}
}

func Read(filepath string) (Folder, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return Folder{}, err
	}

	folder, err := ChecksumBytes(bytes).decode()
	if err != nil {
		return Folder{}, err
	}
	return folder, nil
}

func ChangedFiles(beforeFiles, afterFiles []File) (changed []File) {
	for _, f1 := range beforeFiles {
		for _, f2 := range afterFiles {
			if f1.Name == f2.Name && f1.Checksum != f2.Checksum {
				changed = append(changed, f1)
				break
			}
		}
	}
	return
}

func CreatedFiles(prev, curr []File) (diffFiles, diffDirs []File) {
	return findDiff(curr, prev)
}

func RemovedFiles(prev, curr []File) (diffFiles, diffDirs []File) {
	return findDiff(prev, curr)
}

func findDiff(files1, files2 []File) (diffFiles, diffDirs []File) {
	for _, f1 := range files1 {
		found := false
		for _, f2 := range files2 {
			if f1.Name == f2.Name {
				found = true
				break
			}
		}
		if !found {
			if f1.IsFolder() {
				diffDirs = append(diffDirs, f1)
			} else {
				diffFiles = append(diffFiles, f1)
			}
		}
	}
	return
}
