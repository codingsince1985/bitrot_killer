package util_test

import (
	"github.com/codingsince1985/bitrot_killer/util"
	"testing"
)

func TestGetFiles(t *testing.T) {
	path := "/home/jerry/Downloads"

	_, err := util.GetFiles(path)
	if err != nil {
		t.Error("GetFiles() failed", err)
	}
}

func TestMakeRemoteDir(t *testing.T) {
	err := util.MakeRemoteDir("smb://192.168.8.140/public/sda1/Jerry/temp")
	if err != nil {
		t.Error("TestMakeRemoteDir() failed", err)
	}
}

func TestCopyRemoteFile(t *testing.T) {
	err := util.CopyRemoteFile("/home/jerry/temp.json", "smb://192.168.8.140/public/sda1/Jerry/temp/temp.json")
	if err != nil {
		t.Error("CopyRemoteFile() failed", err)
	}
}

func TestChangeRemoteFile(t *testing.T) {
	err := util.CopyRemoteFile("/home/jerry/temp.json", "smb://192.168.8.140/public/sda1/Jerry/temp/temp.json")
	if err != nil {
		t.Error("TestChangeRemoteFile() failed", err)
	}
}

func TestDeleteRemoteFile(t *testing.T) {
	err := util.DeleteRemoteFile("smb://192.168.8.140/public/sda1/Jerry/temp/temp.json")
	if err != nil {
		t.Error("TestDeleteRemoteFile() failed", err)
	}

	err = util.DeleteRemoteFile("smb://192.168.8.140/public/sda1/Jerry/temp")
	if err != nil {
		t.Error("TestDeleteRemoteFile() failed", err)
	}
}
