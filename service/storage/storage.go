package storage

import (
	"path/filepath"
	"strings"
)

// Storage is an abstraction and can be any storage
// (filesystem, S3 bucket, CDN, in-memory, ...)
type Storage interface {
	// SaveFile gets the name of a file to save and the data as bytes.
	// It saves the file somewhere (according to the implementation),
	// then returns the locationUrl.
	//
	// The locationUrl can be a full HTTP URL or a relative URL starting with /.
	// In case it begins with /, the prefix is intended to be the current server.
	SaveFile(path string, data []byte) (locationUrl string, err error)

	// DeleteFile deletes a stored file given its path.
	// It should have been saved using the same Storage implementation.
	DeleteFile(path string) error

	// GetRoot returns the root path (on this device or somewhere else) that this storage implementation
	// is using to save given data.
	GetRoot() string
}

type filePathParts struct {
	pathPrefix    string
	filename      string
	fileExtension string
}

func parseFilePath(path string) filePathParts {
	pathParts := strings.Split(path, "/")
	pathPrefix := filepath.Dir(path)

	var filename, fileExtension string
	fullFileName := pathParts[len(pathParts)-1]
	if strings.Contains(fullFileName, ".") {
		// Filename is composed of both filename and extension
		fileNameParts := strings.Split(fullFileName, ".")
		filename = strings.Join(fileNameParts[:len(fileNameParts)-1], ".")
		fileExtension = fileNameParts[len(fileNameParts)-1]
	} else {
		// Filename doesn't have an extension
		filename = fullFileName
		fileExtension = ""
	}

	return filePathParts{
		pathPrefix:    pathPrefix,
		filename:      filename,
		fileExtension: fileExtension,
	}
}
