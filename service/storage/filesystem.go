package storage

import (
	"os"
	"strings"
)

// FilesystemStorage implements Storage using the disk filesystem as a storage media
type FilesystemStorage struct {
}

func (FilesystemStorage) SaveFile(name string, data []byte) error {
	pathTree := strings.Split(name, "/")
	dirTree := pathTree[:len(pathTree)-1]
	dirPath := strings.Join(dirTree, "/")

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(name, data, 0o644&os.ModePerm)
}

func (FilesystemStorage) DeleteFile(path string) error {
	return os.Remove(path)
}
