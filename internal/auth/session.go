package auth

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func IsExpired() (bool, error) {
	session, err := GetSession()
	if err != nil {
		return true, err
	}

	sv := strings.Split(session, "; ")
	for _, v := range sv {
		kv := strings.Split(v, "=")
		if len(kv) < 2 {
			continue
		}

		key := kv[0]
		if key != "Expires" {
			continue
		}

		value := kv[1]
		expiredDate, err := time.Parse(time.RFC1123, value)
		if err != nil {
			return true, err
		}

		if time.Now().After(expiredDate) {
			return true, nil
		}
	}

	return false, err
}

func ClearSession() error {
	path, err := GetSessionFilePath()
	if err != nil {
		return err
	}

	return os.Remove(path)
}
