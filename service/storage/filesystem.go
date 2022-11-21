package storage

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilesystemStorage implements Storage using the disk filesystem as a storage media
type FilesystemStorage struct {
}

const fsStorageRootDir = "static/user_content"

func (fs FilesystemStorage) SaveFile(name string, data []byte) (string, error) {
	// Create missing directories
	path := fs.pathForFileWithData(name, data)
	pathTree := strings.Split(path, "/")
	dirTree := pathTree[:len(pathTree)-1]

	err := os.MkdirAll(strings.Join(dirTree, "/"), os.ModePerm)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(path, data, 0o644&os.ModePerm)

	// Return a path starting with / to refer to the current server
	return "/" + path, err
}

func (fs FilesystemStorage) DeleteFile(path string) error {
	// Delete matching file
	fileToDelete := fs.pathForExistingFile(path)
	if fileToDelete == "" {
		return nil // Nothing to delete
	}

	// Delete found file
	if err := os.Remove(fileToDelete); err != nil {
		return err
	}

	// Delete parent directory, if empty
	dirPath := filepath.Dir(fileToDelete)
	if err := os.Remove(dirPath); err != nil {
		return err
	}

	return nil
}

// pathForExistingFile can only return paths for existing files.
//
// Since two different logical files with same filename cannot exist,
// we remove the only file with the matching extension
// in the directory that corresponds to the filename without extension.
func (fs FilesystemStorage) pathForExistingFile(applicationFileName string) string {
	// Parse the application provided file name / path
	pathParts := parseFilePath(applicationFileName)

	// List files inside corresponding directory
	fileDirEntries, err := os.ReadDir(fs.getFileDirPath(pathParts))
	if err != nil {
		return ""
	}

	for _, fileEntry := range fileDirEntries {
		if fileEntry.Type().IsRegular() && strings.HasSuffix(fileEntry.Name(), "."+pathParts.fileExtension) {
			// File found!
			return fs.getFileDirPath(pathParts) + "/" + fileEntry.Name()
		}
	}

	return "" // Not found
}

func (fs FilesystemStorage) pathForFileWithData(applicationFileName string, data []byte) string {
	// Parse the application provided file name / path
	pathParts := parseFilePath(applicationFileName)

	// Add SHA512 inside the filename
	hashSum := sha512.Sum512(data)
	hash := base64.URLEncoding.EncodeToString(hashSum[:])

	return fs.getFilePath(pathParts, hash)
}

func (FilesystemStorage) getFileDirPath(pathParts filePathParts) string {
	return fmt.Sprintf("%s/%s/%s", fsStorageRootDir, pathParts.pathPrefix, pathParts.filename)
}

// getFilePath returns the actual location.
// Stored inside <root>/<pathPrefix>/<filenameWithoutExtension>/<hash>.<extension>
func (fs FilesystemStorage) getFilePath(pathParts filePathParts, hash string) string {
	return fmt.Sprintf("%s/%s.%s", fs.getFileDirPath(pathParts), hash, pathParts.fileExtension)
}
