package util_test

import (
	"github.com/codingsince1985/bitrot_killer/util"
	"testing"
)

func TestGetFiles(t *testing.T) {
	path := "/home/jerry/Downloads"
	if _, err := util.GetFiles(path); err != nil {
		t.Error("GetFiles() failed", err)
	}
}

func TestMakeRemoteDir(t *testing.T) {
	t.SkipNow()
	if err := util.MakeRemoteDir("smb://192.168.8.140/public/sda1/Jerry/temp"); err != nil {
		t.Error("TestMakeRemoteDir() failed", err)
	}
}

func TestCopyRemoteFile(t *testing.T) {
	t.SkipNow()
	if err := util.CopyRemoteFile("/home/jerry/temp.json", "smb://192.168.8.140/public/sda1/Jerry/temp/temp.json"); err != nil {
		t.Error("CopyRemoteFile() failed", err)
	}
}

func TestChangeRemoteFile(t *testing.T) {
	t.SkipNow()
	if err := util.CopyRemoteFile("/home/jerry/temp.json", "smb://192.168.8.140/public/sda1/Jerry/temp/temp.json"); err != nil {
		t.Error("TestChangeRemoteFile() failed", err)
	}
}

func TestDeleteRemoteFile(t *testing.T) {
	t.SkipNow()
	if err := util.DeleteRemoteFile("smb://192.168.8.140/public/sda1/Jerry/temp/temp.json"); err != nil {
		t.Error("TestDeleteRemoteFile() failed", err)
	}

	if err := util.DeleteRemoteFile("smb://192.168.8.140/public/sda1/Jerry/temp"); err != nil {
		t.Error("TestDeleteRemoteFile() failed", err)
	}
}
