package app

import (
	"errors"
	"os"
	"runtime"
)

type Memory struct {
	token string

	serverSideToken string
}

type MYOS int64

const (
	MAC MYOS = iota
	LINUX
	WINDOWS
)

const tmpfilename = "/tmp/termitalk"

func (am *Memory) GetAuthId() string {
	return am.token
}

func (am *Memory) SetAuthId(token string) {
	am.token = token
}

func (am *Memory) DecodeUserInfo(i string) {
	am.serverSideToken = i
}

func (am *Memory) GetServerSideToken() string {
	return am.serverSideToken
}

/*
 * Module to manage token in tmp folder
 */

func (am *Memory) WriteToTmpMetadataFile(s string) error {
	if s == "" {
		return errors.New("input to tmp file is empty")
	}

	_, err := file.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}

func (am *Memory) MyOs() MYOS {
	switch runtime.GOOS {
	case "windows":
		return WINDOWS
	case "darwin":
		return MAC
	case "linux":
		return LINUX
	default:
		panic("unrecognized os")
	}
}

func (am *Memory) IsTmpMetadataFile() bool {
	_, err := os.Stat(tmpfilename)
	switch {
	case err == nil:
		return true
	case os.IsNotExist(err):
		return false
	default:
		return false
	}
}

func (am *Memory) CreateTmpMetadataFile() {
	file, err := os.Create(tmpfilename)
	if err != nil {
		panic("cannot create metadata file")
	}

	defer file.Close()
}
