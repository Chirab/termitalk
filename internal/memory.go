package app

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Memory struct {
	oauthId string

	serverSideToken string
	Username        string
}

type MYOS int64

const (
	MAC MYOS = iota
	LINUX
	WINDOWS
)

const tmpfilename = "termitalk"

func (am *Memory) GetAuthId() string {
	return am.oauthId
}

func (am *Memory) SetAuthId(oauthId string) {
	am.oauthId = oauthId
}

func (am *Memory) SetServerSideToken(i string) {
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

	file, err := os.OpenFile(tmpfilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(s)
	if err != nil {
		return err
	}
	return nil
}

func getTmpFilePath(filename string) string {
	tmpDir := os.TempDir()
	return filepath.Join(tmpDir, filename)
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

func (am *Memory) createTmpMetadataFile(s string) {
	file, err := os.Create(s)
	if err != nil {
		panic("cannot create metadata file")
	}

	defer file.Close()
}

func (am *Memory) InitMetadataStorage() error {
	tmpFilePath := getTmpFilePath(tmpfilename)

	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		file, err := os.Create(tmpFilePath)
		if err != nil {
			return err
		}

		defer file.Close()

	}

	data, err := ioutil.ReadFile(tmpFilePath)
	if err != nil {
		return err
	}

	am.SetServerSideToken(string(data))

	return nil
}
