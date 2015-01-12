package util_test

import (
	"github.com/codingsince1985/bitrot_killer/util"
	"testing"
)

func TestEncodeAndWrite(t *testing.T) {
	files := []util.File{
		{"foo/bar", "asdf"},
		{"bar/foo", "fdsa"},
	}

	folder := util.Folder{
		Folder:    "/home/jerry/",
		Algorithm: "md5",
		Files:     files,
	}

	jsonBytes, err := folder.Encode()
	if err != nil {
		t.Error("TestEncodeAndWrite() failed", err)
	}

	expectedJsonDoc := `{
  "folder": "/home/jerry/",
  "algorithm": "md5",
  "files": [
    {
      "name": "foo/bar",
      "checksum": "asdf"
    },
    {
      "name": "bar/foo",
      "checksum": "fdsa"
    }
  ]
}`

	if string(jsonBytes) != expectedJsonDoc {
		t.Error("TestEncodeAndWrite() failed")
	}

	folder.Write("/home/jerry/temp.json")
}

func TestReadAndDecode(t *testing.T) {
	folder, err := util.Read("/home/jerry/temp.json")

	if err != nil {
		t.Error("TestReadAndDecode() failed", err)
	}

	if len(folder.Files) != 2 {
		t.Error("TestReadAndDecode() failed")
	}
}

func TestCompare(t *testing.T) {
	beforeFiles := []util.File{
		{"foo/bar", "asdf"},
		{"bar/foo", "fdsa"},
		{"bar/bar", "bar"},
		{"foobar", ""},
	}

	afterFiles := []util.File{
		{"foo/bar", "asdf"},
		{"bar/foo", "asdf"},
		{"foo/foo", "foo"},
		{"barfoo", ""},
	}

	changedFiles := util.ChangedFiles(beforeFiles, afterFiles)
	if len(changedFiles) != 1 || changedFiles[0].Name != "bar/foo" {
		t.Error("ChangedFiles() returned wrong changed files")
	}

	createdFiles, createdDirs := util.CreatedFiles(beforeFiles, afterFiles)
	if len(createdFiles) != 1 || createdFiles[0].Name != "foo/foo" || createdFiles[0].Checksum != "foo" {
		t.Error("CreatedFiles() returned wrong created files")
	}
	if len(createdDirs) != 1 || createdDirs[0].Name != "barfoo" || createdDirs[0].Checksum != "" {
		t.Error("CreatedFiles() returned wrong created dirs")
	}

	removedFiles, removedDirs := util.RemovedFiles(beforeFiles, afterFiles)
	if len(removedFiles) != 1 || removedFiles[0].Name != "bar/bar" || removedFiles[0].Checksum != "bar" {
		t.Error("RemovedFiles() returned wrong removed files")
	}
	if len(removedDirs) != 1 || removedDirs[0].Name != "foobar" || removedDirs[0].Checksum != "" {
		t.Error("RemovedFiles() returned wrong removed dirs")
	}
}
