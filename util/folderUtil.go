package util

import (
	"encoding/json"
	"io/ioutil"
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
		return ioutil.WriteFile(filepath, bytes, 0644)
	}
}

func Read(filepath string) (Folder, error) {
	bytes, err := ioutil.ReadFile(filepath)
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

func CreatedFiles(beforeFiles, afterFiles []File) (createdFiles, createdDirs []File) {
	for _, f1 := range afterFiles {
		found := false
		for _, f2 := range beforeFiles {
			if f1.Name == f2.Name {
				found = true
				break
			}
		}
		if !found {
			if f1.Checksum == "" {
				createdDirs = append(createdDirs, f1)
			} else {
				createdFiles = append(createdFiles, f1)
			}
		}
	}
	return
}

func RemovedFiles(beforeFiles, afterFiles []File) (removedFiles, removedDirs []File) {
	for _, f1 := range beforeFiles {
		found := false
		for _, f2 := range afterFiles {
			if f1.Name == f2.Name {
				found = true
				break
			}
		}
		if !found {
			if f1.Checksum == "" {
				removedDirs = append(removedDirs, f1)
			} else {
				removedFiles = append(removedFiles, f1)
			}
		}
	}
	return
}
