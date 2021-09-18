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
