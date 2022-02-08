package auth

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const DIR = "atc"
const SESSION_FILE = "session.txt"

func GetSessionFilePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(cacheDir, DIR, SESSION_FILE)

	return path, nil
}

func CreateSessionDir() error {
	dir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, DIR)
	return os.MkdirAll(path, os.ModePerm)
}

func GetSession() (string, error) {
	path, err := GetSessionFilePath()
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	s := string(bytes)

	return strings.TrimRight(s, "\n"), nil
}

func StoreSession(byte []byte) error {
	if err := CreateSessionDir(); err != nil {
		return err
	}

	path, err := GetSessionFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(byte)
	if err != nil {
		return err
	}

	return nil
}

func ClearSession() error {
	path, err := GetSessionFilePath()
	if err != nil {
		return err
	}

	return os.Remove(path)
}
