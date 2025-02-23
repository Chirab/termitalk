package app

import (
	"fmt"
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

const tmpDirName = "termitalk"
const tmpTokenfile = "token"

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

func getTmpDirPath(dirname string) string {
	return filepath.Join(os.TempDir(), dirname)
}

func (am *Memory) WriteToTokenFile(s string) error {
	fp := am.getTokenFile()
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(s)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil

}

func (am *Memory) getTokenFile() string {
	tmpDirPath := getTmpDirPath(tmpDirName)
	filePath := filepath.Join(tmpDirPath, tmpTokenfile)
	return filePath
}

func (am *Memory) InitMetadataStorage() error {
	tmpDirPath := getTmpDirPath(tmpDirName)
	filePath := filepath.Join(tmpDirPath, tmpTokenfile)
	if _, err := os.Stat(tmpDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(tmpDirPath, 0755)
		if err != nil {
			return err
		}

		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
		return nil
	}

	data, err := ioutil.ReadFile("example.txt")
	if err != nil {
		return err
	}
	am.SetServerSideToken(string(data))

	return nil
}
